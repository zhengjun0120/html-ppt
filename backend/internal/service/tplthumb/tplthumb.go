// Package tplthumb 模板卡片缩略图：把 demo 第 1 页用无头 Chrome 拍成 PNG，
// 画廊/我的模板页以 <img> 展示——替代「每张卡片一个活 iframe」（112 个活
// 文档的 CSS/JS/字体子资源是选模板页加载成本的大头，2026-10-02 用户拍板）。
//
// 失效模型：内容寻址。文件名带模板内容版本前缀（registry.ContentVersion），
// 编辑 → 重挂 → 版本变 → 清单下发新 Thumb → 前端 URL 变 → 旧图自然废弃。
// 渲染两条路：启动预热（Prewarm，内置模板后台渲完）+ 按需兜底（第一个
// 访问者等 1-3s），之后全走磁盘与浏览器长缓存。
package tplthumb

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/vision"
)

// Service 模板缩略图服务。渲染走 vision.CaptureV2 对公开 preview URL 截图，
// 并发由 sem 限流（Chrome 实例又重又贵，成百张卡同时请求不能全量并发）。
type Service struct {
	reg        *template.Registry
	baseURL    string // 自身可达地址（渲染浏览器从这里拉 preview 页）
	chromePath string
	dir        string // PNG 落盘目录（<data>/template-thumbs）
	timeout    time.Duration
	sem        chan struct{}
	// renderFn 渲染出口（单测注入假实现）；nil = 真无头 Chrome（vision.CaptureV2）。
	renderFn func(ctx context.Context, url string) (*vision.Deck2, error)
}

// MaxConcurrentRenders 同时进行的 Chrome 渲染数。108 张卡全冷时按 1-3s/张
// 排队灌满，先到先得；磁盘命中不占信号量。预热（Prewarm）一次只占一个槽。
const MaxConcurrentRenders = 2

func New(reg *template.Registry, baseURL, chromePath, dir string, timeout time.Duration) *Service {
	return &Service{
		reg:        reg,
		baseURL:    strings.TrimRight(baseURL, "/"),
		chromePath: chromePath,
		dir:        dir,
		timeout:    timeout,
		sem:        make(chan struct{}, MaxConcurrentRenders),
	}
}

// PNG 返回模板第 1 页缩略图。version 是清单下发的内容版本前缀（调用方已
// 校验与当前一致）；磁盘命中直接回，未命中渲染后落盘。Chrome 不可用或渲染
// 失败返回错误——调用方（handler）映射成 404，前端回退活 iframe。
//
// userToken 是发起请求的用户的登录态（img 的 ?token= 原样透传）：草稿态
// ut-* 的 demo 与 style.css 只对属主开放，渲染浏览器必须带着同一个 token
// 才能拿到带样式的页面——demo 路由会把 token 续写进 style.css 的 href，
// 子资源请求随之通过。内置模板公开，token 用不上。传空且模板是草稿 →
// 渲染出无样式页或 404，属调用方未按契约传 token。
func (s *Service) PNG(ctx context.Context, id, version, userToken string) ([]byte, error) {
	if s.reg == nil {
		return nil, fmt.Errorf("模板库不可用")
	}
	if _, err := s.reg.Get(id); err != nil {
		return nil, err
	}
	path := s.pathFor(id, version)
	if png, err := os.ReadFile(path); err == nil {
		return png, nil
	}

	// 渲染 URL：内置走公开 preview?slide=1（单页）；ut-* 走鉴权 demo 路由
	//（无 trim，CaptureV2 只拍第 1 页；token 经 href 重写传给 style.css）。
	var renderURL string
	if strings.HasPrefix(id, "ut-") {
		renderURL = fmt.Sprintf("%s/api/user-templates/%s/demo", s.baseURL, id)
		if userToken != "" {
			renderURL += "?token=" + userToken
		}
	} else {
		renderURL = fmt.Sprintf("%s/api/templates/%s/preview?slide=1", s.baseURL, id)
	}
	return s.renderStore(ctx, path, renderURL)
}

// renderStore 渲染一张并存盘（PNG 与预热共用）：限流拿信号量 → 拿到后再查
// 一次盘（排在前面的同款请求可能刚写完）→ 截图（失败重试一次）→ 原子落盘。
func (s *Service) renderStore(ctx context.Context, path, renderURL string) ([]byte, error) {
	s.sem <- struct{}{}
	defer func() { <-s.sem }()
	if png, err := os.ReadFile(path); err == nil {
		return png, nil
	}

	deck, err := s.capture(ctx, renderURL)
	if err != nil {
		// 冷启动竞态兜底：全新 Chrome + 首次加载的 demo，runtime 偶发未就绪
		//（实测 1.7s 快速失败、重跑同 URL 即成功）。短暂等待后重试一次。
		time.Sleep(2 * time.Second)
		deck, err = s.capture(ctx, renderURL)
		if err != nil {
			return nil, fmt.Errorf("缩略图渲染失败: %w", err)
		}
	}
	if len(deck.Slides) == 0 || len(deck.Slides[0].PNG) == 0 {
		return nil, fmt.Errorf("缩略图渲染结果为空")
	}
	png := deck.Slides[0].PNG
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		log.Printf("[warn] tplthumb: 建缓存目录失败（图不落盘，仅本次返回） err:%v", err)
		return png, nil
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, png, 0o644); err == nil {
		_ = os.Rename(tmp, path) // 原子落盘：并发读到半张图比多渲一次更糟
	} else {
		log.Printf("[warn] tplthumb: 缩略图落盘失败 err:%v", err)
	}
	return png, nil
}

// capture 渲染出口：真 Chrome 走 vision.CaptureV2，只拍第 1 页。
func (s *Service) capture(ctx context.Context, renderURL string) (*vision.Deck2, error) {
	if s.renderFn != nil {
		return s.renderFn(ctx, renderURL)
	}
	return vision.CaptureV2(ctx, vision.OptionsV2{
		URL:        renderURL,
		ChromePath: s.chromePath,
		Timeout:    s.timeout,
		Shoot:      []int{0},
	})
}

// Prewarm 后台预热：服务启动后把内置模板的第 1 页缩略图提前渲好，让
// 「第一次全量浏览」也全走缓存命中（2026-10-02 用户拍板）。只做内置——
// ut-* 的 demo 路由对草稿只认属主 token，服务端没有 token，渲出来是
// 无样式页且会污染磁盘缓存；ut 数量少，维持按需渲染。
//
// 对线上流量谦让：只在两个渲染槽都空闲时才占一个、一次只渲一张，用户
// 请求永远排在预热前面。连续 3 张双试均失败视为环境问题（无 Chrome、
// 端口不通），直接中止别把全库硬刷一遍。
func (s *Service) Prewarm(ctx context.Context) {
	if s.reg == nil {
		return
	}
	if !s.waitSelf(ctx) {
		return
	}
	var pending []string
	hits := 0
	for _, m := range s.reg.List() {
		if strings.HasPrefix(m.ID, "ut-") {
			continue
		}
		v, err := s.reg.ContentVersion(m.ID)
		if err != nil || len(v) < 8 {
			continue
		}
		if _, err := os.ReadFile(s.pathFor(m.ID, v[:8])); err == nil {
			hits++
			continue
		}
		pending = append(pending, m.ID)
	}
	log.Printf("[info] tplthumb: 预热开始 内置待渲染 %d 张，磁盘命中 %d 张", len(pending), hits)
	if len(pending) == 0 {
		return
	}

	start := time.Now()
	ok, fail, consec := 0, 0, 0
	for _, id := range pending {
		if ctx.Err() != nil {
			return
		}
		// 空闲门控：两个渲染槽都空才轮到预热（用户请求优先）。
		for len(s.sem) > 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
			}
		}
		v, err := s.reg.ContentVersion(id)
		if err != nil || len(v) < 8 {
			continue
		}
		path := s.pathFor(id, v[:8])
		if _, err := os.ReadFile(path); err == nil {
			ok++ // 排队期间被真实请求渲好了
			continue
		}
		renderURL := fmt.Sprintf("%s/api/templates/%s/preview?slide=1", s.baseURL, id)
		if _, err := s.renderStore(ctx, path, renderURL); err != nil {
			fail++
			consec++
			log.Printf("[warn] tplthumb: 预热渲染失败 %s err:%v", id, err)
			if consec >= 3 {
				log.Printf("[warn] tplthumb: 连续 %d 张失败，预热中止（疑似环境问题）", consec)
				return
			}
			continue
		}
		consec = 0
		ok++
	}
	log.Printf("[info] tplthumb: 预热完成 成功 %d 失败 %d（磁盘命中 %d），耗时 %s",
		ok, fail, hits, time.Since(start).Round(time.Second))
}

// waitSelf 等自身 HTTP 可达：预热在监听 goroutine 发射后立即启动，存在
// 竞态；渲染浏览器要从本服务拉页面，起不来就免谈。最多等 60s。
func (s *Service) waitSelf(ctx context.Context) bool {
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/api/health", nil)
		if err == nil {
			resp, err := http.DefaultClient.Do(req)
			if err == nil {
				resp.Body.Close()
				return true
			}
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(500 * time.Millisecond):
		}
	}
	log.Printf("[warn] tplthumb: 预热放弃——服务自身 60s 内不可达")
	return false
}

// pathFor 缩略图落盘路径：<dir>/<id>-<version>.png。id 与 version 都经过
// 调用方白名单（registry id 白名单 + ContentVersion 前缀比对），拼路径安全。
func (s *Service) pathFor(id, version string) string {
	return filepath.Join(s.dir, id+"-"+version+".png")
}
