package usertpl

// 从空白新建（用户需求 2026-09-28）：不派生自任何内置模板，直接落一份能通过
// 注册表校验的最小脚手架——中性 token、两个最小说明版式（封面/内容页）。
// 用户随后在定制工作台里用对话（write_style/write_demo）、样式面板、手动编辑
// 把它长成自己的模板；历史系统从 v1（空白基线）开始兜底。
//
// 脚手架经 go:embed 进二进制：不进 templates/ 目录（NewRegistry 会把那里的
// 每个子目录当内置模板加载，空白模板不该出现在画廊/选模板里）。

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"html-ppt/backend/internal/store"
)

//go:embed scaffold/*.json scaffold/*.html scaffold/*.css scaffold/*.md
var scaffoldFS embed.FS

// scaffoldFiles 与 forkCopies 同一套（ADAPTATION.md 由 CreateBlank 现写）。
var scaffoldFiles = []string{"template.json", "index.html", "style.css", "layouts.md", "rules.md"}

// CreateBlank 从内置空白脚手架新建一份私有模板。name 为空用「空白模板」。
// 与 Fork 同一条纪律：落盘 → 建行 → MountUser（校验失败连目录一起清）→ 记起点版本。
func (s *Service) CreateBlank(userID uint, name string) (*store.UserTemplate, error) {
	if strings.TrimSpace(name) == "" {
		name = "空白模板"
	}
	id, err := newID()
	if err != nil {
		return nil, err
	}
	dir := s.Dir(id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	for _, f := range scaffoldFiles {
		content, err := scaffoldFS.ReadFile("scaffold/" + f)
		if err != nil {
			_ = os.RemoveAll(dir)
			return nil, fmt.Errorf("脚手架缺 %s: %w", f, err)
		}
		if f == "template.json" {
			c := string(content)
			c = strings.ReplaceAll(c, "{{ID}}", id)
			c = strings.ReplaceAll(c, "{{NAME}}", name)
			content = []byte(c)
		}
		if err := os.WriteFile(filepath.Join(dir, f), content, 0o644); err != nil {
			_ = os.RemoveAll(dir)
			return nil, fmt.Errorf("写 %s 失败: %w", f, err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "ADAPTATION.md"), []byte(fmt.Sprintf(
		"# %s\n\n- 从空白脚手架新建（不派生自任何内置模板），可自由定制。\n"+
			"- 定制入口：对话改 token / 重写 style.css 与 demo 结构；结构契约（版式/类名）在 layouts.md 登记后生效。\n- 创建时间：%s\n",
		name, time.Now().Format("2006-01-02 15:04"))), 0o644); err != nil {
		_ = os.RemoveAll(dir)
		return nil, err
	}

	row := &store.UserTemplate{
		ID: id, UserID: userID, BaseID: "",
		Name: name, Description: "从空白新建的模板：中性灰阶起点，等你定义",
		Visibility: "private", Status: "draft",
	}
	if err := s.st.DB.Create(row).Error; err != nil {
		_ = os.RemoveAll(dir)
		return nil, err
	}
	if err := s.reg.MountUser(dir); err != nil {
		_ = s.st.DB.Delete(row)
		_ = os.RemoveAll(dir)
		return nil, fmt.Errorf("空白脚手架没过校验: %w", err)
	}
	_ = s.recordVersionUT(id, OpFork, "从空白新建")
	return row, nil
}
