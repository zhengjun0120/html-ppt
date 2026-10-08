package handler

// 列表/文档类大响应的协商缓存（2026-09-30 预览缓存化第三批）：
// 内容哈希 ETag + Cache-Control: no-cache，If-None-Match 命中回 304 零传输。
//
// 为什么是 no-cache 而不是 max-age：这些端点的内容会随用户操作变
// （fork 挂新模板、定制对话改 style.css），强缓存会让改动在浏览器侧滞后；
// no-cache = 每次回源问一句，内容没变 304、变了拿新版——新鲜度与省流量兼得。
// 适用判据：响应体 ≥几十 KB 且重访命中率高的才值得（模板清单 296KB、
// ut demo/style.css）；小的 JSON（列表页分页之类）省不了多少，别套。

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/response"
)

// etagJSON 协商缓存版成功响应：序列化 → 哈希 → 条件请求判定。
// 与 response.OK 同语义（200 + 数据本体），只是会算 ETag。
func etagJSON(c *gin.Context, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	sum := sha256.Sum256(body)
	etag := `"` + hex.EncodeToString(sum[:16]) + `"`
	c.Header("Cache-Control", "no-cache")
	c.Header("ETag", etag)
	if c.GetHeader("If-None-Match") == etag {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", body)
}

// etagData 文件流版协商缓存：与 etagJSON 同一套头与 304 判定，正文由调用方
// 给定（已是最终字节——重写/注入完成后的内容，token 之类会话态在字节里，
// 换会话 ETag 自然变，不会串）。
func etagData(c *gin.Context, contentType string, body []byte) {
	sum := sha256.Sum256(body)
	etag := `"` + hex.EncodeToString(sum[:16]) + `"`
	c.Header("Cache-Control", "no-cache")
	c.Header("ETag", etag)
	if c.GetHeader("If-None-Match") == etag {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}
	c.Data(http.StatusOK, contentType, body)
}
