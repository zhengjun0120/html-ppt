package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RenderDeck(c *gin.Context) {
	uid, deckID, ok := h.renderGrants.Take(c.Param("nonce"))
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	html, err := h.decks.GetHTML(uid, deckID)
	if err != nil {
		// 这里少了 return 的话，404 之后还会继续往下写一个 200 的空页面——
		// 无头浏览器拿到 200 + 空 HTML，Reveal 起不来，报错却指向"JS 报错？"，
		// 排查时会绕远路（实测：状态码和内容都对不上时最难查）
		c.Status(http.StatusNotFound)
		return
	}
	deckPageHeaders(c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
