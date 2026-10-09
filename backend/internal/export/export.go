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

	"html-ppt/backend/internal/thumbs"
	"html-ppt/backend/internal/vision"
)

// Requester 导出需要的 deck 读取能力（由 deck.Service 实现，避免包循环）。
type Requester interface {
	PreviewHTML(userID uint, id string) (string, error)
	ExportDir(id string) string
	PageCount(id string) int
	// EnsureOwner 归属校验（别人的 deck 与不存在的 deck 同样报错）。
	EnsureOwner(userID uint, id string) error
	ThumbsDir(id string) string
	IndexPath(id string) (string, error)
}

type Service struct {
	deck       Requester
	baseURL    string // 回环地址（渲染/打印走本服务的 /api/decks/:id/file?token= 或 nonce 通道）
	chromePath string
	timeout    time.Duration
	// grants 复用视觉审查的一次性门票：导出走同一条公开渲染通道，免鉴权头问题
	grants TicketIssuer
	// assetsDir deck-v2 运行时资产目录（web/assets）：单文件导出时读取
	// runtime/base.css/animations/fonts 做内联（见 singlefile.go）
	assetsDir string
}

// TicketIssuer 一次性门票（与 vision.Grants 同一实例，main 装配）。
type TicketIssuer interface {
	Issue(uid uint, deckID string) (string, error)
}

func New(deck Requester, baseURL, chromePath string, timeout time.Duration, grants TicketIssuer, assetsDir string) *Service {
	return &Service{deck: deck, baseURL: baseURL, chromePath: chromePath, timeout: timeout, grants: grants, assetsDir: assetsDir}
}

// Result 导出产物。
type Result struct {
	Path    string `json:"path"` // 相对 deck 目录（exports/<文件名>）
	Filename string `json:"filename"`
	Size    int64  `json:"size"`
	// DownloadAs 下载落盘名（跟随文稿标题）；磁盘仍是 Filename，命名发生在下载头
	DownloadAs string `json:"download_as,omitempty"`
}

// Export 按格式导出。v1 同步实现（8 页实测 < 30s）。
func (s *Service) Export(ctx context.Context, uid uint, deckID, format string) (any, error) {
	if err := s.deck.EnsureOwner(uid, deckID); err != nil {
		return Result{}, err
	}
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

// exportSingleFile 读预览 HTML（style.css 已内联），再把翻页三件套与常用
// 字体内联成真正自包含的单文件：下载到本地、挪去任意静态托管，打开都是
// 完整的单页翻页版（见 singlefile.go）。等宽 CJK 全量字体除外，走系统回退。
func (s *Service) exportSingleFile(uid uint, deckID, dir string) (Result, error) {
	html, err := s.deck.PreviewHTML(uid, deckID)
	if err != nil {
		return Result{}, err
	}
	if styleLinkRe.MatchString(html) {
		return Result{}, fmt.Errorf("style.css 未内联（异常状态），拒绝导出不完整产物")
	}
	if html, err = s.makeSelfContained(html); err != nil {
		return Result{}, err
	}
	name := "deck.html"
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(html), 0o644); err != nil {
		return Result{}, err
	}
	return Result{Path: "exports/" + name, Filename: name, Size: int64(len(html))}, nil
}

// EnsureThumbs 缩略图缓存就绪：有效直接返回目录；过期/缺失就整本重渲一次
//（约 10-20s，之后按 deck 内容版本缓存，页级写入会触发失效）。
// 返回目录路径，调用方按 <no>.png 取图。
func (s *Service) EnsureThumbs(ctx context.Context, uid uint, deckID string) (string, error) {
	if err := s.deck.EnsureOwner(uid, deckID); err != nil {
		return "", err
	}
	dir := s.deck.ThumbsDir(deckID)
	indexPath, err := s.deck.IndexPath(deckID)
	if err != nil {
		return "", err
	}
	// 没生成过页面的 deck（还在澄清/大纲/选模板阶段）没有 index.html：
	// 快速失败，别把 chromedp 发出去渲占位页白等 2s 再报"渲染失败"
	if _, err := os.Stat(indexPath); err != nil {
		return "", fmt.Errorf("deck 尚未生成页面（index.html 不存在）")
	}
	if thumbs.Valid(dir, indexPath) {
		return dir, nil
	}
	nonce, err := s.grants.Issue(uid, deckID)
	if err != nil {
		return "", fmt.Errorf("缩略图发票据失败: %w", err)
	}
	url := strings.TrimRight(s.baseURL, "/") + "/api/render/" + nonce
	d, err := vision.CaptureV2(ctx, vision.OptionsV2{URL: url, ChromePath: s.chromePath, Timeout: s.timeout})
	if err != nil {
		return "", fmt.Errorf("缩略图渲染失败: %w", err)
	}
	pngs := make(map[int][]byte, len(d.Slides))
	for _, sl := range d.Slides {
		pngs[sl.Index+1] = sl.PNG
	}
	if err := thumbs.WriteAll(dir, indexPath, pngs); err != nil {
		return "", err
	}
	return dir, nil
}
