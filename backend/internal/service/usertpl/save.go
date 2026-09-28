package usertpl

// 编辑器手动保存模板文件（docs/user-template-history-plan.md §3.2/§3.3）。
//
// 历史系统上线前的旧机制（滚动备份 5 版到 data/user-template-backups/）已退役：
// 每次保存现在记一条 edit 版本（bundle 五件套整包），回滚走 /history/:version/restore。
// published 拒绝的语义不变（D2）。

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ErrPublished 已发布模板拒绝一切写操作（手动编辑/对话定制/回滚）：
// 已发布模板是社区可见内容，绕过发布门禁改它不可接受；要改先下架
// （下架后 status 回 draft，写操作放行）。
var ErrPublished = errors.New("模板已发布，先下架再编辑")

// SaveIndexHTML 编辑器全量保存 index.html，成功后记一条 edit 版本。
// 只动 index.html，template.json / style.css / layouts.md 不碰（各自路径的领域）。
func (s *Service) SaveIndexHTML(userID uint, id, html string) error {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return err
	}
	if row.Status == "published" {
		return ErrPublished
	}
	if err := scanHTML(html); err != nil {
		return err
	}
	if err := s.writeTemplateFile(id, "index.html", html); err != nil {
		return err
	}
	if err := s.recordVersionUT(id, OpEdit, "手动编辑 index.html"); err != nil {
		return err
	}
	// 重挂失败不阻断保存：文件与版本已落地，预览读磁盘不受注册表影响；
	// 校验问题由发布门禁终审（日志留痕）。
	if err := s.remountUT(id); err != nil {
		log.Printf("[warn] 模板 %s 保存后重挂失败: %v", id, err)
	}
	return nil
}

// SaveStyleCSS 手动保存 style.css（工作台"样式"面板）：安全预检 → 整文件写入 →
// 记 edit 版本。重挂失败不阻断（同上）。
func (s *Service) SaveStyleCSS(userID uint, id, css string) error {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return err
	}
	if row.Status == "published" {
		return ErrPublished
	}
	if err := scanCSS(css); err != nil {
		return err
	}
	if err := s.writeTemplateFile(id, "style.css", css); err != nil {
		return err
	}
	if err := s.recordVersionUT(id, OpEdit, "手动编辑 style.css"); err != nil {
		return err
	}
	if err := s.remountUT(id); err != nil {
		log.Printf("[warn] 模板 %s 样式保存后重挂失败: %v", id, err)
	}
	return nil
}

// writeTemplateFile 整文件写入模板目录内的白名单文件。
// 刻意不含重挂：注册表校验的是整个模板目录——demo 里的一个坏版式会让
// 之后所有 style.css 保存都 400（文件已写却报错、版本不记，盘与历史脱节，
// 实测事故）。写入/记版必须先落地；重挂由调用方按语义处理（见 remountUT）。
func (s *Service) writeTemplateFile(id, name, content string) error {
	switch name {
	case "index.html", "style.css":
	default:
		return fmt.Errorf("不允许写 %s", name)
	}
	return atomicWriteFile(filepath.Join(s.Dir(id), name), []byte(content))
}

// remountUT 重挂注册表：生成管线读注册表快照，文件改完要重挂才生效。
// 失败上抛给调用方决定呈现方式（对话工具转告模型自修；手动路径降级为日志）。
func (s *Service) remountUT(id string) error {
	if s.reg == nil {
		return nil
	}
	return s.reg.MountUser(s.Dir(id))
}

// atomicWriteFile 同目录临时文件 + rename，避免读到半截（与 deck 包同名 helper 同思路）。
func atomicWriteFile(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
