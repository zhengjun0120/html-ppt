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

// PendingAsk GET /api/chat/pending?session_id=N
//
// 返回暂停中那条提问的原始参数，让页面在**刷新之后**把提问卡片重建出来。
//
// 为什么必须有它：agent 提问时是暂停态，而暂停期间后端只接受"回答"
//（新消息会被 agent.guardNewMessage 拦下）。前端的状态却活不过一次刷新——
// 卡片丢了，用户既答不上（不知道问的是什么）又发不出新消息，
// 这个会话就只能弃掉。有了它，刷新后卡片能原样回来，硬拦截才不构成死路。
//
// 没有暂停时 questions 返回空串（不是错误）：那是正常状态。
func (h *Handler) PendingAsk(c *gin.Context){
	uid,ok := authctx.UserID(c.Request.Context())
	if !ok{
		response.Err(c,http.StatusUnauthorized,"未登录")
		return
	}

	sessionID,err := queryUint(c,"session_id")
	if err !=nil{
		response.Err(c,400,"session_id 必须是数字")
		return
	}
	if sessionID == 0{
		response.Err(c,400,"缺少 session_id")
		return
	}

	if err := h.agent.EnsureSessionOwner(uid,sessionID);err !=nil{
		response.Err(c,http.StatusNotFound,"会话不存在")
		return
	}

	questions,err := h.agent.PendingAskQuestions(c.Request.Context(),uid,sessionID)
	if err != nil{
		// 会话存在但读不出提问（消息损坏之类）：不该 500 让页面卡住，
		// 返回"没有待答提问"即可——真有问题用户发消息时会拿到明确报错
		response.OK(c,gin.H{"questions":""})
		return
	}
	response.OK(c,gin.H{"questions":questions})
}