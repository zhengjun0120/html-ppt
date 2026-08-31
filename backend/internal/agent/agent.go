package agent

import (
	"context"
	"errors"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

func (as *AgentServer) TestChat(ctx context.Context,messages []openai.ChatCompletionMessageParamUnion)(string,error){
	res,err := as.ModelClient.Chat.Completions.New(ctx,openai.ChatCompletionNewParams{
		Model: openai.ChatModel(as.ModelID),
		Messages: messages,
		MaxCompletionTokens: openai.Int(as.MaxToken),
		ReasoningEffort: shared.ReasoningEffortMax,
	})
	if err !=nil{
		var apierr *openai.Error
		if errors.As(err,&apierr){
			return "",fmt.Errorf("api 错误 %w",err)
		}else{
			return "",fmt.Errorf("网络或其他错误，%w",err)
		}
	}

	return res.Choices[0].Message.Content,nil
}


func (as *AgentServer) StreamChat(ctx context.Context,messages []openai.ChatCompletionMessageParamUnion,emit func(StreamEvent))error{

}