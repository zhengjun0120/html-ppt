package agent

import (
	"context"
	"encoding/json"
	"errors"
	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
	"log"
	"strconv"

	"fmt"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

const maxTurns = 20 //最大允许调用20轮llm请求

// func (as *AgentService) buildMessages(userContent,deckID string) ([]openai.ChatCompletionMessageParamUnion,error){
// 	//TODO 从数据库中获取对话记录，如果数据库中没有 则先拼接系统提示词 并新建数据

// 	//TODO 然后拼接用户对话

// 	// Test
// 	messages := make([]openai.ChatCompletionMessageParamUnion,0)
// 	systemMessage := systemPrompt
// 	if deckID != "" && deck.IsValidID(deckID){
// 		systemMessage += fmt.Sprintf(`当前用户正在预览的演示文稿是 {"deck_id":"%s"}，涉及它的修改直接用这个 deck_id，不要向用户询问`,deckID)
// 	}
// 	messages = append(messages, openai.SystemMessage(systemMessage))
// 	messages = append(messages, openai.UserMessage(userContent))
// 	fmt.Println("系统提示词："+systemMessage)
// 	return messages,nil
// }

// StreamChat 以 userID 的身份执行一轮 agent 循环。
// userID 来自 JWT 中间件（经 request context 传入）：
//   - 注入 ctx 后，工具函数从 ctx 取身份做归属校验（LLM 参数里永远没有 user 概念）
//   - clientFor 据此决定用用户的 BYOK key 还是服务器默认 key
// func (as *AgentService) StreamChat(ctx context.Context, userID uint, userContent, deckID string, emit func(StreamEvent) error) error {
// 	ctx = authctx.WithUser(ctx, userID)

// 	client := as.clientFor(ctx)

// 	messages, err := as.buildMessages(userContent, deckID)
// 	if err != nil {
// 		return err
// 	}

// 	opt := as.setChatOpts()
// 	opt.Messages = messages
// 	opt.StreamOptions = openai.ChatCompletionStreamOptionsParam{
// 		IncludeUsage: openai.Bool(true),
// 	}
// 	opt.Tools = as.Tools

// 	var fullText strings.Builder

// 	var(
// 		promptTokens int64 = 0		//输入token
// 		completionTokens int64 =0   //输出token
// 		totalTokens int64 =0		//总token
// 		cachedTokens int64 =0		//缓存命中token
// 	)

// 	for i:=0;;i++{
// 		if i>maxTurns{
// 			opt.Messages = append(opt.Messages, openai.UserMessage("工具调用预算已用完：不要再调用任何工具，直接基于以上获取的信息给出最终回答"))
// 			opt.Tools = nil
// 			stream := client.Chat.Completions.NewStreaming(ctx,opt)
// 			acc := openai.ChatCompletionAccumulator{}
// 			for stream.Next(){
// 				chunk := stream.Current()
// 				if !acc.AddChunk(chunk){
// 					continue
// 				}
// 				if len(chunk.Choices) == 0{
// 					continue
// 				}
// 				delta := chunk.Choices[0].Delta

// 				if delta.Content!=""{
// 					fullText.WriteString(delta.Content)
// 					if emitErr := emit(StreamEvent{Type: EventTypeDelta,Content: delta.Content});emitErr!=nil{
// 						stream.Close()
// 						return fmt.Errorf("事件推送失败 :%w",emitErr)
// 					}
// 				}

// 				if f,ok := delta.JSON.ExtraFields["reasoning_content"];ok{
// 					var reasoning string
// 					err := json.Unmarshal([]byte(f.Raw()),&reasoning);
// 					if err==nil && reasoning!="" {
// 						if emitErr := emit(StreamEvent{Type:EventTypeThink,Content: reasoning});emitErr!=nil{
// 							stream.Close()
// 							return fmt.Errorf("事件推送失败 :%w",emitErr)
// 						}
// 					}
// 					if err !=nil{
// 						stream.Close()
// 						return fmt.Errorf("流式响应不完整 err:%w",err)
// 					}
// 				}

// 			}
// 			if err := stream.Err(); err !=nil{
// 				if emitErr := emit(StreamEvent{Type: EventTypeError,Content: err.Error()});emitErr!=nil{
// 					stream.Close()
// 					return fmt.Errorf("事件推送失败 :%w",emitErr)
// 				}
// 				stream.Close()
// 				return fmt.Errorf("流式输出错误 err:%w",err)
// 			}
// 			stream.Close()

// 			if len(acc.Choices) == 0{
// 				if emitErr := emit(StreamEvent{Type:EventTypeError,Content: "流式响应不完整"});emitErr!=nil{
// 					return fmt.Errorf("事件推送失败 :%w",emitErr)
// 				}
// 				return fmt.Errorf("流式响应不完整")
// 			}

// 			promptTokens+= acc.Usage.PromptTokens
// 			completionTokens += acc.Usage.CompletionTokens
// 			totalTokens += acc.Usage.TotalTokens
// 			cachedTokens += acc.Usage.PromptTokensDetails.CachedTokens

// 			if emitErr := emit(StreamEvent{Type:EventTypeDone,Content: fullText.String(),PromptTokens: promptTokens,CompletionTokens: completionTokens,TotalTokens: totalTokens,CachedTokens: cachedTokens});emitErr!=nil{
// 				return fmt.Errorf("事件推送失败 :%w",emitErr)
// 			}
// 			return nil
// 		}

// 		stream := client.Chat.Completions.NewStreaming(ctx,opt)
// 		acc := openai.ChatCompletionAccumulator{}
// 		for stream.Next(){
// 			chunk := stream.Current()
// 			if !acc.AddChunk(chunk){
// 				continue
// 			}
// 			if len(chunk.Choices)==0{
// 				continue
// 			}
// 			delta := chunk.Choices[0].Delta
// 			if delta.Content!=""{
// 				fullText.WriteString(delta.Content)
// 				if emitErr := emit(StreamEvent{Type:EventTypeDelta,Content: delta.Content});emitErr!=nil{
// 					stream.Close()
// 					return fmt.Errorf("事件推送失败 :%w",emitErr)
// 				}
// 			}
// 			if f,ok := delta.JSON.ExtraFields["reasoning_content"];ok{
// 				var reasoning string
// 				err := json.Unmarshal([]byte(f.Raw()),&reasoning);
// 				if err==nil && reasoning!="" {
// 					if emitErr := emit(StreamEvent{Type:EventTypeThink,Content: reasoning});emitErr!=nil{
// 						stream.Close()
// 						return fmt.Errorf("事件推送失败 :%w",emitErr)
// 					}
// 				}
// 				if err !=nil{
// 					stream.Close()
// 					return fmt.Errorf("流式响应不完整 err:%w",err)
// 				}
// 			}
// 		}
// 		if err := stream.Err(); err !=nil{
// 			if emitErr := emit(StreamEvent{Type: EventTypeError,Content: err.Error()});emitErr!=nil{
// 				stream.Close()
// 				return fmt.Errorf("事件推送失败 :%w",emitErr)
// 			}
// 			stream.Close()
// 			return fmt.Errorf("流式输出错误 err:%w",err)
// 		}
// 		stream.Close()

// 		if len(acc.Choices) == 0 || acc.Choices[0].FinishReason == ""{
// 			if emitErr := emit(StreamEvent{Type: EventTypeError,Content: "流式响应不完整，请重试，本轮不计入上下文"});emitErr!=nil{
// 				return fmt.Errorf("事件推送失败 :%w",emitErr)
// 			}
// 			return fmt.Errorf("流式响应不完整，请重试，本轮不计入上下文")
// 		}

// 		//计算token
// 		choice := acc.Choices[0]
// 		promptTokens+= acc.Usage.PromptTokens
// 		completionTokens += acc.Usage.CompletionTokens
// 		totalTokens += acc.Usage.TotalTokens
// 		cachedTokens += acc.Usage.PromptTokensDetails.CachedTokens

// 		//拿到拼接好的回答
// 		msg := choice.Message
// 		// 没有工具调用，回答结束
// 		if len(msg.ToolCalls)==0{
// 			if emitErr := emit(StreamEvent{Type:EventTypeDone,Content: fullText.String(),PromptTokens: promptTokens,CompletionTokens: completionTokens,TotalTokens: totalTokens,CachedTokens: cachedTokens});emitErr!=nil{
// 				return fmt.Errorf("事件推送失败 :%w",emitErr)
// 			}
// 			return nil
// 		}

// 		//拼接上下文
// 		opt.Messages = append(opt.Messages, msg.ToParam())

// 		//执行工具
// 		for _,tool := range msg.ToolCalls{
// 			var result string
// 			t,ok := as.Exec[tool.Function.Name]
// 			if !ok{
// 				result = "未知工具"
// 				if emitErr := emit(StreamEvent{Type:EventTypeToolError,Content: result,ToolName: tool.Function.Name});emitErr!=nil{
// 					return fmt.Errorf("事件推送失败 :%w",emitErr)
// 				}
// 			}else{
// 				res ,err := t(ctx,tool.Function.Arguments)
// 				if err!=nil{
// 					result = fmt.Sprintf("工具调用失败 err:%s",err.Error())
// 					if emitErr := emit(StreamEvent{Type:EventTypeToolError,Content: result,ToolName: tool.Function.Name});emitErr!=nil{
// 						return fmt.Errorf("事件推送失败 :%w",emitErr)
// 					}
// 				}else{
// 					result = res
// 					if emitErr := emit(StreamEvent{Type:EventTypeToolCall,Content: result,ToolName: tool.Function.Name});emitErr!=nil{
// 						return fmt.Errorf("事件推送失败 :%w",emitErr)
// 					}
// 				}
// 			}
// 			opt.Messages = append(opt.Messages, openai.ToolMessage(result,tool.ID))
// 		}
// 		//继续下一轮推理
// 	}
// }

// weekdayCN 给日期配一个中文星期：不配的话「上周末」「这周三」这类说法模型算不出来。
// 数组下标直接就是 time.Weekday（周日=0），不需要转换。
var weekdayCN = [...]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}

func (as *AgentService) buildSystemMessage(deckID string) string {
	msg := systemPrompt
	if deckID != "" && deck.IsValidID(deckID) {
		msg += fmt.Sprintf(`当前用户正在预览的演示文稿是 {"deck_id":"%s"},涉及它的修改直接使用这个 deck_id, 不要向用户询问`, deckID)
	}

	// 当前日期追加在**整个系统提示词的最后**，这个位置是刻意的，不要往前挪：
	// 前缀缓存要求"完整匹配一个已持久化的前缀单元"，所以变化的内容越靠后，被打断的部分越少。
	// 实测（约 2000 token 的静态前缀 + 末尾日期）：把末尾日期改一天，cached_tokens 从
	// 2560/2696 只掉到 2432/2696——前面那 2560 token 照样命中。反过来把日期放开头，
	// 就是第一个 token 就分叉、后面全部重算。
	//
	// 这也是它不写进 systemPrompt.md 的原因：那个文件是 //go:embed 编译进二进制的，
	// 写死一个日期等于"发布即过期"，而且没有任何东西会报出来。
	//
	// 最后那句"不要为了确认日期去联网搜索"不是客套：没有明确授权时模型会在 thinking 里
	// 纠结"我的知识截止 2024-06、不能访问实时网络"，实测它会真的发起一次计费搜索去问今天几号。
	now := time.Now()
	msg += fmt.Sprintf("\n\n当前日期：%s（%s）。涉及「今天」「本月」「最近」这类时间说法时以它为准——"+
		"你的训练数据有截止时间，不要按它推断当前时间，也不要为了确认日期去联网搜索。",
		now.Format("2006-01-02"), weekdayCN[int(now.Weekday())])

	return msg
}

func (as *AgentService) NewStreamChat(ctx context.Context, userID, sessionID uint, userContent, deckID string, emit func(StreamEvent) error) (uint, error) {
	ctx = authctx.WithUser(ctx, userID)
	client := as.clientFor(ctx)

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

	//首轮构建系统提示词
	if sessionID == 0 {
		messages = append(messages, openai.SystemMessage(as.buildSystemMessage(deckID)))
	}

	messages = append(messages, openai.UserMessage(userContent))
	if err := as.persistSession(sess, messages); err != nil {
		return sess.ID, err
	}
	if emitErr := emit(StreamEvent{Type: EventTypeSession, Content: strconv.FormatUint(uint64(sess.ID), 10)}); emitErr != nil {
		return sess.ID, fmt.Errorf("事件推送失败: %w", emitErr)
	}

	paused, err := as.runLoop(ctx, client, sess, messages, emit, runTraceInfo{
		UserID:      userID,
		DeckID:      deckID,
		UserContent: userContent,
	})
	if paused {
		return sess.ID, ErrPaused
	}
	if err != nil {
		return sess.ID, err
	}
	return sess.ID, nil
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
	// parent_run_id 指回被暂停的那次，观测页据此把被打断的对话串起来看
	paused, err := as.runLoop(ctx, client, sess, messages, emit, runTraceInfo{
		UserID:      userID,
		DeckID:      sess.DeckID,
		UserContent: "(用户回答了 ask_user 的提问)",
		ParentRunID: pending.RunID,
	})
	if paused {
		return sess.ID, ErrPaused
	}
	if err != nil {
		return sess.ID, err
	}
	return sess.ID, nil
}

func (as *AgentService) runLoop(ctx context.Context, client *openai.Client, sess *store.ChatSession, messages []openai.ChatCompletionMessageParamUnion, emit func(StreamEvent) error, info runTraceInfo) (paused bool, err error) {
	rr := newRunRecorder()
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
	opt.Tools = as.Tools

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
		//预算花完了
		if i > maxTurns {
			messages = append(messages, openai.UserMessage("工具调用预算已用完：不要再调用任何工具，直接基于以上获取的信息给出最终回答"))
			opt.Messages = messages
			opt.Tools = nil
			err := as.persistSession(sess, messages)
			if err != nil {
				return false, err
			}
		}
		msg, usage, streamErr := as.streamOnce(turnCtx, client, opt, &fullText, emit)
		if streamErr != nil {
			return false, streamErr
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
				return false, err
			}
			// 汇总带上分项：扁平那四个字段保持"主循环口径"不动（前端契约），
			// 真实总成本看 Usage.Total
			summary := rec.Summary()
			if emitErr := emit(StreamEvent{Type: EventTypeDone, Content: fullText.String(), PromptTokens: promptTokens, CompletionTokens: completionTokens, TotalTokens: totalTokens, CachedTokens: cachedTokens, Usage: &summary}); emitErr != nil {
				return false, fmt.Errorf("事件推送失败: %w", emitErr)
			}
			endRun(trace.StatusOK)
			return false, nil
		}

		messages = append(messages, msg.ToParam())
		opt.Messages = messages
		if err := as.persistSession(sess, messages); err != nil {
			return false, err
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
							return false, err
						}
						continue
					}
					//如果提问的问题大于6个小于1个 同上
					result := "ask_user 参数不合法: questions 必须是 1~6 个问题"
					trace.Emit(askCtx, trace.Event{Kind: trace.KindToolResult, Result: result, Error: "questions 数量不在 1~6"})
					messages = append(messages, openai.ToolMessage(result, tool.ID))
					opt.Messages = messages
					if err := as.persistSession(sess, messages); err != nil {
						return false, err
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
					return false, err
				}
				if emitErr := emit(StreamEvent{Type: EventTypeAskUser, ToolCallID: tool.ID, Content: tool.Function.Arguments}); emitErr != nil {
					return false, fmt.Errorf("事件推送失败: %w", emitErr)
				}
				trace.Emit(askCtx, trace.Event{Kind: trace.KindToolResult, Result: `{"note":"已向用户提问，本轮暂停等回答"}`})
				endRun(trace.StatusPaused)
				return true, nil
			}
			//普通工具调用
			result, execErr := as.execTool(turnCtx, tool, emit, rr)
			if execErr != nil {
				return false, execErr
			}
			messages = append(messages, openai.ToolMessage(result, tool.ID))
			opt.Messages = messages
			if err := as.persistSession(sess, messages); err != nil {
				return false, err
			}
		}
		// 继续下一轮
	}
}

func (as *AgentService) execTool(ctx context.Context, tool openai.ChatCompletionMessageToolCallUnion, emit func(StreamEvent) error, rr *runRecorder) (string, error) {
	// 把归属挂进 ctx：工具内部（联网搜索的每一次子搜索、视觉审查的每一页截图）
	// 因此不必知道自己在第几轮、call_id 是什么，只要 trace.Emit 就会自动归位
	ctx = trace.WithTool(ctx, tool.Function.Name, tool.ID)
	trace.Emit(ctx, trace.Event{Kind: trace.KindToolCall, Args: tool.Function.Arguments})

	var (
		result    string
		callErr   string // 空串 = 没失败。刻意不用 error 类型：它同时要进 JSON
		startedAt = time.Now()
	)

	t, ok := as.Exec[tool.Function.Name]
	switch {
	case !ok:
		result = "未知工具"
		callErr = result
	default:
		res, err := t(ctx, tool.Function.Arguments)
		if err != nil {
			result = fmt.Sprintf("工具调用失败 err:%s", err.Error())
			callErr = err.Error()
		} else {
			result = res
			rr.noteToolCall(tool.Function.Name, tool.Function.Arguments, res) //版本备注只记成功的
		}
	}

	// 失败也要记。原来 runRecorder 只记成功调用，于是"模型连试三次都被拒"这种
	// 最需要解释的现象，在记录里完全看不到——你只看得到它最后放弃了。
	trace.Emit(ctx, trace.Event{
		Kind: trace.KindToolResult, Result: result, Error: callErr,
		DurationMS: time.Since(startedAt).Milliseconds(),
	})

	if callErr != "" {
		if emitErr := emit(StreamEvent{Type: EventTypeToolError, Content: result, ToolName: tool.Function.Name}); emitErr != nil {
			return result, fmt.Errorf("事件推送失败: %w", emitErr)
		}
	} else {
		if emitErr := emit(StreamEvent{Type: EventTypeToolCall, Content: result, ToolName: tool.Function.Name}); emitErr != nil {
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
		ReasoningEffort: shared.ReasoningEffortMax,
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
