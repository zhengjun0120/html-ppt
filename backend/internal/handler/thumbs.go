package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
	"html-ppt/backend/internal/thumbs"
)

// DeckThumbs GET /api/decks/:id/thumbs —— 缩略图清单。
//
// 首次访问会触发整本渲染（无头浏览器，约 10-20s），之后按 deck 内容版本缓存
// （页面任何写入都会失效缓存）。返回已缓存的页码列表；前端预览栏与文稿列表
// 封面都吃这里。
func (h *Handler) DeckThumbs(c *gin.Context) {
	if h.exporter == nil {
		response.Err(c, http.StatusServiceUnavailable, "渲染服务不可用（Chrome 未配置）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	dir, err := h.exporter.EnsureThumbs(c.Request.Context(), uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "deck 不存在或缩略图渲染失败")
		return
	}
	pages := thumbs.Pages(dir)
	if pages == nil {
		response.Err(c, http.StatusInternalServerError, "缩略图为空")
		return
	}
	response.OK(c, gin.H{"pages": pages, "count": len(pages)})
}

// DeckThumb GET /api/decks/:id/thumbs/:no —— 单页缩略图 PNG。
// 与 file 接口同为文件流：404 统一表示"拿不到"，不区分原因。
func (h *Handler) DeckThumb(c *gin.Context) {
	if h.exporter == nil {
		response.Err(c, http.StatusServiceUnavailable, "渲染服务不可用（Chrome 未配置）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	no, err := strconv.Atoi(c.Param("no"))
	if err != nil || no < 1 || no > 999 {
		response.Err(c, http.StatusNotFound, "页码不合法")
		return
	}
	dir, err := h.exporter.EnsureThumbs(c.Request.Context(), uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "deck 不存在或缩略图渲染失败")
		return
	}
	p := filepath.Join(dir, strconv.Itoa(no)+".png")
	if _, err := os.Stat(p); err != nil {
		// 缓存有效但单页缺失（partial write 等异常态）：下次整本重渲兜底
		response.Err(c, http.StatusNotFound, "该页缩略图不存在")
		return
	}
	// 缓存由写路径主动失效，短缓存只是挡连点
	c.Header("Cache-Control", "private, max-age=60")
	c.Data(http.StatusOK, "image/png", mustRead(p))
}

func mustRead(p string) []byte {
	b, _ := os.ReadFile(p)
	return b
}
