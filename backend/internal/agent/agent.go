package agent

import (
	"context"
	"encoding/json"
	"html-ppt/backend/internal/service/deck"

	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

const maxTurns = 20 //最大允许调用20轮llm请求

func (as *AgentService) buildMessages(userContent,deckID string) ([]openai.ChatCompletionMessageParamUnion,error){
	//TODO 从数据库中获取对话记录，如果数据库中没有 则先拼接系统提示词 并新建数据



	//TODO 然后拼接用户对话

	// Test
	messages := make([]openai.ChatCompletionMessageParamUnion,0)
	systemMessage := "你是一个专业设计ppt的Ai，你需要帮助用户使用reveal.js构建一个html的演示文稿"
	if deckID != "" && deck.IsValidID(deckID){
		systemMessage += fmt.Sprintf(`当前用户正在预览的演示文稿是 {"deck_id":"%s"}，涉及它的修改直接用这个 deck_id，不要向用户询问`,deckID)
	}
	messages = append(messages, openai.SystemMessage(systemMessage))
	messages = append(messages, openai.UserMessage(userContent))

	return messages,nil
}



func (as *AgentService) StreamChat(ctx context.Context, userContent,deckID string, emit func(StreamEvent) error) error {

	messages,err := as.buildMessages(userContent,deckID)
	if err !=nil{
		return err
	}
	
	opt := as.setChatOpts()
	opt.Messages = messages
	opt.StreamOptions = openai.ChatCompletionStreamOptionsParam{
		IncludeUsage: openai.Bool(true),
	}
	opt.Tools = as.Tools

	var fullText strings.Builder

	var(
		promptTokens int64 = 0		//输入token
		completionTokens int64 =0   //输出token
		totalTokens int64 =0		//总token
		cachedTokens int64 =0		//缓存命中token
	)

	for i:=0;;i++{
		if i>maxTurns{
			opt.Messages = append(opt.Messages, openai.UserMessage("工具调用预算已用完：不要再调用任何工具，直接基于以上获取的信息给出最终回答"))
			opt.Tools = nil
			stream := as.ModelClient.Chat.Completions.NewStreaming(ctx,opt)
			acc := openai.ChatCompletionAccumulator{}
			for stream.Next(){
				chunk := stream.Current()
				if !acc.AddChunk(chunk){
					continue
				}
				if len(chunk.Choices) == 0{
					continue
				}
				delta := chunk.Choices[0].Delta

				if delta.Content!=""{
					fullText.WriteString(delta.Content)
					if emitErr := emit(StreamEvent{Type: EventTypeDelta,Content: delta.Content});emitErr!=nil{
						stream.Close()
						return fmt.Errorf("事件推送失败 :%w",emitErr)
					}
				}

				
				if f,ok := delta.JSON.ExtraFields["reasoning_content"];ok{
					var reasoning string
					err := json.Unmarshal([]byte(f.Raw()),&reasoning); 
					if err==nil && reasoning!="" {
						if emitErr := emit(StreamEvent{Type:EventTypeThink,Content: reasoning});emitErr!=nil{
							stream.Close()
							return fmt.Errorf("事件推送失败 :%w",emitErr)
						}
					}
					if err !=nil{
						stream.Close()
						return fmt.Errorf("流式响应不完整 err:%w",err)
					}
				}

			}
			if err := stream.Err(); err !=nil{
				if emitErr := emit(StreamEvent{Type: EventTypeError,Content: err.Error()});emitErr!=nil{
					stream.Close()
					return fmt.Errorf("事件推送失败 :%w",emitErr)
				}
				stream.Close()
				return fmt.Errorf("流式输出错误 err:%w",err)
			}
			stream.Close()

			if len(acc.Choices) == 0{
				if emitErr := emit(StreamEvent{Type:EventTypeError,Content: "流式响应不完整"});emitErr!=nil{
					return fmt.Errorf("事件推送失败 :%w",emitErr)
				}
				return fmt.Errorf("流式响应不完整")
			}

			promptTokens+= acc.Usage.PromptTokens
			completionTokens += acc.Usage.CompletionTokens
			totalTokens += acc.Usage.TotalTokens
			cachedTokens += acc.Usage.PromptTokensDetails.CachedTokens

			if emitErr := emit(StreamEvent{Type:EventTypeDone,Content: fullText.String(),PromptTokens: promptTokens,CompletionTokens: completionTokens,TotalTokens: totalTokens,CachedTokens: cachedTokens});emitErr!=nil{
				return fmt.Errorf("事件推送失败 :%w",emitErr)
			}
			return nil
		}

		stream := as.ModelClient.Chat.Completions.NewStreaming(ctx,opt)
		acc := openai.ChatCompletionAccumulator{}
		for stream.Next(){
			chunk := stream.Current()
			if !acc.AddChunk(chunk){
				continue
			}
			if len(chunk.Choices)==0{
				continue
			}
			delta := chunk.Choices[0].Delta
			if delta.Content!=""{
				fullText.WriteString(delta.Content)
				if emitErr := emit(StreamEvent{Type:EventTypeDelta,Content: delta.Content});emitErr!=nil{
					stream.Close()
					return fmt.Errorf("事件推送失败 :%w",emitErr)
				}
			}
			if f,ok := delta.JSON.ExtraFields["reasoning_content"];ok{
				var reasoning string
				err := json.Unmarshal([]byte(f.Raw()),&reasoning); 
				if err==nil && reasoning!="" {
					if emitErr := emit(StreamEvent{Type:EventTypeThink,Content: reasoning});emitErr!=nil{
						stream.Close()
						return fmt.Errorf("事件推送失败 :%w",emitErr)
					}
				}
				if err !=nil{
					stream.Close()
					return fmt.Errorf("流式响应不完整 err:%w",err)
				}
			}
		}
		if err := stream.Err(); err !=nil{
			if emitErr := emit(StreamEvent{Type: EventTypeError,Content: err.Error()});emitErr!=nil{
				stream.Close()
				return fmt.Errorf("事件推送失败 :%w",emitErr)
			}
			stream.Close()
			return fmt.Errorf("流式输出错误 err:%w",err)
		}
		stream.Close()

		if len(acc.Choices) == 0 || acc.Choices[0].FinishReason == ""{
			if emitErr := emit(StreamEvent{Type: EventTypeError,Content: "流式响应不完整，请重试，本轮不计入上下文"});emitErr!=nil{
				return fmt.Errorf("事件推送失败 :%w",emitErr)
			}
			return fmt.Errorf("流式响应不完整，请重试，本轮不计入上下文")
		}

		//计算token
		choice := acc.Choices[0]
		promptTokens+= acc.Usage.PromptTokens
		completionTokens += acc.Usage.CompletionTokens
		totalTokens += acc.Usage.TotalTokens
		cachedTokens += acc.Usage.PromptTokensDetails.CachedTokens

		//拿到拼接好的回答
		msg := choice.Message
		// 没有工具调用，回答结束
		if len(msg.ToolCalls)==0{
			if emitErr := emit(StreamEvent{Type:EventTypeDone,Content: fullText.String(),PromptTokens: promptTokens,CompletionTokens: completionTokens,TotalTokens: totalTokens,CachedTokens: cachedTokens});emitErr!=nil{
				return fmt.Errorf("事件推送失败 :%w",emitErr)
			}
			return nil
		}

		//拼接上下文
		opt.Messages = append(opt.Messages, msg.ToParam())
		
		//执行工具
		for _,tool := range msg.ToolCalls{
			var result string
			t,ok := as.Exec[tool.Function.Name]
			if !ok{
				result = "未知工具"
				if emitErr := emit(StreamEvent{Type:EventTypeToolError,Content: result,ToolName: tool.Function.Name});emitErr!=nil{
					return fmt.Errorf("事件推送失败 :%w",emitErr)
				}
			}else{
				res ,err := t(ctx,tool.Function.Arguments)
				if err!=nil{
					result = fmt.Sprintf("工具调用失败 err:%s",err.Error())
					if emitErr := emit(StreamEvent{Type:EventTypeToolError,Content: result,ToolName: tool.Function.Name});emitErr!=nil{
						return fmt.Errorf("事件推送失败 :%w",emitErr)
					}
				}else{
					result = res
					if emitErr := emit(StreamEvent{Type:EventTypeToolCall,Content: result,ToolName: tool.Function.Name});emitErr!=nil{
						return fmt.Errorf("事件推送失败 :%w",emitErr)
					}
				}
			}
			opt.Messages = append(opt.Messages, openai.ToolMessage(result,tool.ID))
		}
		//继续下一轮推理
	}


	

	
}



func (as *AgentService) setChatOpts() openai.ChatCompletionNewParams{
	opt := openai.ChatCompletionNewParams{
		Model: openai.ChatModel(as.ModelID),
		ReasoningEffort: shared.ReasoningEffortMax,
	}
	return opt
}
