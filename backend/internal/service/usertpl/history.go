package usertpl

// 用户模板的历史版本系统（docs/user-template-history-plan.md §3.1）。
//
// 机制镜像 deck/history.go + deck/v2_history.go：bundle 快照 + index.json 索引 +
// 单调序号 + 回滚也记版 + 分层裁剪。刻意不为复用去泛化 deck 包——那条链路已验证，
// 保持零改动；这里按模板的形状（五件套 bundle）对称实现一份。
//
// 与 deck 版的差异：
//   - bundle 是模板五件套（template.json/index.html/style.css/layouts.md/rules.md），
//     ADAPTATION.md 是溯源说明不进快照；
//   - VersionMeta 多一个 Changed（与上一版相比变化的文件名，记录时算好，前端免算）；
//   - prune 额外保证起点版本（最早的 fork 基线）永不裁——"回得去起点"是本特性的卖点；
//   - 恢复时把 bundle.template.json 里的名称/描述同步回 DB（元数据以快照为准）。

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"html-ppt/backend/internal/store"
)

const (
	maxVersionsPerUT = 100
	keepFullDays     = 7
)

const (
	OpFork    = "fork"    // 克隆初始化（起点基线）
	OpChat    = "chat"    // 定制对话（一轮有实际文件改动 = 一版）
	OpEdit    = "edit"    // 手动保存（index.html / style.css）
	OpMeta    = "meta"    // 改名/描述
	OpRestore = "restore" // 回滚本身也记版，历史永不重写
)

var utVersionPattern = regexp.MustCompile(`^v\d{6}$`)

// utFiles 快照覆盖的文件（与 forkCopies 同一套）。顺序即 Changed 的输出顺序。
var utFiles = []string{"template.json", "index.html", "style.css", "layouts.md", "rules.md"}

// snapshotUT 五件套整包 bundle 的磁盘形状（<version>.json）。
// 与 deck 同一条设计理由：记档/恢复都是单文件原子操作，不存在拷一半的中间态。
type snapshotUT struct {
	TemplateJSON string `json:"template_json,omitempty"`
	IndexHTML    string `json:"index_html,omitempty"`
	StyleCSS     string `json:"style_css,omitempty"`
	LayoutsMD    string `json:"layouts_md,omitempty"`
	RulesMD      string `json:"rules_md,omitempty"`
}

func (a snapshotUT) equal(b snapshotUT) bool {
	return a.TemplateJSON == b.TemplateJSON && a.IndexHTML == b.IndexHTML &&
		a.StyleCSS == b.StyleCSS && a.LayoutsMD == b.LayoutsMD && a.RulesMD == b.RulesMD
}

// changedFiles 与上一版逐文件比对，返回内容有出入的文件名（首个 bundle 全量非空文件）。
func (a snapshotUT) changedFiles(b snapshotUT) []string {
	var out []string
	for _, f := range utFiles {
		if a.fileOf(f) != b.fileOf(f) {
			out = append(out, f)
		}
	}
	return out
}

func (b snapshotUT) fileOf(name string) string {
	switch name {
	case "template.json":
		return b.TemplateJSON
	case "index.html":
		return b.IndexHTML
	case "style.css":
		return b.StyleCSS
	case "layouts.md":
		return b.LayoutsMD
	case "rules.md":
		return b.RulesMD
	}
	return ""
}

// UTVersionMeta 一条版本的元信息（index.json 的 versions 数组元素，新→旧）。
type UTVersionMeta struct {
	Version   string   `json:"version"`
	Time      int64    `json:"time"`
	Operation string   `json:"operation"`
	Detail    string   `json:"detail"`
	Changed   []string `json:"changed,omitempty"`
}

type utHistoryIndex struct {
	NextSeq  int             `json:"next_seq"`
	Versions []UTVersionMeta `json:"versions"`
}

func (s *Service) historyDir(id string) string {
	return filepath.Join(s.Dir(id), "history")
}

func (s *Service) readUTHistoryIndex(id string) utHistoryIndex {
	data, err := os.ReadFile(filepath.Join(s.historyDir(id), "index.json"))
	if err != nil {
		return utHistoryIndex{NextSeq: 1}
	}
	var idx utHistoryIndex
	if json.Unmarshal(data, &idx) != nil || idx.NextSeq < 1 {
		return utHistoryIndex{NextSeq: 1}
	}
	return idx
}

func (s *Service) writeUTHistoryIndex(id string, idx utHistoryIndex) error {
	if err := os.MkdirAll(s.historyDir(id), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(idx)
	if err != nil {
		return err
	}
	return atomicWriteFile(filepath.Join(s.historyDir(id), "index.json"), data)
}

// readUTBundle 读模板目录现状为 bundle（缺失文件=空串，与 deck 的口径一致）。
func (s *Service) readUTBundle(id string) snapshotUT {
	var b snapshotUT
	for _, f := range utFiles {
		data, err := os.ReadFile(filepath.Join(s.Dir(id), f))
		if err != nil {
			continue
		}
		b.fileSet(f, string(data))
	}
	return b
}

func (s *Service) readUTBundleAt(id, version string) (snapshotUT, bool) {
	data, err := os.ReadFile(filepath.Join(s.historyDir(id), version+".json"))
	if err != nil {
		return snapshotUT{}, false
	}
	var b snapshotUT
	if json.Unmarshal(data, &b) != nil {
		return snapshotUT{}, false
	}
	return b, true
}

func (b *snapshotUT) fileSet(name, content string) {
	switch name {
	case "template.json":
		b.TemplateJSON = content
	case "index.html":
		b.IndexHTML = content
	case "style.css":
		b.StyleCSS = content
	case "layouts.md":
		b.LayoutsMD = content
	case "rules.md":
		b.RulesMD = content
	}
}

// lockUT 模板级互斥（对话与手动保存可能并发写同一模板）。
// 返回解锁函数，与 deck 的 lockDeck 同一形态。
func (s *Service) lockUT(id string) func() {
	s.locksMu.Lock()
	defer s.locksMu.Unlock()
	if s.locks == nil {
		s.locks = map[string]*sync.Mutex{}
	}
	m, ok := s.locks[id]
	if !ok {
		m = &sync.Mutex{}
		s.locks[id] = m
	}
	m.Lock()
	return m.Unlock
}

// recordVersionUT 记一笔版本（锁内调用）。幂等：与上一版完全一致不记。
// 历史是锦上添花，失败只 log 不阻塞主流程（与 deck recordVersionV2 同纪律）。
func (s *Service) recordVersionUT(id, operation, detail string) error {
	idx := s.readUTHistoryIndex(id)
	bundle := s.readUTBundle(id)

	meta := UTVersionMeta{
		Version:   fmt.Sprintf("v%06d", idx.NextSeq),
		Time:      time.Now().Unix(),
		Operation: operation,
		Detail:    detail,
	}
	if len(idx.Versions) > 0 {
		prev, ok := s.readUTBundleAt(id, idx.Versions[0].Version)
		if ok {
			if prev.equal(bundle) {
				return nil // 无实质变化，不产生噪音版本
			}
			meta.Changed = bundle.changedFiles(prev)
		}
	}
	if meta.Changed == nil {
		// 首版没有上一版可比：changed = 全部非空文件（基线的形状）
		for _, f := range utFiles {
			if bundle.fileOf(f) != "" {
				meta.Changed = append(meta.Changed, f)
			}
		}
	}

	dir := s.historyDir(id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Printf("[warn] 模板 %s 历史目录创建失败: %v", id, err)
		return nil
	}
	data, err := json.Marshal(&bundle)
	if err != nil {
		log.Printf("[warn] 模板 %s 快照序列化失败: %v", id, err)
		return nil
	}
	if err := atomicWriteFile(filepath.Join(dir, meta.Version+".json"), data); err != nil {
		log.Printf("[warn] 模板 %s 快照 %s 写入失败: %v", id, meta.Version, err)
		return nil
	}

	idx.NextSeq++
	idx.Versions = append([]UTVersionMeta{meta}, idx.Versions...)
	kept, dropped := pruneUTVersions(idx.Versions)
	for _, v := range dropped {
		_ = os.Remove(filepath.Join(dir, v+".json"))
	}
	idx.Versions = kept
	if err := s.writeUTHistoryIndex(id, idx); err != nil {
		log.Printf("[warn] 模板 %s 历史索引写入失败: %v", id, err)
	}
	return nil
}

// pruneUTVersions 分层裁剪：近 keepFullDays 天全留；更早的每天只留当天最新一条；
// 总数超 maxVersionsPerUT 截尾。起点版本（传入切片的最后一条 = 最早）永不裁——
// "回得去起点"是本特性的承诺，为此允许比上限多留一条。
func pruneUTVersions(versions []UTVersionMeta) (kept []UTVersionMeta, dropped []string) {
	if len(versions) == 0 {
		return nil, nil
	}
	oldest := versions[len(versions)-1]
	cutoff := time.Now().AddDate(0, 0, -keepFullDays)
	dayOf := func(t int64) string { return time.Unix(t, 0).Format("2006-01-02") }

	seenOldDays := make(map[string]bool)
	for _, m := range versions {
		recent := time.Unix(m.Time, 0).After(cutoff)
		if recent || !seenOldDays[dayOf(m.Time)] || m.Version == oldest.Version {
			kept = append(kept, m)
			if !recent {
				seenOldDays[dayOf(m.Time)] = true
			}
		} else {
			dropped = append(dropped, m.Version)
		}
	}
	if len(kept) > maxVersionsPerUT {
		tail := kept[maxVersionsPerUT:]
		keepOldest := false
		for _, m := range tail {
			if m.Version == oldest.Version {
				keepOldest = true
				continue
			}
			dropped = append(dropped, m.Version)
		}
		kept = kept[:maxVersionsPerUT]
		if keepOldest {
			kept = append(kept, oldest) // 比上限多一条，是起点版本的豁免位
		}
	}
	return kept, dropped
}

// ListUTVersions 某模板的版本历史（新→旧）。
func (s *Service) ListUTVersions(userID uint, id string) ([]UTVersionMeta, error) {
	if _, err := s.GetOwned(userID, id); err != nil {
		return nil, err
	}
	idx := s.readUTHistoryIndex(id)
	if idx.Versions == nil {
		idx.Versions = []UTVersionMeta{}
	}
	return idx.Versions, nil
}

// RestoreUTVersion 回滚到指定版本：bundle 整包写回 + 元数据同步 DB + 重挂注册表 +
// 记一条 restore 版本。published 拒绝（D2：公开内容不能经回滚悄悄改变）。
func (s *Service) RestoreUTVersion(userID uint, id, version string) error {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return err
	}
	if row.Status == "published" {
		return ErrPublished
	}
	if !utVersionPattern.MatchString(version) {
		return fmt.Errorf("版本号 %q 不合法", version)
	}

	unlock := s.lockUT(id)
	defer unlock()

	idx := s.readUTHistoryIndex(id)
	found := false
	for _, m := range idx.Versions {
		if m.Version == version {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("版本 %s 不存在（可能已被删除或裁剪）", version)
	}
	bundle, ok := s.readUTBundleAt(id, version)
	if !ok {
		return fmt.Errorf("快照 %s 读取失败或已损坏", version)
	}

	for _, f := range utFiles {
		if content := bundle.fileOf(f); content != "" {
			if err := atomicWriteFile(filepath.Join(s.Dir(id), f), []byte(content)); err != nil {
				return fmt.Errorf("写回 %s 失败: %w", f, err)
			}
		}
	}
	// 元数据以快照为准同步回 DB（展示名与 template.json 保持单一真相）
	var meta struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal([]byte(bundle.TemplateJSON), &meta); err == nil {
		updates := map[string]any{}
		if meta.Name != "" {
			updates["name"] = meta.Name
		}
		if meta.Description != "" {
			updates["description"] = meta.Description
		}
		if len(updates) > 0 {
			if err := s.st.DB.Model(&store.UserTemplate{}).Where("id = ?", id).Updates(updates).Error; err != nil {
				return fmt.Errorf("同步元数据失败: %w", err)
			}
		}
	}
	// 重挂注册表：恢复后的模板必须整体有效；失败时报错但文件已写回，
	// 用户可再回滚一版（历史不受影响）。
	if s.reg != nil {
		if err := s.reg.MountUser(s.Dir(id)); err != nil {
			return fmt.Errorf("恢复后的模板未通过注册表校验：%w（可再回滚到其他版本）", err)
		}
	}
	return s.recordVersionUT(id, OpRestore, "恢复到 "+version)
}

// DeleteUTVersion 删除单条历史（NextSeq 不清零，编号永不复用）。
func (s *Service) DeleteUTVersion(userID uint, id, version string) error {
	if _, err := s.GetOwned(userID, id); err != nil {
		return err
	}
	if !utVersionPattern.MatchString(version) {
		return fmt.Errorf("版本号 %q 不合法", version)
	}
	unlock := s.lockUT(id)
	defer unlock()

	idx := s.readUTHistoryIndex(id)
	kept := make([]UTVersionMeta, 0, len(idx.Versions))
	found := false
	for _, m := range idx.Versions {
		if m.Version == version {
			found = true
			continue
		}
		kept = append(kept, m)
	}
	if !found {
		return fmt.Errorf("版本 %s 不存在", version)
	}
	if err := os.Remove(filepath.Join(s.historyDir(id), version+".json")); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除快照失败: %w", err)
	}
	idx.Versions = kept
	return s.writeUTHistoryIndex(id, idx)
}

// ClearUTHistory 清空全部历史（起点版本也删——用户显式要求；NextSeq 不清零）。
func (s *Service) ClearUTHistory(userID uint, id string) (int, error) {
	if _, err := s.GetOwned(userID, id); err != nil {
		return 0, err
	}
	unlock := s.lockUT(id)
	defer unlock()

	idx := s.readUTHistoryIndex(id)
	n := len(idx.Versions)
	for _, m := range idx.Versions {
		_ = os.Remove(filepath.Join(s.historyDir(id), m.Version+".json"))
	}
	idx.Versions = []UTVersionMeta{}
	return n, s.writeUTHistoryIndex(id, idx)
}
