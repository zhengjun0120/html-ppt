package agent

import (
	"context"
	"encoding/json"
	"errors"
	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/thumbs"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"
	"log"
	"strconv"

	"fmt"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

const maxTurns = 20 //最大允许调用20轮llm请求

// weekdayCN 给日期配一个中文星期：不配的话「上周末」「这周三」这类说法模型算不出来。
// 数组下标直接就是 time.Weekday（周日=0），不需要转换。
var weekdayCN = [...]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}

// runScope 一次 run 的执行作用域：工具集、配额表、轮数预算。
type runScope struct {
	tools     []openai.ChatCompletionToolUnionParam
	exec      map[string]ToolFunc
	maxPerRun map[string]int
	maxTurns  int
}

// v2 各阶段的默认轮数预算（config deck_v2.max_turns 可覆盖）。
var defaultStageMaxTurns = map[string]int{
	"clarifying":     8,
	"outlining":      10,
	"outline_review": 8,
	"generating":     24,
	"iterating":      20,
}

// maxTurnsLegacy 兜底轮数（未识别阶段/无配置时）。
const maxTurnsLegacy = 20

// scopeForStage 按阶段装配 v2 作用域。
func (as *AgentService) scopeForStage(stage string) *runScope {
	toolMap := as.buildToolsV2(stage)
	if toolMap == nil {
		toolMap = map[string]Tool{}
	}
	scope := &runScope{exec: map[string]ToolFunc{}, maxPerRun: map[string]int{}, maxTurns: maxTurnsLegacy}
	if mt, ok := as.StageMaxTurns[stage]; ok && mt > 0 {
		scope.maxTurns = mt
	} else if mt, ok := defaultStageMaxTurns[stage]; ok {
		scope.maxTurns = mt
	}
	for name, t := range toolMap {
		scope.tools = append(scope.tools, t.Definition)
		scope.exec[name] = t.Execute
		if t.MaxPerRun > 0 {
			scope.maxPerRun[name] = t.MaxPerRun
		}
	}
	return scope
}

// ErrStageLocked v2 阶段守卫的拒绝（handler 把它映射成 409 + 可执行的提示）。
type ErrStageLocked struct{ Msg string }

func (e ErrStageLocked) Error() string { return e.Msg }

// scopeDecision 一次 chat run 的路由结论：作用域 + 系统提示词。
// v2 会话每个 run 重刷系统提示词（阶段会跨轮迁移，旧 system 必须跟着换）。
type scopeDecision struct {
	scope *runScope
	sys   string
	isV2  bool
}

// scopeForSession 决定一次 chat run 的作用域与系统提示词（deck-v2 的阶段路由核心）。
//
//	deckID 空：冷启动会话。首轮 = clarifying（提问）；ask_user 恢复轮 = outlining
//	（答案已到手，该产出大纲了）。
//	deckID 是 v2 deck：按 deck.stage 路由，selecting_template/generating 阶段拒绝对话。
//	deckID 不是 v2：旧格式文稿已随旧栈下线，明确报错（D7）。
func (as *AgentService) scopeForSession(userID uint, sess *store.ChatSession, deckIDParam string, continuation bool) (scopeDecision, error) {
	deckID := sess.DeckID
	if deckID == "" {
		deckID = deckIDParam
	}
	if deckID == "" {
		if continuation {
			// 冷启动的 ask_user 恢复轮：澄清的答案已到手，该产出大纲了
			return scopeDecision{scope: as.scopeForStage("outlining"), sys: BuildStagePrompt("outlining", "", nil, nil), isV2: true}, nil
		}
		return scopeDecision{scope: as.scopeForStage("clarifying"), sys: BuildStagePrompt("", "", nil, nil), isV2: true}, nil
	}
	if !as.DeckService.IsV2(deckID) {
		return scopeDecision{}, ErrStageLocked{Msg: "这份文稿是旧格式（deck-v1），已随新版下线；请新建文稿生成"}
	}

	df, err := as.DeckService.GetDeckV2(userID, deckID)
	if err != nil {
		return scopeDecision{}, err
	}
	switch df.Stage {
	case deck.StageSelectingTemplate:
		return scopeDecision{}, ErrStageLocked{Msg: "请先在界面上选择模板（这一步在对话里做不了）"}
	case deck.StageGenerating:
		return scopeDecision{}, ErrStageLocked{Msg: "deck 正在生成中：请等生成完成，或刷新页面后用「继续生成」恢复"}
	case deck.StageOutlining, deck.StageOutlineReview:
		return scopeDecision{scope: as.scopeForStage(df.Stage), sys: BuildStagePrompt(df.Stage, deckID, as.templateOrNil(df), nil), isV2: true}, nil
	case deck.StageIterating:
		o, _ := as.DeckService.ReadOutline(userID, deckID)
		return scopeDecision{scope: as.scopeForStage(deck.StageIterating), sys: BuildStagePrompt(deck.StageIterating, deckID, as.templateOrNil(df), o), isV2: true}, nil
	default:
		return scopeDecision{}, fmt.Errorf("deck 处于未知阶段 %q", df.Stage)
	}
}

// templateOrNil deck → 模板（未选模板/库不可用时 nil，提示词相应少一段）。
func (as *AgentService) templateOrNil(df *deck.DeckFile) *template.Template {
	if as.Templates == nil || df.TemplateID == "" {
		return nil
	}
	tpl, err := as.Templates.Get(df.TemplateID)
	if err != nil {
		return nil
	}
	return tpl
}

// refreshSystem 把会话消息数组里的系统提示词替换为当前阶段的版本。
// 不做类型探测：v2 会话的 messages[0] 恒为 system（第一次 refreshSystem 时前置的，
// 之后只会被覆盖），这个不变式由 v2 会话只走 refreshSystem 保证。
func refreshSystem(messages []openai.ChatCompletionMessageParamUnion, systemPrompt string) []openai.ChatCompletionMessageParamUnion {
	if len(messages) == 0 {
		return append([]openai.ChatCompletionMessageParamUnion{openai.SystemMessage(systemPrompt)}, messages...)
	}
	messages[0] = openai.SystemMessage(systemPrompt)
	return messages
}

func (as *AgentService) NewStreamChat(ctx context.Context, userID, sessionID uint, userContent, deckID string, emit func(StreamEvent) error) (uint, error) {
	ctx = authctx.WithUser(ctx, userID)
	client := as.clientFor(ctx)

	// 读取会话历史
	sess, messages, err := as.loadOrCreateSession(ctx, userID, sessionID, userContent, deckID)
	if err != nil {
		return 0, err
	}

	// 闸门：上一条提问还没回答时，不接受新的用户消息。
	//
	// 不拦的后果比"这一轮报错"严重得多：暂停态的消息数组结尾是一条带 tool_calls 的
	// assistant 消息（ask_user 的应答要等用户回答才补上）。此时追加用户消息就破坏了
	// "tool_calls 后面必须紧跟 tool 消息"这条协议不变式，API 回
	//   400 An assistant message with 'tool_calls' must be followed by tool messages…
	// 而这轮的用户消息**已经被下面的 persistSession 写进库了**，于是这个会话之后
	// 每次请求都带着坏形状 400，连正常回答都救不回来（AnswerChat 把答案追加到末尾，
	// 仍然不是紧跟 tool_calls 的位置）。所以这里一律拒绝：用户要么把提问答完，
	// 要么开新对话。
	if err := guardNewMessage(sess); err != nil {
		return sess.ID, err
	}

	// deck-v2 阶段路由：作用域 + 系统提示词（v1 会话走旧路径，行为不变）
	dec, err := as.scopeForSession(userID, sess, deckID, false)
	if err != nil {
		return sess.ID, err
	}
	if sess.DeckID == "" && deckID != "" {
		sess.DeckID = deckID // 首轮把请求里的 deck 绑到会话（v1 语义保留）
	}
	if dec.isV2 {
		messages = refreshSystem(messages, dec.sys)
	} else if sessionID == 0 {
		messages = append(messages, openai.SystemMessage(dec.sys))
	}

	messages = append(messages, openai.UserMessage(userContent))
	if err := as.persistSession(sess, messages); err != nil {
		return sess.ID, err
	}
	if emitErr := emit(StreamEvent{Type: EventTypeSession, Content: strconv.FormatUint(uint64(sess.ID), 10)}); emitErr != nil {
		return sess.ID, fmt.Errorf("事件推送失败: %w", emitErr)
	}

	paused, rr, err := as.runLoop(ctx, client, sess, messages, emit, runTraceInfo{
		UserID:      userID,
		DeckID:      sess.DeckID,
		UserContent: userContent,
	}, dec.scope)
	as.bindDeckFromRun(sess, rr)
	if paused {
		return sess.ID, ErrPaused
	}
	if err != nil {
		return sess.ID, err
	}
	return sess.ID, nil
}

// isLegacyScope 已移除：isV2 标志在 scopeDecision 上，v1 会话不再触碰系统提示词。

// bindDeckFromRun 会话绑定：冷启动会话里 submit_outline 创建了 deck，
// run 结束（含暂停）后把它绑回会话——之后的轮次才有 stage 可路由。
func (as *AgentService) bindDeckFromRun(sess *store.ChatSession, rr *runRecorder) {
	if sess == nil || rr == nil || sess.DeckID != "" || as.st == nil {
		return
	}
	ids := rr.deckIDs()
	if len(ids) != 1 {
		return
	}
	if err := as.st.DB.Model(&store.ChatSession{}).Where("id = ?", sess.ID).Update("deck_id", ids[0]).Error; err != nil {
		log.Printf("[warn] 会话 %d 绑定 deck %s 失败: %v", sess.ID, ids[0], err)
		return
	}
	sess.DeckID = ids[0]
}

// AnswerChat 用户回答 ask_user 后恢复循环
func (as *AgentService) AnswerChat(ctx context.Context, userID, sessionID uint, answersJSON string, emit func(StreamEvent) error) (uint, error) {
	ctx = authctx.WithUser(ctx, userID)
	client := as.clientFor(ctx)

	sess, messages, err := as.loadSession(ctx, userID, sessionID)
	if err != nil {
		return 0, err
	}
	if sess.PendingAsk == "" {
		return sess.ID, errors.New("会话不在等待回答状态")
	}
	pending := parsePendingAsk(sess)
	if pending.ToolCallID == "" {
		return sess.ID, errors.New("暂停态损坏 err: tool_call_id 为空")
	}

	// 用"设置答案"而不是"追加一条 tool 消息"：加载时已经修复过，
	// 所以这个 tool_call_id 必然已经**有一条**应答（正常路径下是修复补的兜底说明，
	// 被历史 bug 弄坏过的会话里则是上一次回答留下的内容）。追加会让同一个 id
	// 出现两条应答，API 同样会拒。
	messages = setToolAnswer(messages, pending.ToolCallID, answersJSON)
	sess.PendingAsk = ""
	if err := as.persistSession(sess, messages); err != nil {
		return sess.ID, err
	}
	if emitErr := emit(StreamEvent{Type: EventTypeSession, Content: strconv.FormatUint(uint64(sess.ID), 10)}); emitErr != nil {
		return sess.ID, fmt.Errorf("事件推送失败 err:%w", emitErr)
	}

	// 恢复也算一次新的 run（消息数组继续用，但 run 的边界是"一次循环执行"）。
	// parent_run_id 指回被暂停的那次，观测页据此把被打断的对话串起来看。
	// 阶段路由在恢复轮重新走：冷启动会话的恢复轮 = outlining（答案已到手）。
	dec, err := as.scopeForSession(userID, sess, "", true)
	if err != nil {
		return sess.ID, err
	}
	if dec.isV2 {
		messages = refreshSystem(messages, dec.sys)
	}
	paused, rr, err := as.runLoop(ctx, client, sess, messages, emit, runTraceInfo{
		UserID:      userID,
		DeckID:      sess.DeckID,
		UserContent: "(用户回答了 ask_user 的提问)",
		ParentRunID: pending.RunID,
	}, dec.scope)
	as.bindDeckFromRun(sess, rr)
	if paused {
		return sess.ID, ErrPaused
	}
	if err != nil {
		return sess.ID, err
	}
	return sess.ID, nil
}

func (as *AgentService) runLoop(ctx context.Context, client *openai.Client, sess *store.ChatSession, messages []openai.ChatCompletionMessageParamUnion, emit func(StreamEvent) error, info runTraceInfo, scope *runScope) (paused bool, rr *runRecorder, err error) {
	rr = newRunRecorder()
	defer func() {
		if err == nil {
			as.recordRunVersions(ctx, rr)
		}
	}()

	// —— 观测（trace）——
	// 顺序很重要：先声明 Close、后声明下面那个 run_end 兜底。defer 是后进先出，
	// 所以兜底会先执行、再关文件。反过来写的话，中途出错那一轮的 run_end 会落进
	// 一个已经关闭的文件（writeLocked 直接丢弃），观测页上那次运行就永远停在
	// "运行中"并一直轮询——而真实情况是它早就挂了。
	info.SessionID = sess.ID
	if uid, ok := authctx.UserID(ctx); ok {
		info.UserID = uid
	}
	rec := as.openTraceRecorder(info, emit)
	defer rec.Close()
	ctx = trace.With(ctx, rec)

	ended := false
	endRun := func(status string) {
		if ended {
			return
		}
		ended = true
		s := rec.Summary()
		trace.Emit(ctx, trace.Event{Kind: trace.KindRunEnd, Status: status, Summary: &s})
	}
	defer func() {
		if ended {
			return
		}
		// 走到这里说明是"中途返回了错误"（流断了、写库失败、事件推不动）。
		// 错误原因必须进 trace：服务端日志会被滚动冲掉，而这条记录是要长期回看的
		if err != nil {
			trace.Emit(ctx, trace.Event{Kind: trace.KindError, Error: err.Error()})
		}
		endRun(trace.StatusError)
	}()

	opt := as.setChatOpts()
	opt.Messages = messages
	opt.StreamOptions = openai.ChatCompletionStreamOptionsParam{IncludeUsage: openai.Bool(true)}
	opt.Tools = scope.tools

	var fullText strings.Builder

	var (
		promptTokens     int64 = 0 //输入token
		completionTokens int64 = 0 //输出token
		totalTokens      int64 = 0 //总token
		cachedTokens     int64 = 0 //缓存命中token
	)

	for i := 0; ; i++ {
		// 这一轮的所有事件（请求、工具、子过程）都归到 turn=i 下。
		// 工具调用从 turnCtx 派生自己的工具作用域，兄弟工具之间不会互相污染
		turnCtx := trace.WithTurn(ctx, i)
		//预算将尽：提前吹哨，让模型把剩余轮次留给硬伤，而不是被硬掐后仓促收尾
		if warn := scope.maxTurns - 4; warn > 0 && i == warn {
			messages = append(messages, openai.UserMessage("工具调用预算还剩 4 轮：只修硬伤（溢出/截断/拒收的页），停止打磨性改动，然后准备收尾汇报"))
			opt.Messages = messages
			if err := as.persistSession(sess, messages); err != nil {
				return false, rr, err
			}
		}
		//预算花完了
		if i > scope.maxTurns {
			messages = append(messages, openai.UserMessage("工具调用预算已用完：不要再调用任何工具，直接基于以上获取的信息给出最终回答"))
			opt.Messages = messages
			opt.Tools = nil
			err := as.persistSession(sess, messages)
			if err != nil {
				return false, rr, err
			}
		}
		msg, usage, streamErr := as.streamOnce(turnCtx, client, opt, &fullText, emit)
		if streamErr != nil {
			return false, rr, streamErr
		}

		promptTokens += usage.PromptTokens
		completionTokens += usage.CompletionTokens
		totalTokens += usage.TotalTokens
		cachedTokens += usage.PromptTokensDetails.CachedTokens

		// 分项用量：主循环这一次调用单独记一笔。过去这里只把四个数字累加起来
		// 最后发一次，既看不出"第 3 轮花了多少"，也把视觉/联网那两个子调用的
		// 用量整个漏掉了（它们是分开计费的另外几次 API 调用）。
		trace.Usage(turnCtx, trace.CompMain, usagePartFrom(usage))

		if len(msg.ToolCalls) == 0 {
			messages = append(messages, msg.ToParam())
			opt.Messages = messages
			err := as.persistSession(sess, messages)
			if err != nil {
				return false, rr, err
			}
			// 汇总带上分项：扁平那四个字段保持"主循环口径"不动（前端契约），
			// 真实总成本看 Usage.Total
			summary := rec.Summary()
			if emitErr := emit(StreamEvent{Type: EventTypeDone, Content: fullText.String(), PromptTokens: promptTokens, CompletionTokens: completionTokens, TotalTokens: totalTokens, CachedTokens: cachedTokens, Usage: &summary}); emitErr != nil {
				return false, rr, fmt.Errorf("事件推送失败: %w", emitErr)
			}
			endRun(trace.StatusOK)
			return false, rr, nil
		}

		messages = append(messages, msg.ToParam())
		opt.Messages = messages
		if err := as.persistSession(sess, messages); err != nil {
			return false, rr, err
		}

		for _, tool := range msg.ToolCalls {
			//判断是不是要调用提问用户
			if tool.Function.Name == "ask_user" {
				// ask_user 由循环拦截、不会走到 execTool，所以它的观测要在这里补：
				// 少了这几条，"这次对话为什么停住了"在 trace 里是个空洞
				askCtx := trace.WithTool(turnCtx, "ask_user", tool.ID)
				trace.Emit(askCtx, trace.Event{Kind: trace.KindToolCall, Args: tool.Function.Arguments})

				var args AskUserArgs
				//解析参数
				if err := json.Unmarshal([]byte(tool.Function.Arguments), &args); err != nil || len(args.Questions) == 0 || len(args.Questions) > 6 {
					//如果json不合法 把错误信息直接加入到上下文中 然后跳过
					if err != nil {
						result := fmt.Sprintf("ask_user 参数json不合法 err:%s", err.Error())
						trace.Emit(askCtx, trace.Event{Kind: trace.KindToolResult, Result: result, Error: err.Error()})
						messages = append(messages, openai.ToolMessage(result, tool.ID))
						opt.Messages = messages
						if err := as.persistSession(sess, messages); err != nil {
							return false, rr, err
						}
						continue
					}
					//如果提问的问题大于6个小于1个 同上
					result := "ask_user 参数不合法: questions 必须是 1~6 个问题"
					trace.Emit(askCtx, trace.Event{Kind: trace.KindToolResult, Result: result, Error: "questions 数量不在 1~6"})
					messages = append(messages, openai.ToolMessage(result, tool.ID))
					opt.Messages = messages
					if err := as.persistSession(sess, messages); err != nil {
						return false, rr, err
					}
					continue
				}
				//如果要调用提问用户 则不可调用其他工具
				for _, other := range msg.ToolCalls {
					if other.ID != tool.ID {
						note := `{"note":"用户被提问打断，本工具本轮未执行，用户回答后如仍需要请重新调用"}`
						messages = append(messages, openai.ToolMessage(note, other.ID))
						// 这些"看起来被调用了、实际没执行"的工具最容易被误读成"模型做了 X"，
						// 所以在 trace 里也要显式留一条，而不是只在上下文里体现
						trace.Emit(trace.WithTool(turnCtx, other.Function.Name, other.ID),
							trace.Event{Kind: trace.KindToolResult, Result: note, Error: "被 ask_user 打断，本轮未执行"})
					}
				}
				opt.Messages = messages
				// run_id 一起存进暂停态：恢复时作为 parent_run_id 带回来，
				// 观测页才能把被打断的这次对话和一个 run 串起来。
				// 注意这条 assistant 消息的 tool_calls 里，ask_user 那个**故意**没有应答
				//（正在等用户回答）——所以此时的 messages 是不满足协议不变式的，
				// 它不能直接发给模型。这就是 NewStreamChat 要设闸门的原因。
				sess.PendingAsk = encodePendingAsk(pendingAsk{ToolCallID: tool.ID, RunID: rec.RunID()})
				if err := as.persistSession(sess, messages); err != nil {
					return false, rr, err
				}
				if emitErr := emit(StreamEvent{Type: EventTypeAskUser, ToolCallID: tool.ID, Content: tool.Function.Arguments}); emitErr != nil {
					return false, rr, fmt.Errorf("事件推送失败: %w", emitErr)
				}
				trace.Emit(askCtx, trace.Event{Kind: trace.KindToolResult, Result: `{"note":"已向用户提问，本轮暂停等回答"}`})
				endRun(trace.StatusPaused)
				return true, rr, nil
			}
			//普通工具调用
			result, execErr := as.execTool(turnCtx, tool, emit, rr, scope)
			if execErr != nil {
				return false, rr, execErr
			}
			messages = append(messages, openai.ToolMessage(result, tool.ID))
			opt.Messages = messages
			if err := as.persistSession(sess, messages); err != nil {
				return false, rr, err
			}
		}
		// 继续下一轮
	}
}

type emitKey struct{}

func withEmit(ctx context.Context, emit func(StreamEvent) error) context.Context {
	return context.WithValue(ctx, emitKey{}, emit)
}

func emitFrom(ctx context.Context) func(StreamEvent) error {
	e, _ := ctx.Value(emitKey{}).(func(StreamEvent) error)
	return e
}

func (as *AgentService) execTool(ctx context.Context, tool openai.ChatCompletionMessageToolCallUnion, emit func(StreamEvent) error, rr *runRecorder, scope *runScope) (string, error) {
	if scope == nil {
		// 白盒测试路径：nil scope 回落全局装配（生产路径恒有 scope）
		scope = &runScope{exec: as.Exec, maxPerRun: as.MaxPerRun}
	}
	// 把归属挂进 ctx：工具内部（联网搜索的每一次子搜索、视觉审查的每一页截图）
	// 因此不必知道自己在第几轮、call_id 是什么，只要 trace.Emit 就会自动归位
	ctx = trace.WithTool(ctx, tool.Function.Name, tool.ID)
	ctx = withEmit(ctx, func(e StreamEvent) error {
		if e.ToolCallID == "" {
			e.ToolCallID = tool.ID
		}
		return emit(e)
	}) // 子调用 可用它把过程增量推给前端
	trace.Emit(ctx, trace.Event{Kind: trace.KindToolCall, Args: tool.Function.Arguments})

	var (
		result    string
		callErr   string // 空串 = 没失败。刻意不用 error 类型：它同时要进 JSON
		startedAt = time.Now()
	)

	t, ok := scope.exec[tool.Function.Name]
	// runTool 执行工具并处理结果。配额分流有两个入口（数字复查跳过配额直接执行、
	// 其余走闸门），但执行与失败处理是同一份，抽出来避免抄两遍。
	runTool := func() {
		res, err := t(ctx, tool.Function.Arguments)
		if err != nil {
			result = fmt.Sprintf("工具调用失败 err:%s", err.Error())
			callErr = err.Error()
		} else {
			result = res
			// review_slides 的配额只数"真的看了图"的调用：看图失败的降级结果
			// 不带 visionReportMarker，这里退还计数（pages 留空的数字复查在
			// 上一个 case 里根本不进配额，无需退还）。实测（trace 1789613158323-30a4）：
			// 两次免费复查加一次看图失败就把 3 次配额耗尽，模型想重试真审查时被拒——
			// 只数有产出的调用才符合"拦住审查循环"的本意。
			if tool.Function.Name == "review_slides" && !strings.Contains(res, visionReportMarker) {
				rr.refund("review_slides")
			}
			rr.noteToolCall(tool.Function.Name, tool.Function.Arguments, res) //版本备注只记成功的
		}
	}
	switch {
	case !ok:
		result = "未知工具"
		callErr = result
	case tool.Function.Name == "review_slides" && reviewMeasureOnlyArgs(tool.Function.Arguments):
		// 数字复查（pages 留空）不进配额、直接执行。必须放在闸门之前：
		// overQuota 是"先计数再比较"，被拦的真审查也会留下计数，若数字复查走闸门，
		// 它会在配额耗尽后被误拦——提示词承诺它"免费、随时可调"。
		// 解析失败按真审查处理（宁可少放行，不可放开真审查的闸）。
		runTool()
	case rr.overQuota(tool.Function.Name, scope.maxPerRun[tool.Function.Name]):
		// 配额闸门（Tool.MaxPerRun）：审查→修复→再审这类循环没有自然出口，
		// 提示词拦不住，就在这里拦。刻意不算 callErr：这是策略结果不是故障，
		// 模型要的是"接下来该怎么办"的指令，前端也不该把它显示成红色失败
		result = fmt.Sprintf("配额用完：%s 在一次 run 里最多执行 %d 次，已用尽（本次调用未执行）。"+
			"不要再调用它——基于已有的结果继续完成任务，剩余想检查/想打磨的点在最终汇报里说明。",
			tool.Function.Name, as.MaxPerRun[tool.Function.Name])
	default:
		runTool()
	}

	// 失败也要记。原来 runRecorder 只记成功调用，于是"模型连试三次都被拒"这种
	// 最需要解释的现象，在记录里完全看不到——你只看得到它最后放弃了。
	trace.Emit(ctx, trace.Event{
		Kind: trace.KindToolResult, Result: result, Error: callErr,
		DurationMS: time.Since(startedAt).Milliseconds(),
	})

	if callErr != "" {
		if emitErr := emit(StreamEvent{Type: EventTypeToolError, Content: result, ToolCallID: tool.ID, ToolName: tool.Function.Name}); emitErr != nil {
			return result, fmt.Errorf("事件推送失败: %w", emitErr)
		}
	} else {
		if emitErr := emit(StreamEvent{Type: EventTypeToolCall, Content: result, ToolCallID: tool.ID, ToolName: tool.Function.Name}); emitErr != nil {
			// 这里原来写的是 `return result, err`，而那个 err 是上层已判为 nil 的变量，
			// 于是"推送失败"被静默吞掉、循环拿着一个根本没发出去的工具结果继续往下跑
			return result, fmt.Errorf("事件推送失败: %w", emitErr)
		}
	}

	return result, nil
}

func (as *AgentService) streamOnce(ctx context.Context, client *openai.Client, opt openai.ChatCompletionNewParams, fullText *strings.Builder, emit func(StreamEvent) error) (openai.ChatCompletionMessage, openai.CompletionUsage, error) {
	// 观测：把这一轮**实际发给模型的东西**记下来。这是"模型为什么这么答"唯一能查的证据——
	// 只记工具的输入输出是不够的：同一个工具结果，在不同上下文里会被理解成不同的意思。
	//
	// trace.Active 这层判断是必要的：序列化整份上下文是 MB 级的开销（每轮重发整份消息），
	// 关掉观测就不该白做一遍。其余地方的 Emit 不用套它——Discard 本身就是空操作。
	if trace.Active(ctx) {
		raw, err := json.Marshal(opt.Messages)
		if err != nil {
			// 上下文序列化不了本身就是个信号（比如混进了不可序列化的类型），
			// 但别因此中断对话：记一条空占位，让"这一轮没留下上下文"是可见的
			log.Printf("[warn] trace: 序列化上下文失败 err: %v", err)
			raw = json.RawMessage(`null`)
		}
		trace.Emit(ctx, trace.Event{
			Kind:         trace.KindLLMRequest,
			Model:        as.ModelID,
			Messages:     raw,
			MessageCount: len(opt.Messages),
			Bytes:        len(raw),
		})
	}

	stream := client.Chat.Completions.NewStreaming(ctx, opt)
	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		if !acc.AddChunk(chunk) {
			continue
		}
		//尾分片 无choices 跳过防panic
		if len(chunk.Choices) == 0 {
			continue
		}
		delta := chunk.Choices[0].Delta
		for _, tc := range delta.ToolCalls {
			idx := tc.Index
			if idx < 0 {
				idx = 0 // 个别网关对单个工具调用用 -1，sdk 累加器内部就归到0
			}
			if tc.ID != "" {
				// 只有首个片段会带id和完整name 后续片段只有index+arguments
				if emitErr := emit(StreamEvent{Type: EventTypeToolStart, ToolIndex: idx, ToolCallID: tc.ID, ToolName: tc.Function.Name}); emitErr != nil {
					stream.Close()
					return openai.ChatCompletionMessage{}, openai.CompletionUsage{}, fmt.Errorf("事件推送失败: %w", emitErr)
				}
			}
			// 流式推送参数
			if tc.Function.Arguments != "" {
				if emitErr := emit(StreamEvent{Type: EventTypeToolDelta, ToolIndex: idx, Content: tc.Function.Arguments}); emitErr != nil {
					stream.Close()
					return openai.ChatCompletionMessage{}, openai.CompletionUsage{}, fmt.Errorf("事件推送失败: %w", emitErr)
				}
			}
		}

		if delta.Content != "" {
			fullText.WriteString(delta.Content)
			if emitErr := emit(StreamEvent{Type: EventTypeDelta, Content: delta.Content}); emitErr != nil {
				stream.Close()
				return openai.ChatCompletionMessage{}, openai.CompletionUsage{}, fmt.Errorf("事件推送失败: %w", emitErr)
			}
		}
		if f, ok := delta.JSON.ExtraFields["reasoning_content"]; ok {
			var reasoning string
			if err := json.Unmarshal([]byte(f.Raw()), &reasoning); err == nil && reasoning != "" {
				if emitErr := emit(StreamEvent{Type: EventTypeThink, Content: reasoning}); emitErr != nil {
					stream.Close()
					return openai.ChatCompletionMessage{}, openai.CompletionUsage{}, fmt.Errorf("事件推送失败: %w", emitErr)
				}
			}
		}
	}
	if err := stream.Err(); err != nil {
		stream.Close()
		trace.Emit(ctx, trace.Event{Kind: trace.KindLLMResponse, Error: "流式输出错误: " + err.Error()})
		return openai.ChatCompletionMessage{}, openai.CompletionUsage{}, fmt.Errorf("流式输出错误: %w", err)
	}
	stream.Close()

	if len(acc.Choices) == 0 || acc.Choices[0].FinishReason == "" {
		trace.Emit(ctx, trace.Event{Kind: trace.KindLLMResponse, Error: "流式响应不完整（没有 finish_reason）"})
		return openai.ChatCompletionMessage{}, openai.CompletionUsage{}, errors.New("流式响应不完整，请重试")
	}

	// 观测：模型回了什么。文本、工具调用意图、finish_reason、以及这一次的用量。
	// 与上面那条 llm_request 配对看，才能回答"同样的工具结果，模型这次为什么改主意了"
	msg := acc.Choices[0].Message
	part := usagePartFrom(acc.Usage)
	trace.Emit(ctx, trace.Event{
		Kind:         trace.KindLLMResponse,
		FinishReason: string(acc.Choices[0].FinishReason),
		Content:      msg.Content,
		ToolCalls:    toolCallsOut(msg.ToolCalls),
		Usage:        &part,
	})
	return msg, acc.Usage, nil
}

func (as *AgentService) setChatOpts() openai.ChatCompletionNewParams {
	opt := openai.ChatCompletionNewParams{
		Model:           openai.ChatModel(as.ModelID),
		ReasoningEffort: shared.ReasoningEffortHigh,
	}
	return opt
}

func (as *AgentService) recordRunVersions(ctx context.Context, rr *runRecorder) {
	uid, ok := authctx.UserID(ctx)
	if !ok {
		return
	}
	for _, deckID := range rr.deckIDs() {
		if err := as.DeckService.RecordRunVersion(uid, deckID, rr.detail(deckID)); err != nil {
			log.Printf("[warn] deck %s 记录历史版本失败 err: %v", deckID, err)
		}
	}
}

// StartGenerationRun 生成 run 的发起（POST /api/decks/:id/generate）。
//
// 与 chat run 的区别：messages 是**全新的一份**（system=generate 提示词 + 一句启动指令），
// 整体覆盖会话消息——大纲产物已持久化在 deck 目录里，生成 run 是自包含的执行，
// 会话回放让位给"生成全记录"。断线恢复：每批 write_pages 落盘即持久化，
// stage 停在 generating，resume=true 的 run 会先对齐现状再续写。
func (as *AgentService) StartGenerationRun(ctx context.Context, userID, sessionID uint, deckIDParam string, resume bool, emit func(StreamEvent) error) (uint, error) {
	ctx = authctx.WithUser(ctx, userID)
	client := as.clientFor(ctx)

	sess, _, err := as.loadSession(ctx, userID, sessionID)
	if err != nil {
		return 0, err
	}
	if sess.DeckID == "" {
		// 会话可能没绑上 deck（ask_user 暂停期间的绑定不发生）。generate 的 URL
		// 里就带着 deck id——这里直接绑，而不是报错让用户重来。
		if deckIDParam == "" || as.DeckService.EnsureOwner(userID, deckIDParam) != nil {
			return sess.ID, errors.New("会话没有关联的 deck：请先完成大纲与模板选择")
		}
		if as.st != nil {
			if err := as.st.DB.Model(&store.ChatSession{}).Where("id = ?", sess.ID).Update("deck_id", deckIDParam).Error; err == nil {
				sess.DeckID = deckIDParam
			}
		}
		if sess.DeckID == "" {
			return sess.ID, errors.New("绑定 deck 失败")
		}
	}
	df, err := as.DeckService.GetDeckV2(userID, sess.DeckID)
	if err != nil {
		return sess.ID, err
	}
	if df.Stage != deck.StageGenerating {
		return sess.ID, ErrStageLocked{Msg: fmt.Sprintf("deck 阶段为 %s，不在生成态", df.Stage)}
	}

	outline, _ := as.DeckService.ReadOutline(userID, sess.DeckID)

	var userMsg string
	if resume {
		userMsg = "继续生成：先 list_slides 对齐已写入的页，然后从缺失的页继续 write_pages；" +
			"若全部页已写入，则按最近一次量测结果修复问题页，修完汇报。"
	} else {
		userMsg = "开始生成：按工作流走——先 plan_pages 全局规划（被打回就调整重提），" +
			"然后分批 write_pages 写完全部页，处理量测与 lint 反馈，需要时 review_slides 看图，最后如实汇报。"
	}
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(appendDate(BuildStagePrompt(deck.StageGenerating, sess.DeckID, as.templateOrNil(df), outline))),
		openai.UserMessage(userMsg),
	}
	if err := as.persistSession(sess, messages); err != nil {
		return sess.ID, err
	}
	if emitErr := emit(StreamEvent{Type: EventTypeSession, Content: strconv.FormatUint(uint64(sess.ID), 10)}); emitErr != nil {
		return sess.ID, fmt.Errorf("事件推送失败: %w", emitErr)
	}

	scope := as.scopeForStage(deck.StageGenerating)
	_, rr, err := as.runLoop(ctx, client, sess, messages, emit, runTraceInfo{
		UserID:      userID,
		DeckID:      sess.DeckID,
		UserContent: userMsg,
	}, scope)
	if err == nil {
		// 正常结束 → iterating。中断/错误时 stage 留在 generating（resume 靠它）。
		if to, terr := as.DeckService.FinishGeneration(userID, sess.DeckID); terr != nil {
			log.Printf("[warn] deck %s 推进到 iterating 失败: %v", sess.DeckID, terr)
		} else {
			as.emitV2Public(emit, EventTypeStage, map[string]any{"deck_id": sess.DeckID, "to": to})
			// 异步预热缩略图：文稿列表封面与预览栏首屏不用等首次点开时的
			// 整本渲染（约 10-20s）。失败静默——懒加载路径会兜底重试。
			go as.warmThumbs(userID, sess.DeckID)
		}
	}
	_ = rr
	return sess.ID, err
}

// warmThumbs 生成结束后整本预渲染缩略图（后台执行，与请求生命周期解耦）。
func (as *AgentService) warmThumbs(userID uint, deckID string) {
	defer func() { _ = recover() }() // 预热是锦上添花，任何异常都不许影响主流程
	if !as.Vision || as.VisionGrants == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	nonce, err := as.VisionGrants.Issue(userID, deckID)
	if err != nil {
		return
	}
	url := strings.TrimRight(as.VisionBaseURL, "/") + "/api/render/" + nonce
	d, err := vision.CaptureV2(ctx, vision.OptionsV2{URL: url, ChromePath: as.ChromePath})
	if err != nil {
		log.Printf("[warn] deck %s 缩略图预热失败（懒加载兜底）: %v", deckID, err)
		return
	}
	ip, err := as.DeckService.IndexPath(deckID)
	if err != nil {
		return
	}
	pngs := make(map[int][]byte, len(d.Slides))
	for _, sl := range d.Slides {
		pngs[sl.Index+1] = sl.PNG
	}
	if err := thumbs.WriteAll(as.DeckService.ThumbsDir(deckID), ip, pngs); err != nil {
		log.Printf("[warn] deck %s 缩略图落盘失败: %v", deckID, err)
	}
}

// emitV2Public 生成 run 的管线事件（runLoop 之外没有 ctx emit 链，直连 emit）。
func (as *AgentService) emitV2Public(emit func(StreamEvent) error, eventType string, payload map[string]any) {
	if emit == nil {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_ = emit(StreamEvent{Type: eventType, Content: string(data)})
}
