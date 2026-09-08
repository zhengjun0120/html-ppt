package handler

import (
	"encoding/json"
	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnswerItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type AnswerRequest struct {
	SessionID uint         `json:"session_id"`
	Answers   []AnswerItem `json:"answers"`
	Note      string       `json:"note"` //用户拒绝回答则为空，如果有新的指示则使用这个字段，如 用户说“你来决定吧”并跳过问题
}

func (h *Handler) AskUser(c *gin.Context){
	uid,ok := authctx.UserID(c.Request.Context())
	if !ok{
		response.Err(c,http.StatusUnauthorized,"未登录")
		return
	}

	var req AnswerRequest
	if err := c.ShouldBindJSON(&req);err !=nil{
		response.ParameterErr(c)
		return
	}

	if req.SessionID==0{
		response.Err(c,400,"缺少 session_id")
		return
	}

	if err := h.agent.EnsureSessionOwner(uid,req.SessionID);err !=nil{
		response.Err(c,http.StatusNotFound,"会话不存在")
		return
	}

	var payload any
	if len(req.Answers) >0{
		payload = map[string]any{"answers":req.Answers}
	}else{
		note := req.Note
		if note == ""{
			note = "用户未作答，请使用推荐答案"
		}
		payload = map[string]any{"note":note}
	}
	data,err := json.Marshal(payload)
	if err !=nil{
		response.Err(c,500,"回答序列化json失败")
		return
	}

	serveAgentSSE(c,func(emit func(agent.StreamEvent) error) (uint, error) {
		return h.agent.AnswerChat(c.Request.Context(),uid,req.SessionID,string(data),emit)
	})
}