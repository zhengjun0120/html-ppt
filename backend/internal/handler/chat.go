package handler

import (
	"encoding/json"
	"fmt"
	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/response"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

// Chat POST /api/chat —— agent 循环的入口，阶段1实现。
//
// 计划中的协议（SSE 流式返回，事件类型）：
//   agent_text    LLM 的流式文字回复
//   tool_progress 工具执行进度（如 "正在读取第 3 页…"）
//   ask_user      结构化提问卡片（问题 + 推荐答案 + 选项），循环暂停等用户作答
//   done          本轮结束
type ChatRequest struct {
	UserContent string 	`json:"user_content"`
	EnableWebSearch bool `json:"enable_web_search,omitempty"` //预留是否开启联网搜索
	DeckID string `json:"deck_id,omitempty"`
}

func (h *Handler) Chat(c *gin.Context) {

	var req ChatRequest
	if err:= c.ShouldBindJSON(&req);err !=nil{
		response.ParameterErr(c)
		return
	}

	if req.UserContent == ""{
		response.Err(c,400,"用户消息不可为空")
		return
	}

	//设置sse响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	ch := make(chan agent.StreamEvent,16)

	go func(){
		defer close(ch)
		
		defer func(){
			if r := recover();r!=nil{
				log.Printf("agent 对话触发panic")
				select{
				case ch<- agent.StreamEvent{Type: agent.EventTypeError,Content: "服务器内部错误，请稍后重试"}:
				case <- c.Request.Context().Done():
				}
			}
		} ()

		err := h.agent.StreamChat(c.Request.Context(),req.UserContent,req.DeckID,func(ev agent.StreamEvent) error{
			select{
			case ch<-ev:
				return nil
			case <-c.Request.Context().Done():
				return c.Request.Context().Err()
			}
		})

		if err !=nil{
			log.Printf("agent流式对话失败 err:%v", err)
			select{
			case ch<- agent.StreamEvent{Type: agent.EventTypeError,Content: "对话服务暂时不可用，请稍后重试"}:
			case <- c.Request.Context().Done():
			}
		}
	}()

	// 流式响应
	c.Stream(func(w io.Writer) bool {
		ev,ok := <- ch
		if !ok{
			return false
		}
		data,_:= json.Marshal(ev)
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, data)
		return true
	})
}
