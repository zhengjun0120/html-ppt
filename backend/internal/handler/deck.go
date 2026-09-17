package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
)

// ListDecks GET /api/decks —— 当前用户的 deck 列表（id + 标题）。
func (h *Handler) ListDecks(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	decks, err := h.decks.List(uid)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, decks)
}

// GetDeckFile GET /api/decks/:id/file —— 返回 deck.html 给 <iframe> 预览。
// 注意：这是文件流不是 JSON API，走 c.Data 原生写法，不经过 response 包。
// 错误统一映射为 404（非法 id、不存在的 id、别人的 id 都当"找不到"，
// 不泄露存在性）。deckPageHeaders 挂沙箱配套的安全响应头。
func (h *Handler) GetDeckFile(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	html, err := h.decks.PreviewHTML(uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "deck 不存在")
		return
	}
	deckPageHeaders(c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func deckPageHeaders(c *gin.Context) {
	c.Header("Content-Security-Policy",
		"default-src 'self'; "+
			"script-src 'self'; "+
			"style-src 'self' 'unsafe-inline'; "+
			// font-src 只允许同源：字体必须自托管（web/assets/fonts/）。
			// 严格说它是 default-src 'self' 的复述，写出来是为了让"字体不走 CDN"
			// 这条决定出现在强制执行它的地方——将来谁想往 <head> 里加一行
			// Google Fonts，撞上的第一堵墙就是这里（而且它会顺便把访问者的 IP
			// 送给第三方，对一个本地工具没有理由）。
			"font-src 'self'; "+
			"img-src 'self' data:; "+
			"connect-src 'none'; "+
			// form-action 不受 connect-src 管辖（表单提交是导航不是 fetch），
			// 不限的话隐藏表单自动提交外站是一条数据渗出通道
			"form-action 'none'; "+
			"object-src 'none'; "+
			"base-uri 'self'; "+
			"frame-ancestors 'self'")
	// ?token= 走 URL 的妥协（见 auth.go）：no-referrer 保证任何外发导航
	// 都不会把带 token 的完整 URL 泄给第三方
	c.Header("Referrer-Policy", "no-referrer")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "no-store")
}

func (h *Handler) DeleteDeckVersion(c *gin.Context) {
	uid,ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c,http.StatusUnauthorized,"未登录")
		return
	}

	deckID := c.Param("id")
	version := c.Param("version")
	if deckID == "" || version == ""{
		response.ParameterErr(c)
		return
	}

	if err := h.decks.EnsureOwner(uid,deckID);err !=nil{
		response.Err(c,http.StatusNotFound,"deck 不存在")
		return
	}
	if err := h.decks.DeleteVersion(uid,deckID,version); err !=nil{
		response.Err(c,http.StatusBadRequest,err.Error())
		return
	}
	response.OK(c,nil)
}

func (h *Handler) ClearDeckHistory(c *gin.Context) {
	uid,ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c,http.StatusUnauthorized,"未登录")
		return
	}
	deckID := c.Param("id")
	if deckID == ""{
		response.ParameterErr(c)
		return
	}

	if err := h.decks.EnsureOwner(uid,deckID); err !=nil{
		response.Err(c,http.StatusNotFound,"deck 不存在 err:"+err.Error())
		return
	}
	n,err := h.decks.ClearHistory(uid,deckID)
	if err !=nil {
		response.Err(c,http.StatusInternalServerError,err.Error())
		return
	}
	response.OK(c,gin.H{"deleted":n})
}

func (h *Handler) ListDeckHistory(c *gin.Context){
	uid,ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c,http.StatusUnauthorized,"未登录")
		return
	}

	deckID := c.Param("id")
	if deckID == ""{
		response.ParameterErr(c)
		return
	}

	if err := h.decks.EnsureOwner(uid,deckID);err !=nil{
		response.Err(c,http.StatusNotFound,"deck 不存在 err:"+err.Error())
		return
	}

	versions,err := h.decks.ListVersions(uid,deckID)
	if err !=nil{
		response.Err(c,http.StatusInternalServerError,err.Error())
		return
	}
	response.OK(c,versions)
}

func (h *Handler) RestoreDeckVersion(c *gin.Context){
	uid,ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c,http.StatusUnauthorized,"未登录")
		return
	}
	deckID:= c.Param("id")
	version := c.Param("version")
	if deckID == "" || version == ""{
		response.ParameterErr(c)
		return
	}

	if err := h.decks.EnsureOwner(uid,deckID); err !=nil{
		response.Err(c,http.StatusNotFound,"deck 不存在 err:" +err.Error())
		return
	}
	if err := h.decks.RestoreVersion(uid,deckID,version);err !=nil{
		response.Err(c,http.StatusBadRequest,err.Error())
		return
	}
	response.OK(c,gin.H{"restored":version})
}