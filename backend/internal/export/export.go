package export

// deck-v2 导出：PDF / 逐页 PNG zip / 单文件 HTML。
//
// 三种产物都从"预览用的那份 HTML"出发（deck.PreviewHTML：style.css 已内联），
// 与用户在浏览器里看到的永远是同一份内容——导出即所见。

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"html-ppt/backend/internal/vision"
)

// Requester 导出需要的 deck 读取能力（由 deck.Service 实现，避免包循环）。
type Requester interface {
	PreviewHTML(userID uint, id string) (string, error)
	ExportDir(id string) string
	PageCount(id string) int
}

type Service struct {
	deck       Requester
	baseURL    string // 回环地址（渲染/打印走本服务的 /api/decks/:id/file?token= 或 nonce 通道）
	chromePath string
	timeout    time.Duration
	// grants 复用视觉审查的一次性门票：导出走同一条公开渲染通道，免鉴权头问题
	grants TicketIssuer
}

// TicketIssuer 一次性门票（与 vision.Grants 同一实例，main 装配）。
type TicketIssuer interface {
	Issue(uid uint, deckID string) (string, error)
}

func New(deck Requester, baseURL, chromePath string, timeout time.Duration, grants TicketIssuer) *Service {
	return &Service{deck: deck, baseURL: baseURL, chromePath: chromePath, timeout: timeout, grants: grants}
}

// Result 导出产物。
type Result struct {
	Path    string `json:"path"`    // 相对 deck 目录（exports/<文件名>）
	Filename string `json:"filename"`
	Size    int64  `json:"size"`
}

// Export 按格式导出。v1 同步实现（8 页实测 < 30s）。
func (s *Service) Export(ctx context.Context, uid uint, deckID, format string) (any, error) {
	nonce, err := s.grants.Issue(uid, deckID)
	if err != nil {
		return Result{}, fmt.Errorf("导出发票据失败: %w", err)
	}
	url := strings.TrimRight(s.baseURL, "/") + "/api/render/" + nonce
	dir := s.deck.ExportDir(deckID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Result{}, err
	}

	switch strings.ToLower(format) {
	case "pdf":
		return s.exportPDF(ctx, url, dir)
	case "png":
		return s.exportPNG(ctx, deckID, url, dir)
	case "html":
		return s.exportSingleFile(uid, deckID, dir)
	default:
		return Result{}, fmt.Errorf("不支持的导出格式 %q（pdf / png / html）", format)
	}
}

func (s *Service) exportPDF(ctx context.Context, url, dir string) (Result, error) {
	pdf, err := vision.PrintPDF2(ctx, url, s.chromePath, s.timeout)
	if err != nil {
		return Result{}, fmt.Errorf("PDF 渲染失败: %w", err)
	}
	name := "deck.pdf"
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, pdf, 0o644); err != nil {
		return Result{}, err
	}
	return Result{Path: "exports/" + name, Filename: name, Size: int64(len(pdf))}, nil
}

func (s *Service) exportPNG(ctx context.Context, deckID, url, dir string) (Result, error) {
	d, err := vision.CaptureV2(ctx, vision.OptionsV2{URL: url, ChromePath: s.chromePath})
	if err != nil {
		return Result{}, fmt.Errorf("页面捕获失败: %w", err)
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, sl := range d.Slides {
		if len(sl.PNG) == 0 {
			continue
		}
		w, err := zw.Create(fmt.Sprintf("page-%02d.png", sl.Index+1))
		if err != nil {
			return Result{}, err
		}
		if _, err := w.Write(sl.PNG); err != nil {
			return Result{}, err
		}
	}
	if err := zw.Close(); err != nil {
		return Result{}, err
	}
	name := "deck-png.zip"
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, buf.Bytes(), 0o644); err != nil {
		return Result{}, err
	}
	return Result{Path: "exports/" + name, Filename: name, Size: int64(buf.Len())}, nil
}

// styleLinkRe 与 deck 包的同一正则：单文件打包时把 style.css 的 link 换成内联 style。
var styleLinkRe = regexp.MustCompile(`<link[^>]*href="style\.css"[^>]*>`)

// exportSingleFile 读预览 HTML（style.css 已内联）直接落盘。
// runtime/base.css/字体保持 /assets 绝对引用——CJK 字体内联会把文件撑到 MB 级，
// 单文件版的定位是"分享与托管"，同源部署下资产可达；离线场景有系统字体回退。
func (s *Service) exportSingleFile(uid uint, deckID, dir string) (Result, error) {
	html, err := s.deck.PreviewHTML(uid, deckID)
	if err != nil {
		return Result{}, err
	}
	if styleLinkRe.MatchString(html) {
		return Result{}, fmt.Errorf("style.css 未内联（异常状态），拒绝导出不完整产物")
	}
	name := "deck.html"
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(html), 0o644); err != nil {
		return Result{}, err
	}
	return Result{Path: "exports/" + name, Filename: name, Size: int64(len(html))}, nil
}
