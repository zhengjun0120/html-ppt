package usertpl

// 编辑器手动保存模板文件（docs/user-template-history-plan.md §3.2/§3.3）。
//
// 历史系统上线前的旧机制（滚动备份 5 版到 data/user-template-backups/）已退役：
// 每次保存现在记一条 edit 版本（bundle 五件套整包），回滚走 /history/:version/restore。
// published 拒绝的语义不变（D2）。

import (
	"errors"
	"fmt"
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
	if err := atomicWriteFile(filepath.Join(s.Dir(id), "index.html"), []byte(html)); err != nil {
		return fmt.Errorf("写 index.html 失败: %w", err)
	}
	return s.recordVersionUT(id, OpEdit, "手动编辑 index.html")
}

// atomicWriteFile 同目录临时文件 + rename，避免读到半截（与 deck 包同名 helper 同思路）。
func atomicWriteFile(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
