// Package tplthumb 模板卡片缩略图：把 demo 第 1 页用无头 Chrome 拍成 PNG，
// 画廊/我的模板页以 <img> 展示——替代「每张卡片一个活 iframe」（112 个活
// 文档的 CSS/JS/字体子资源是选模板页加载成本的大头，2026-10-02 用户拍板）。
//
// 失效模型：内容寻址。文件名带模板内容版本前缀（registry.ContentVersion），
// 编辑 → 重挂 → 版本变 → 清单下发新 Thumb → 前端 URL 变 → 旧图自然废弃。
// 渲染按需（第一个访问者等 1-3s，之后走磁盘与浏览器长缓存）。
package tplthumb

import (
	"context"
	"fmt"
	"log"
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
}

// MaxConcurrentRenders 同时进行的 Chrome 渲染数。108 张卡全冷时按 1-3s/张
// 排队灌满，先到先得；磁盘命中不占信号量。
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
func (s *Service) PNG(ctx context.Context, id, version string) ([]byte, error) {
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

	// 渲染全程限流；拿到信号量后再查一次盘——排在前面的同款请求可能刚写完
	s.sem <- struct{}{}
	defer func() { <-s.sem }()
	if png, err := os.ReadFile(path); err == nil {
		return png, nil
	}

	// preview?slide=1：单页 HTML（第 1 页），内置与 ut-*（挂载层）同一条路。
	// Shoot [0]：只要这一页的 PNG。
	deck, err := vision.CaptureV2(ctx, vision.OptionsV2{
		URL:        fmt.Sprintf("%s/api/templates/%s/preview?slide=1", s.baseURL, id),
		ChromePath: s.chromePath,
		Timeout:    s.timeout,
		Shoot:      []int{0},
	})
	if err != nil {
		return nil, fmt.Errorf("缩略图渲染失败: %w", err)
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

// pathFor 缩略图落盘路径：<dir>/<id>-<version>.png。id 与 version 都经过
// 调用方白名单（registry id 白名单 + ContentVersion 前缀比对），拼路径安全。
func (s *Service) pathFor(id, version string) string {
	return filepath.Join(s.dir, id+"-"+version+".png")
}
