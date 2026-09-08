package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
	"html-ppt/backend/internal/service/auth"
)

// registerRequest / loginRequest 共用一个绑定结构。
// 哨兵错误 → HTTP 状态码的映射全在这里，service 层只管领域语义。
type authRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
	APIKey   string `json:"api_key"`
}

// RequestCode POST /api/auth/code —— 发送注册验证码。
func (h *Handler) RequestCode(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParameterErr(c)
		return
	}
	if err := h.auth.RequestCode(c.Request.Context(), req.Email); err != nil {
		respondAuthErr(c, err)
		return
	}
	// 不回显"是否已注册"：这个接口本身不查库，注册时才暴露，收敛枚举面
	response.OK(c, gin.H{"sent": true})
}

// Register POST /api/auth/register —— 验证码 + 密码建号，成功即登录（直接发 token）。
func (h *Handler) Register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParameterErr(c)
		return
	}
	uid, err := h.auth.Register(c.Request.Context(), req.Email, req.Code, req.Password)
	if err != nil {
		respondAuthErr(c, err)
		return
	}
	response.OK(c, gin.H{"token": h.auth.Token(uid)})
}

// Login POST /api/auth/login。
func (h *Handler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParameterErr(c)
		return
	}
	token, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		respondAuthErr(c, err)
		return
	}
	response.OK(c, gin.H{"token": token})
}

// SetAPIKey POST /api/auth/apikey —— 用户自带 LLM Key（BYOK），传空串清除。
func (h *Handler) SetAPIKey(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParameterErr(c)
		return
	}
	if err := h.auth.SetAPIKey(c.Request.Context(), uid, req.APIKey); err != nil {
		respondAuthErr(c, err)
		return
	}
	response.OK(c, nil)
}

// Me GET /api/auth/me —— 当前登录用户摘要。
func (h *Handler) Me(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	email, hasKey, err := h.auth.Me(c.Request.Context(), uid)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"email": email, "has_api_key": hasKey})
}

func respondAuthErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, auth.ErrEmailInvalid), errors.Is(err, auth.ErrPasswordWeak),
		errors.Is(err, auth.ErrCodeInvalid), errors.Is(err, auth.ErrAPIKeyInvalid):
		response.Err(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, auth.ErrCodeCooldown):
		response.Err(c, http.StatusTooManyRequests, err.Error())
	case errors.Is(err, auth.ErrEmailTaken):
		response.Err(c, http.StatusConflict, err.Error())
	case errors.Is(err, auth.ErrBadCredentials):
		response.Err(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, auth.ErrCryptoUnavailable), errors.Is(err, auth.ErrStorageUnavailable):
		response.Err(c, http.StatusServiceUnavailable, err.Error())
	default:
		response.Err(c, http.StatusInternalServerError, err.Error())
	}
}
