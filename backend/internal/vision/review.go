package vision

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

const reviewPrompt = `你是这套幻灯片框架的版面审查员。下面按页给你：一页渲染截图 + 这一页的量测数字。

量测字段：fit=适配兜底缩放系数（1.00 表示没被缩过）；最小字号（px）；溢出比（>1.02 = 内容超出画布）；
子元素 = section 直接子元素个数；字数 = 整页可见文字数。

程序已经做完的硬判定列在最前面，**不要复核这些数字**。你的任务只有两件：
1. 用眼睛确认硬判定在画面上真的看得出来，并说清是哪一处（引用页面上的实际文字）。
2. 找出数字看不出来的问题：文字被裁、元素重叠、内容贴边、明显不对齐、
   对比度不足、装饰抢戏、以及这一页的图和标题讲的不是一回事。

不要提"建议补充数据""内容可以更丰富"这类内容意见，只说版面。不要重写内容，只指出问题。

输出格式（严格照做）：
第 N 页｜问题类别｜具体是什么（引用页面实际文字）｜改哪里
每页最多一行，只写有问题的页。最后一行为：合计：N 页有问题
都没问题则只输出一行：未发现问题`

func Review(ctx context.Context,client *openai.Client,model string,d *Deck)(string,error){
	var parts []openai.ChatCompletionContentPartUnionParam
	
	head := reviewPrompt + "\n\n程序已判定的问题：\n"
	if f := HardFindings(d); len(f)>0{
		head += strings.Join(f,"\n")
	}else {
		head += "无(所有硬指标都在范围内)"
	}
	head += "\n\n逐页量测：\n"+Digest(d)
	parts = append(parts, openai.TextContentPart(head))

	for _,s := range d.Slides {
		parts = append(parts, openai.TextContentPart(fmt.Sprintf("\n第 %d 页：",s.Index+1)))
		parts = append(parts,openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL:"data:image/png;base64,"+base64.StdEncoding.EncodeToString(s.PNG),
			Detail: "auto",
		}))
	}

	resp,err := client.Chat.Completions.New(ctx,openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		ReasoningEffort: shared.ReasoningEffortHigh,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(parts),
		},
	})

	if err != nil{
		return "",err
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content),nil
}