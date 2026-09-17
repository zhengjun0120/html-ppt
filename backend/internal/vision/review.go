package vision

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

const reviewPrompt = `你是这套幻灯片框架的版面审查员。下面按页给你：一页渲染截图 + 这一页的量测数字。

量测字段：fit=适配兜底缩放系数（1.00 表示没被缩过）；最小字号（px）；溢出比（>1.02 = 内容超出画布）；
子元素 = section 直接子元素个数；字数 = 整页可见文字数。

**只看给你的这几页**：没给你截图的页不要在报告里出现（你可能在硬判定里看到它的页号，那是整份 deck
层面的判定，不要替它写结论）。**也不要复核程序给出的数字**。你的任务只有两件：
1. 用眼睛确认硬判定在画面上真的看得出来，并说清是哪一处（引用页面上的实际文字）。
2. 找出数字看不出来的问题：文字被裁、元素重叠、内容贴边、明显不对齐、对比度不足、装饰抢戏、以及这一页的图和标题讲的不是一回事。

不要提"建议补充数据""内容可以更丰富"这类内容意见，只说版面。不要重写内容，只指出问题。

**两条纪律，违反的代价是主模型陷入"改→再审→改"的无限返工：**
- **每条问题分级**：【硬】= 投屏会翻车的：文字被裁、元素重叠、看不清、明显错位、断行尴尬。
  【软】= 审美打磨：对比度略闷、间距略宽、差几像素的对齐。拿不准归哪类时算【软】。
- **没把握就不报**：一页没问题就写没问题，不要为了凑数硬找。每一条软意见都可能换来一轮
  改版面，而每次改版面又可能带出新的问题——"报得准"远比"报得多"有价值。

输出格式（严格照做）：
第 N 页｜【硬或软】｜问题类别｜具体是什么（引用页面实际文字）｜期望效果
每页最多一行，只写有问题的页。最后一行为：合计：N 页有问题（硬 M 页）
都没问题则只输出一行：未发现问题

"期望效果"只写**改成什么样**，不要指挥实现手段——这套组件库有你看不到的约束
（比如禁止内联改字号、颜色只能用主题变量），指定手段经常做不到，还会把主模型带偏：
写"让该标签与标题同级可读"，不要写"把字号提到 31px"。`

// ReviewResult 看图的产出。把"发了什么提示词、附了几张图、花了多少"一并交出来，
// 是因为这三样都是观测要用的：审查报告的每个结论都要能回溯到"它当时看到了什么"，
// 而看图是整个 agent 里最贵的一次调用（几张 1244×700 的图，图片 token 很值钱）。
type ReviewResult struct {
	Report string                 // 给主模型看的审查报告
	Prompt string                 // 实际发出去的文本部分（图片作为独立 part 追加在后面）
	Images int                    // 附了几张图
	Usage  openai.CompletionUsage // 这次调用的用量；过去拿到了就扔，于是 token 统计系统性偏小
}

// Review 让模型看 sel 里这几页的截图。
//
// findings **必须由调用方从整份 deck 上算出来再传进来**（vision.ScopeFindings），
// 不能拿 sel 现算：版式重复是 deck 级规则，"同一版式用了 3 次"在只选了 2 页的子集里
// 永远数不出来，而模型会照着"没有版式重复"去确认画面——它确认的是一个假结论。
func Review(ctx context.Context, client *openai.Client, model string, sel *Deck, findings []string,onDelta func(string)) (*ReviewResult, error) {


	if onDelta == nil{
		onDelta = func(string){}
	}

	//拼接提示词
	var parts []openai.ChatCompletionContentPartUnionParam

	head := reviewPrompt + fmt.Sprintf("\n\n你要看的页：第 %s 页（共 %d 张截图）\n", pageList(sel), len(sel.Slides))
	//把有问题的页输出
	head += "\n程序已判定的问题：\n"
	if len(findings) > 0 {
		head += strings.Join(findings, "\n")
	} else {
		head += "无(所有硬指标都在范围内)"
	}
	//输出每页量测的结果
	head += "\n\n逐页量测：\n" + Digest(sel)
	parts = append(parts, openai.TextContentPart(head))

	//拼接页和该页截图
	for _, s := range sel.Slides {
		parts = append(parts, openai.TextContentPart(fmt.Sprintf("\n第 %d 页：", s.Index+1)))
		parts = append(parts, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL:    "data:image/png;base64," + base64.StdEncoding.EncodeToString(s.PNG),
			Detail: "auto",
		}))
	}
	// 空报告重试一次。推理型模型偶尔会把整个 completion 预算全部花在 reasoning 上、
	// 一字正文不吐（实测 trace 1789613158323-30a4：completion=8000 全是 reasoning），
	// 这是概率行为不是稳定行为——同样的输入再发一次大多能出正文。截图的编码、
	// 提示词的拼接都只做一次，重试只重发请求；两次的用量都算这次审查的花费。
	tryReview := func() (string, openai.CompletionUsage, error) {
		stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Model: openai.ChatModel(model),
			// Low 而不是 High：这个调用是同步阻塞的（以前阻塞在 write_deck 里，现在阻塞在
			// review_slides 里）。实测 High 要约 48 秒，而 Low 明显更快、同一份 deck
			// 抓到的是同一批问题（孤字换行、三列不齐、字号偏小）。
			ReasoningEffort: shared.ReasoningEffortLow,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(parts),
			},
			// 显式给预算：reasoning 先吃 token、吃完就没有 content，而"content 为空"不会报错
			// ——上层会照样写一行"看图审查："然后什么都没有（我第一次探针就是 max_tokens=200
			// 撞出来的：content 空、completion 却用满 200）。实测 9 页用掉 4382；
			// 现在单次最多 6 页，8000 的余量更宽。
			MaxTokens: openai.Int(8000),

			StreamOptions: openai.ChatCompletionStreamOptionsParam{IncludeUsage: openai.Bool(true)},
		})

		acc := openai.ChatCompletionAccumulator{}
		var sb strings.Builder
		for stream.Next() {
			chunk := stream.Current()
			if !acc.AddChunk(chunk) {
				continue
			}
			if len(chunk.Choices) == 0 {
				continue
			}
			delta := chunk.Choices[0].Delta
			if delta.Content != "" {
				sb.WriteString(delta.Content)
				onDelta(delta.Content)
			}
		}

		if err := stream.Err(); err != nil {
			stream.Close()
			return "", openai.CompletionUsage{}, err
		}
		stream.Close()
		if len(acc.Choices) == 0 || acc.Choices[0].FinishReason == "" {
			return "", openai.CompletionUsage{}, errors.New("看图审查流式响应不完整（没有 finish_reason）")
		}
		return strings.TrimSpace(sb.String()), acc.Usage, nil
	}

	report, usage, err := tryReview()
	if err == nil && report == "" {
		log.Printf("[warn] 视觉审查：模型返回空报告（completion=%d 全为 reasoning），重试一次",
			usage.CompletionTokens)
		report, usage, err = tryReview()
	}
	if err != nil {
		return nil, err
	}
	if report == "" {
		// 空报告不能当成功返回：否则上层会留一段空的"看图审查："，
		// 看起来像"审过了、没问题"，实际是这次调用什么都没产出。
		// 重试过仍为空才走到这里——别再无限重试，把决定权交回上层
		//（它会把量测数字还给模型并退还配额，模型可以直接重试 review_slides）。
		return nil, fmt.Errorf("模型返回了空报告（completion=%d，其中 reasoning=%d），重试一次仍为空",
			usage.CompletionTokens, usage.CompletionTokensDetails.ReasoningTokens)
	}
	return &ReviewResult{
		Report: report,
		Prompt: head,
		Images: len(sel.Slides),
		Usage:  usage,
	}, nil
}

func pageList(d *Deck) string {
	ps := make([]string, 0, len(d.Slides))
	for _, s := range d.Slides {
		ps = append(ps, strconv.Itoa(s.Index+1))
	}
	return strings.Join(ps, "、")
}
