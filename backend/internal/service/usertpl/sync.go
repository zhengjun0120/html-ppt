package usertpl

// fork 模板的 layouts.md 与内置 base 的同步。放独立文件：这不是模板生命周期的
// 常规操作，而是"base 侧修了结构契约、旧 fork 也要吃到"的补偿机制。

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"

	"html-ppt/backend/internal/store"
)

// SyncLayoutsFromBase 把 fork 模板的 layouts.md 对齐到内置 base 的最新版（启动时跑一次）。
//
// fork 的 ADAPTATION.md 写明「结构契约（版式/类名）继承 base」，产品也没有编辑
// layouts.md 的入口——base 侧修了骨架/数量契约后，旧 fork 永远停在 fork 那一刻的
// 旧文件上。实测：ut-f044bad5 的 how-it-works 停在旧版裸文本骨架，三行短文撑不起
// 1080 高的画布、填充率永远 26%，agent 每轮量测都被追着改。
//
// 只同步 layouts.md：style.css / template.json 是用户定制（对话改 token）的落点，
// 不动。css 的判定见 cssCovers——fork 的 css 允许值上的漂移，缺类才拦；shellCSS
// 是 deck 外壳 base.css（骨架的 grid/notes/mt-l 等公共类由它提供），读不到就传
// nil，退化为"全部类都要求 fork css 有"的保守口径。
func (s *Service) SyncLayoutsFromBase(shellCSS []byte) {
	if s == nil || s.st == nil || s.reg == nil {
		return
	}
	var rows []store.UserTemplate
	if err := s.st.DB.Find(&rows).Error; err != nil {
		return
	}
	for _, row := range rows {
		if row.BaseID == "" {
			continue
		}
		baseDir, err := s.reg.BuiltinDir(row.BaseID)
		if err != nil {
			continue // base 已下线：fork 保持原样
		}
		dir := s.Dir(row.ID)
		baseL, err := os.ReadFile(filepath.Join(baseDir, "layouts.md"))
		if err != nil {
			continue
		}
		if fileEq(filepath.Join(dir, "layouts.md"), filepath.Join(baseDir, "layouts.md")) {
			continue // 已是最新
		}
		forkCSS, err := os.ReadFile(filepath.Join(dir, "style.css"))
		if err != nil || !cssCovers(forkCSS, shellCSS, string(baseL)) {
			log.Printf("[warn] 用户模板 %s 的 style.css 覆盖不了 base 新骨架的类，layouts.md 保持旧版（请重新 fork）", row.ID)
			continue
		}
		if err := os.WriteFile(filepath.Join(dir, "layouts.md"), baseL, 0o644); err != nil {
			log.Printf("[warn] 用户模板 %s 同步 layouts.md 失败: %v", row.ID, err)
			continue
		}
		log.Printf("[info] 用户模板 %s（fork 自 %s）的 layouts.md 已同步到 base 最新版", row.ID, row.BaseID)
	}
}

// fileEq 两个文件都存在且字节一致（任一读失败 = 不一致，调用方按"不同步"处理）。
func fileEq(a, b string) bool {
	ab, err1 := os.ReadFile(a)
	bb, err2 := os.ReadFile(b)
	return err1 == nil && err2 == nil && bytes.Equal(ab, bb)
}

// cssCovers 判断 layouts.md 骨架用到的类是否都有着落：deck 外壳 base.css（shellCSS）
// 提供的公共类直接放行；壳没有的类必须出现在 fork 的 style.css。类名取自骨架的
// class=" 属性（跳过首段——那不是骨架），数值漂移不管、缺类才拦。
func cssCovers(forkCSS, shellCSS []byte, layouts string) bool {
	segs := strings.Split(layouts, `class="`)
	for _, seg := range segs[1:] {
		i := strings.IndexByte(seg, '"')
		if i < 0 {
			continue
		}
		for _, c := range strings.Fields(seg[:i]) {
			if bytes.Contains(shellCSS, []byte("."+c)) {
				continue
			}
			if !bytes.Contains(forkCSS, []byte("."+c)) {
				return false
			}
		}
	}
	return true
}
