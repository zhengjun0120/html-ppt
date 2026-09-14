package agent

import (
	"context"
	"encoding/json"
	"errors"
	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/store"
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

	paused, err := as.runLoop(ctx, client, sess, messages, emit)
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
	var pending struct {
		ToolCallID string `json:"tool_call_id"`
	}
	if err := json.Unmarshal([]byte(sess.PendingAsk), &pending); err != nil || pending.ToolCallID == "" {
		if err != nil {
			return sess.ID, fmt.Errorf("暂停态损坏 err:%w", err)
		} else {
			return sess.ID, fmt.Errorf("暂停态损坏 err: tool_call_id 为空")
		}
	}

	messages = append(messages, openai.ToolMessage(answersJSON, pending.ToolCallID))
	sess.PendingAsk = ""
	if err := as.persistSession(sess, messages); err != nil {
		return sess.ID, err
	}
	if emitErr := emit(StreamEvent{Type: EventTypeSession, Content: strconv.FormatUint(uint64(sess.ID), 10)}); emitErr != nil {
		return sess.ID, fmt.Errorf("事件推送失败 err:%w", emitErr)
	}

	paused, err := as.runLoop(ctx, client, sess, messages, emit)
	if paused {
		return sess.ID, ErrPaused
	}
	if err != nil {
		return sess.ID, err
	}
	return sess.ID, nil
}

func (as *AgentService) runLoop(ctx context.Context, client *openai.Client, sess *store.ChatSession, messages []openai.ChatCompletionMessageParamUnion, emit func(StreamEvent) error) (paused bool, err error) {
	rr := newRunRecorder()
	defer func(){
		if err == nil{
			as.recordRunVersions(ctx,rr)
		}
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
		msg, usage, streamErr := as.streamOnce(ctx, client, opt, &fullText, emit)
		if streamErr != nil {
			return false, streamErr
		}

		promptTokens += usage.PromptTokens
		completionTokens += usage.CompletionTokens
		totalTokens += usage.TotalTokens
		cachedTokens += usage.PromptTokensDetails.CachedTokens

		if len(msg.ToolCalls) == 0 {
			messages = append(messages, msg.ToParam())
			opt.Messages = messages
			err := as.persistSession(sess, messages)
			if err !=nil{
				return false,err
			}
			if emitErr := emit(StreamEvent{Type: EventTypeDone, Content: fullText.String(), PromptTokens: promptTokens, CompletionTokens: completionTokens, TotalTokens: totalTokens, CachedTokens: cachedTokens}); emitErr != nil {
				return false, fmt.Errorf("事件推送失败: %w", emitErr)
			}
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
				var args AskUserArgs
				//解析参数
				if err := json.Unmarshal([]byte(tool.Function.Arguments), &args); err != nil || len(args.Questions) == 0 || len(args.Questions) > 6 {
					//如果json不合法 把错误信息直接加入到上下文中 然后跳过
					if err != nil {
						messages = append(messages, openai.ToolMessage(fmt.Sprintf("ask_user 参数json不合法 err:%s", err.Error()), tool.ID))
						opt.Messages = messages
						if err := as.persistSession(sess, messages); err != nil {
							return false, err
						}
						continue
					}
					//如果提问的问题大于6个小于1个 同上
					messages = append(messages, openai.ToolMessage("ask_user 参数不合法: questions 必须是 1~6 个问题", tool.ID))
					opt.Messages = messages
					if err := as.persistSession(sess, messages); err != nil {
						return false, err
					}
					continue
				}
				//如果要调用提问用户 则不可调用其他工具
				for _, other := range msg.ToolCalls {
					if other.ID != tool.ID {
						messages = append(messages, openai.ToolMessage(`{"note":"用户被提问打断，本工具本轮未执行，用户回答后如仍需要请重新调用"}`, other.ID))
					}
				}
				opt.Messages = messages
				pending, _ := json.Marshal(map[string]string{"tool_call_id": tool.ID})
				sess.PendingAsk = string(pending)
				if err := as.persistSession(sess, messages); err != nil {
					return false, err
				}
				if emitErr := emit(StreamEvent{Type: EventTypeAskUser, ToolCallID: tool.ID, Content: tool.Function.Arguments}); emitErr != nil {
					return false, fmt.Errorf("事件推送失败: %w", emitErr)
				}
				return true, nil
			}
			//普通工具调用
			result, execErr := as.execTool(ctx, tool, emit,rr)
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

func (as *AgentService) execTool(ctx context.Context, tool openai.ChatCompletionMessageToolCallUnion, emit func(StreamEvent) error,rr *runRecorder) (string, error) {
	var result string
	if t, ok := as.Exec[tool.Function.Name]; !ok {
		result = "未知工具"
		if emitErr := emit(StreamEvent{Type: EventTypeToolError, Content: result, ToolName: tool.Function.Name}); emitErr != nil {
			return result, fmt.Errorf("事件推送失败: %w", emitErr)
		}
	} else if res, err := t(ctx, tool.Function.Arguments); err != nil {
		result = fmt.Sprintf("工具调用失败 err:%s", err.Error())
		if emitErr := emit(StreamEvent{Type: EventTypeToolError, Content: result, ToolName: tool.Function.Name}); emitErr != nil {
			return result, fmt.Errorf("事件推送失败: %w", emitErr)
		}
	} else {
		result = res
		rr.noteToolCall(tool.Function.Name,tool.Function.Arguments,res) //只记录成功的
		if emitErr := emit(StreamEvent{Type: EventTypeToolCall, Content: result, ToolName: tool.Function.Name}); emitErr != nil {
			return result, err
		}

		
	}

	return result, nil
}

func (as *AgentService) streamOnce(ctx context.Context, client *openai.Client, opt openai.ChatCompletionNewParams, fullText *strings.Builder, emit func(StreamEvent) error) (openai.ChatCompletionMessage, openai.CompletionUsage, error) {
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
		return openai.ChatCompletionMessage{}, openai.CompletionUsage{}, fmt.Errorf("流式输出错误: %w", err)
	}
	stream.Close()

	if len(acc.Choices) == 0 || acc.Choices[0].FinishReason == "" {
		return openai.ChatCompletionMessage{}, openai.CompletionUsage{}, errors.New("流式响应不完整，请重试")
	}
	return acc.Choices[0].Message, acc.Usage, nil
}

func (as *AgentService) setChatOpts() openai.ChatCompletionNewParams {
	opt := openai.ChatCompletionNewParams{
		Model:           openai.ChatModel(as.ModelID),
		ReasoningEffort: shared.ReasoningEffortMax,
	}
	return opt
}

func (as *AgentService) recordRunVersions(ctx context.Context,rr *runRecorder){
	uid,ok := authctx.UserID(ctx)
	if !ok {
		return
	}
	for _,deckID := range rr.deckIDs() {
			if err := as.DeckService.RecordRunVersion(uid,deckID,rr.detail(deckID));err != nil{
			log.Printf("[warn] deck %s 记录历史版本失败 err: %v",deckID,err)
		}		
	}
}
