package deck

import "fmt"

// v2 手动编辑保存链路（docs/deck-editor-plan.md §4.1）。
//
// 与 agent 写入（write_pages → FinishGeneration → RecordRunVersion）走同一条
// 「锁内写 index.html → 记版本 → 缩略图失效」纪律：手动保存的 HTML 就是新的
// 当前内容，快照 bundle 打包的是刚落盘的文件，恢复链路（RestoreVersionV2）
// 对 edit 版本与 run 版本一视同仁。stage 不动（编辑不回退流程状态机）。

// SaveHTML 编辑器全量保存：覆盖 index.html 并记一条 OpEdit 版本。
// detail 进版本历史的人读汇总，空则给默认文案。
func (s *Service) SaveHTML(userID uint, id, html, detail string) error {
	if err := s.authorize(userID, id); err != nil {
		return err
	}
	if !s.IsV2(id) {
		return fmt.Errorf("deck %s 不支持编辑（仅 v2 格式）", id)
	}
	unlock := s.lockDeck(id)
	defer unlock()

	ip, err := s.IndexPathV2(id)
	if err != nil {
		return err
	}
	if err := atomicWriteFile(ip, []byte(html)); err != nil {
		return fmt.Errorf("写回 index.html 失败: %w", err)
	}
	s.invalidateThumbs(id)
	if detail == "" {
		detail = "手动编辑"
	}
	return s.recordVersionV2(id, OpEdit, detail)
}
