package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RenderDeck(c *gin.Context){
	uid,deckID,ok := h.renderGrants.Take(c.Param("nonce"))
	if !ok{
		c.Status(http.StatusNotFound)
		return
	}
	html,err := h.decks.GetHTML(uid,deckID)
	if err != nil{
		c.Status(http.StatusNotFound)
	}
	deckPageHeaders(c)
	c.Data(http.StatusOK,"text/html; charset=utf-8",[]byte(html))
}