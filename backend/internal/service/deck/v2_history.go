package deck

// deck-v2 的版本历史：快照对象是三件套 bundle（index.html + style.css + outline.json）。
//
// 为什么不是三个文件：历史目录里一个版本一个文件，恢复是"整包换回"——bundle JSON
// 让"记档"和"恢复"都是单文件原子操作，不存在"拷了一半"的中间态。
// runtime.js/base.css/字体是共享资产（/assets/deck-v2/*，版本化路径逃生舱），不进快照。

import (
	"encoding/json"
	"log"
	"time"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// snapshotV2 v2 快照 bundle 的磁盘形状（<version>.json）。
type snapshotV2 struct {
	IndexHTML  string `json:"index_html"`
	StyleCSS   string `json:"style_css,omitempty"`
	OutlineJSON string `json:"outline_json,omitempty"`
}

// recordVersionV2 记一笔 v2 快照（锁内调用，与 v1 的 recordVersion 同一套索引）。
func (s *Service) recordVersionV2(deckID, operation, detail string) error {
	dir := s.historyDir(deckID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Printf("deck %s 历史目录创建失败 err: %v", deckID, err)
		return nil // 历史是锦上添花，不阻塞主流程——与 v1 recordVersion 同一取舍
	}

	idx := s.readHistoryIndex(deckID)

	bundle := snapshotV2{}
	var err error
	if bundle.IndexHTML, err = s.readV2File(deckID, "index.html"); err != nil {
		return err
	}
	// style.css / outline.json 允许缺失（实例化前后的阶段不同）
	bundle.StyleCSS, _ = s.readV2File(deckID, "style.css")
	bundle.OutlineJSON, _ = s.readV2File(deckID, "outline.json")

	data, err := json.Marshal(&bundle)
	if err != nil {
		return err
	}

	meta := VersionMeta{
		Version:   fmt.Sprintf("v%06d", idx.NextSeq),
		Time:      time.Now().Unix(),
		Operation: operation,
		Detail:    detail,
		Slides:    countSlidesV2(bundle.IndexHTML),
	}
	if err := atomicWriteFile(filepath.Join(dir, meta.Version+".json"), data); err != nil {
		log.Printf("deck %s 快照 %s 写入失败 err:%v", deckID, meta.Version, err)
		return nil
	}
	idx.NextSeq++
	idx.Versions = append([]VersionMeta{meta}, idx.Versions...)

	kept, dropped := pruneVersions(idx.Versions)
	for _, v := range dropped {
		// v2 快照是 .json；兼容清理可能存在的 v1 残留 .html
		_ = os.Remove(filepath.Join(dir, v+".json"))
		_ = os.Remove(filepath.Join(dir, v+".html"))
	}
	idx.Versions = kept
	if err := s.writeHistoryIndex(deckID, idx); err != nil {
		log.Printf("deck %s 历史索引写入失败 err: %v", deckID, err)
	}
	return nil
}

// readV2File 读 deck 目录内的文件（不存在返回空串 + err）。
func (s *Service) readV2File(deckID, name string) (string, error) {
	data, err := os.ReadFile(filepath.Join(s.decksDir, deckID, name))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// RestoreVersionV2 v2 的版本恢复：bundle 整包写回 + 立即记一条 restore。
// 阶段不回退——恢复的是"内容"，不是流程状态机（流程只会向前）。
func (s *Service) RestoreVersionV2(userID uint, deckID, version string) error {
	if err := s.authorize(userID, deckID); err != nil {
		return err
	}
	if !versionPattern.MatchString(version) {
		return fmt.Errorf("版本号 %q 不合法", version)
	}
	unlock := s.lockDeck(deckID)
	defer unlock()

	idx := s.readHistoryIndex(deckID)
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

	data, err := os.ReadFile(filepath.Join(s.historyDir(deckID), version+".json"))
	if err != nil {
		// v1 快照（.html）不能用 v2 恢复：格式不同，报清楚而不是给一半
		if _, err1 := os.Stat(filepath.Join(s.historyDir(deckID), version+".html")); err1 == nil {
			return fmt.Errorf("版本 %s 是旧格式快照，v2 deck 无法恢复它", version)
		}
		return fmt.Errorf("快照文件读取失败 err: %w", err)
	}
	var bundle snapshotV2
	if err := json.Unmarshal(data, &bundle); err != nil {
		return fmt.Errorf("快照损坏 err: %w", err)
	}

	ip, err := s.IndexPathV2(deckID)
	if err != nil {
		return err
	}
	if err := atomicWriteFile(ip, []byte(bundle.IndexHTML)); err != nil {
		return fmt.Errorf("写回 index.html 失败: %w", err)
	}
	if bundle.StyleCSS != "" {
		if err := atomicWriteFile(s.stylePath(deckID), []byte(bundle.StyleCSS)); err != nil {
			return fmt.Errorf("写回 style.css 失败: %w", err)
		}
	}
	if bundle.OutlineJSON != "" {
		if err := atomicWriteFile(s.outlineFilePath(deckID), []byte(bundle.OutlineJSON)); err != nil {
			return fmt.Errorf("写回 outline.json 失败: %w", err)
		}
	}
	return s.recordVersionV2(deckID, OpRestore, "恢复到 "+version)
}

// countSlidesV2 v2 的页数统计：.deck 直接子级 section。
func countSlidesV2(html string) int {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("解析 index.html 失败 err:%v", err)
		return 0
	}
	return doc.Find(".deck > section").Length()
}
