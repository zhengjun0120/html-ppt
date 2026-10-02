package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/response"
	"html-ppt/backend/internal/service/template"
)

// ListTemplates GET /api/templates —— 画廊的模板清单（公开）。
//
// 为什么公开：模板元数据不含任何用户数据，画廊在登录页之外的展示场景
// （落地页、分享预览）也需要它。体积小（每模板一份元数据），不做分页。
//
// 协商缓存：完整 Meta 含版式数组，102 个模板实测 ~296KB——选模板页/我的模板
// 页每次进页都拉一遍太重。内容哈希 ETag + no-cache：清单没变（绝大多数重访）
// 304 零传输；fork/发布/挂载后内容变，立刻拿新版。
func (h *Handler) ListTemplates(c *gin.Context) {
	if h.templates == nil {
		response.Err(c, http.StatusServiceUnavailable, "模板库不可用")
		return
	}
	etagJSON(c, h.templates.ListSummaries())
}

// GetTemplate GET /api/templates/:id —— 单个模板详情（含版式索引，前端向导展示用）。
func (h *Handler) GetTemplate(c *gin.Context) {
	if h.templates == nil {
		response.Err(c, http.StatusServiceUnavailable, "模板库不可用")
		return
	}
	t, err := h.templates.Get(c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "模板不存在")
		return
	}
	response.OK(c, t.Meta)
}

// PreviewTemplate GET /api/templates/:id/preview?variant=<vid>&slide=<n> —— demo 页 HTML。
//
// 选模板页的实时换肤预览：variant 非空时服务端把变体 class 挂到 body，
// 前端切变体只改 iframe src 的 query。demo 里的相对引用（style.css）改写为
// 公开静态路由 /templates/<id>/ 的绝对路径——本端点自身带鉴权，资产走静态。
//
// slide>=1 是缩略模式：只返回第 n 页。画廊一屏十几张卡，每张都为整本 demo
// 付解析+布局的 CPU 账，缩略卡只看得到第一页，没必要把其余十几页也渲染了
// （TrimToSlide 的 fail-open 语义见 service 层注释）。
func (h *Handler) PreviewTemplate(c *gin.Context) {
	if h.templates == nil {
		response.Err(c, http.StatusServiceUnavailable, "模板库不可用")
		return
	}
	t, err := h.templates.Get(c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "模板不存在")
		return
	}
	html := t.DemoHTML(c.Query("variant"))
	if v := c.Query("slide"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 {
			html = template.TrimToSlide(html, n)
		}
	}
	// 相对引用 → 对应静态路由的绝对引用（内置 = /templates/<id>/，用户 = /user-templates/<id>/）
	base := "/templates/" + t.ID + "/"
	if strings.HasPrefix(t.ID, "ut-") {
		base = "/user-templates/" + t.ID + "/"
	}
	html = strings.ReplaceAll(html, `href="style.css"`, `href="`+base+`style.css"`)

	// 缓存策略分叉：
	//   内置模板（绝大多数卡）——demo 运行期不可变，1 小时强缓存。重进选模板页
	//     时上百张缩略的 HTML、css、js 全部零请求，加载近零；过期后 ETag 协商
	//     兜底（内容变 ETag 必变）。开发期改内置模板文件后 Ctrl+F5 强刷即可。
	//   用户模板（ut-）——定制对话随时改 style.css，保持 no-cache 逐次校验
	//     （304 零传输），编辑完回画廊立刻看到新样子。
	// 原先统一 no-store：选模板页每次进页全量重拉上百张缩略（合计 4-10MB），
	// 2026-09-30 用户实测反馈"每次进来都会重新加载"。
	sum := sha256.Sum256([]byte(html))
	etag := `"` + hex.EncodeToString(sum[:16]) + `"`
	c.Header("ETag", etag)
	if strings.HasPrefix(t.ID, "ut-") {
		c.Header("Cache-Control", "no-cache")
		if c.GetHeader("If-None-Match") == etag {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}
	} else {
		c.Header("Cache-Control", "private, max-age=3600")
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// TemplateThumb GET /api/templates/:id/thumb?v=<版本> —— 卡片缩略图 PNG
// （tplthumb 按需渲染 + 磁盘缓存）。v 是清单下发的 Thumb 版本：不匹配（缺/
// 旧）一律 404——img 地址内容寻址，版本对了才可能有对应文件；前端对 404
// 回退活 iframe。URL 含版本 → 浏览器长缓存 immutable，内容变 URL 自然变。
// Chrome 不可用（服务未装配）同样 404，前端整体回退，画廊照常可用。
func (h *Handler) TemplateThumb(c *gin.Context) {
	if h.thumbs == nil || h.templates == nil {
		c.Status(http.StatusNotFound)
		return
	}
	id := c.Param("id")
	v := c.Query("v")
	cur, err := h.templates.ContentVersion(id)
	if err != nil || len(cur) < 8 || v != cur[:8] {
		c.Status(http.StatusNotFound)
		return
	}
	png, err := h.thumbs.PNG(c.Request.Context(), id, cur[:8])
	if err != nil {
		log.Printf("[warn] tplthumb: %s 渲染失败 err:%v", id, err)
		c.Status(http.StatusNotFound)
		return
	}
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Data(http.StatusOK, "image/png", png)
}
