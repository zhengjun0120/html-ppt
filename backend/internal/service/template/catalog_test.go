package template

import (
	"os"
	"path/filepath"
	"testing"
)

// ompTemplates oh-my-ppt 风格导入工程（docs/ohmyppt-styles-import-plan.md）产出的
// 全部模板 id。本测试是这份目录资产的回归清单：任何一个被误删、改坏契约或漏带
// 归属声明都会在这里炸出来。上游 arcsin1/oh-my-ppt（Apache-2.0）共 84 个风格，
// 其中 soft-pastel 与本仓库既有内置模板同名同定位，按拍板跳过，不在此列。
var ompTemplates = []string{
	// 试点批（中文风 12）
	"ink-wash-jiangnan", "palace-ink-red", "chinese-porcelain-rose", "chinese-cream-blossom",
	"song-rain-poetic", "chinese-ink-landscape", "celadon-bamboo", "chinese-fresh-trio",
	"chinese-pastel-spring", "indigo-lotus", "oriental-poetic-illustration", "gold-ivory",
	// 波 1（暗色科技 / 效果戏剧 24）
	"terminal-green", "dracula", "nord", "gruvbox-dark", "rose-pine", "tokyo-night",
	"catppuccin-mocha", "catppuccin-latte", "solarized-light", "sharp-mono",
	"cyberpunk-neon", "vaporwave", "y2k-chrome", "neon-haze", "arctic-cool",
	"glassmorphism", "aurora", "gradient-cosmic", "starlight-fireworks",
	"blueprint", "engineering-whiteprint", "geography-classroom",
	"burgundy-premium", "olive-elegant",
	// 波 2（学术商务 / 设计宣言 / 杂志编辑 26）
	"academic-navy", "academic-paper", "sage-academia", "red-research-framework",
	"corporate-clean", "blue-orange-analytics", "blue-white-chart", "classic-duo-blue",
	"acid-blue-business", "pitch-deck-vc", "swiss-international", "swiss-grid",
	"bauhaus", "neo-brutalism", "memphis-pop", "minimal-white", "japanese-minimal",
	"industrial-kaizen", "red-gold-ceremony", "premium-color-blocking", "magazine-bold",
	"e-ink-editorial", "editorial-serif", "midcentury", "news-broadcast", "retro-tv",
	// 波 3（手绘治愈 / 暖色杂项 21）
	"hand-drawn-autumn", "handdrawn-watercolor", "mountain-green-literary",
	"children-warm-orange", "macaron-mist", "sakura-soft-healing", "healing-color-card",
	"cream-pastel", "xiaohongshu-white", "mint-fresh", "summer-warm-color", "sunset-warm",
	"orange-sea", "dreamy-pink-gradient", "dreamy-romance", "dopamine-clash",
	"splash-abstract", "starry-dust", "rainbow-gradient", "amber-aurora", "cobalt-sunshine",
}

// TestOMPTemplateCatalog 导入工程的每一份模板都必须：注册成功、版式契约健全
// （9-10 版式、角色覆盖、视觉模式多样性）、归属声明完整（Apache-2.0 + NOTICE 义务）。
func TestOMPTemplateCatalog(t *testing.T) {
	dir, assets := repoTemplatesDir(t)
	reg, err := NewRegistry(dir, assets)
	if err != nil {
		t.Fatalf("模板注册失败: %v", err)
	}
	for _, id := range ompTemplates {
		t.Run(id, func(t *testing.T) {
			tp, err := reg.Get(id)
			if err != nil {
				t.Fatalf("未注册: %v", err)
			}
			if n := len(tp.Layouts); n != 9 && n != 10 {
				t.Errorf("版式数 %d，应为 9（基础集）或 10（含 code）", n)
			}
			roles := map[string]bool{}
			for _, l := range tp.Layouts {
				for _, r := range l.Roles {
					roles[r] = true
				}
			}
			for _, must := range []string{"cover", "content", "thanks"} {
				if !roles[must] {
					t.Errorf("roles 缺少 %q（生成侧大纲必用）", must)
				}
			}
			if tp.DistinctPatterns() < 4 {
				t.Errorf("视觉模式仅 %d 种（<4），节奏守卫会无解", tp.DistinctPatterns())
			}
			if tp.Source.License != "Apache-2.0" {
				t.Errorf("source.license = %q，归属声明被改", tp.Source.License)
			}
			if len(tp.Source.DerivedFrom) == 0 {
				t.Error("source.derived_from 为空，来源不可追溯")
			}
			if _, err := os.Stat(filepath.Join(dir, id, "UPSTREAM-README.md")); err != nil {
				t.Error("缺 UPSTREAM-README.md（上游 NOTICE 归属义务）")
			}
		})
	}
}
