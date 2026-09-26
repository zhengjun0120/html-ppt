package usertpl

// 编辑器手动保存模板 demo 的 index.html（docs/deck-editor-plan.md §4.2）。
//
// 覆盖前把当前版滚动备份到 data/user-template-backups/<ut-id>/，留最近 5 版。
// 备份目录刻意放在模板目录外：Publish 的 securityScan 与公开静态服务都只看
// 模板目录本身，备份混进去既会被扫描也会被当 demo 伺服出去。
// 定制对话（customize）只改 style.css，与本链路无文件交集。

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ErrPublished 已发布模板拒绝手改：已发布模板是社区可见内容，绕过发布门禁
// 改它不可接受；要改先下架（下架后 status 回 draft，保存放行）。
var ErrPublished = errors.New("模板已发布，先下架再编辑")

// keepTemplateBackups 滚动备份保留份数（§1 决策 8）。
const keepTemplateBackups = 5

// SaveIndexHTML 编辑器全量保存 index.html，返回滚动备份后的现存备份数。
// status == published 时拒绝；只动 index.html，template.json / style.css /
// layouts.md 不碰（定制对话的领域）。
func (s *Service) SaveIndexHTML(userID uint, id, html string) (int, error) {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return 0, err
	}
	if row.Status == "published" {
		return 0, ErrPublished
	}

	cur := filepath.Join(s.Dir(id), "index.html")
	kept := 0
	if data, err := os.ReadFile(cur); err == nil {
		bd := s.backupDir(id)
		if err := os.MkdirAll(bd, 0o755); err == nil {
			// 编辑器场景会毫秒级连保存：UnixMilli 撞名会让后一次备份覆盖前一次，
			// 悄悄丢备份——加递增后缀（%02d 保证同毫秒内字典序=时间序）。
			name := nextBackupName(bd, time.Now().UnixMilli())
			if err := os.WriteFile(filepath.Join(bd, name), data, 0o644); err == nil {
				kept = pruneTemplateBackups(bd)
			}
		}
		// 当前 index.html 读不到（首次保存前的异常态）不算失败：覆盖照常进行
	}

	if err := atomicWriteFile(cur, []byte(html)); err != nil {
		return kept, fmt.Errorf("写 index.html 失败: %w", err)
	}
	return kept, nil
}

// backupDir 备份根：data/user-template-backups/<ut-id>/（s.root = data/user-templates）。
func (s *Service) backupDir(id string) string {
	return filepath.Clean(filepath.Join(s.root, "..", "user-template-backups", id))
}

// nextBackupName 生成不与现存文件重名的备份名：index-<unixms>-<n>.html。
func nextBackupName(dir string, ts int64) string {
	for n := 0; ; n++ {
		name := fmt.Sprintf("index-%d-%02d.html", ts, n)
		if _, err := os.Stat(filepath.Join(dir, name)); os.IsNotExist(err) {
			return name
		}
	}
}

// pruneTemplateBackups 按文件名（unixms 时间戳字典序即时间序）保留最新 N 份。
func pruneTemplateBackups(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "index-") && strings.HasSuffix(e.Name(), ".html") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for len(names) > keepTemplateBackups {
		_ = os.Remove(filepath.Join(dir, names[0]))
		names = names[1:]
	}
	return len(names)
}

// atomicWriteFile 同目录临时文件 + rename，避免读到半截（与 deck 包同名 helper 同思路）。
func atomicWriteFile(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
