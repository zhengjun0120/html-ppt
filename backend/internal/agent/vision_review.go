package agent

import (
	"context"
	"html-ppt/backend/internal/vision"
	"log"
	"strings"
)

const reviewChars = 1500 //限制返回给模型的报告字数长度

//运行视觉审查
func (a *AgentService) runVisionReview(ctx context.Context,uid uint,deckID string)string{
	if !a.Vision || a.VisionGrants == nil{
		return ""
	} 
	//获取门票
	nonce,err := a.VisionGrants.Issue(uid,deckID)
	if err != nil{
		log.Printf("[warn] 视觉审查：发票据失败,跳过 err：%v",err)
		return ""
	}
	//预览路由
	url := strings.TrimRight(a.VisionBaseURL,"/")+"/render/"+nonce

	//拿到截图和量测结果
	deck,err := vision.Capture(ctx,vision.Options{URL: url,ChromePath: a.ChromePath})
	if err != nil{
		log.Printf("[warn] 视觉审查：渲染/量测失败，跳过 err：%v",err)
		return ""
	}

	var b strings.Builder
	b.WriteString("\n\n视觉审查（无头浏览器实拍 + 量测）：\n")
	b.WriteString(deck.Text())

	client := a.clientFor(ctx)
	report,err := vision.Review(ctx,client,a.ModelID,deck)
	if err != nil{
		log.Printf("[warn] 视觉审查：看图失败，只回量测 err:%v",err)
		b.WriteString("\n (看图部分失败，以上是程序量测结果)")
		return clipRunes(b.String(),reviewChars) 
	}

	b.WriteString("\n看图审查：\n")
	b.WriteString(report)
	return clipRunes(b.String(),reviewChars)
}

func clipRunes(s string,n int) string{
	r :=[]rune(s)
	if len(r) <=n{
		return s
	}
	return string(r[:n])+"\n(报告过长已截断)"
}