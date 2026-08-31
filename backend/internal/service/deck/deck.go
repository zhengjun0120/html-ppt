package deck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// idPattern 是 deck id 的白名单。所有来自外部（URL、LLM 工具调用）的 id
// 都是"不可信输入"：不校验就拼路径，攻击者可以传 ../../ 绕出数据目录
// （路径穿越攻击）。白名单正则是这里唯一的安全闸门。
var idPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// Meta 是列表接口返回的 deck 摘要。
type Meta struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Service 负责 deck 文件的存取。目前只有读（列表 + 单个），
// 阶段2 加入页级读写时，原子写入 / 版本快照逻辑也放这里。
type Service struct {
	decksDir string
}

func New(dataDir string) *Service {
	return &Service{decksDir: filepath.Join(dataDir, "decks")}
}

// List 扫描数据目录，返回所有 deck 的 id 和标题。
func (s *Service) List() ([]Meta, error) {
	entries, err := os.ReadDir(s.decksDir)
	if os.IsNotExist(err) {
		return []Meta{}, nil // 目录还没建过 = 空列表，不是错误
	}
	if err != nil {
		return nil, err
	}
	var out []Meta
	for _, e := range entries {
		if !e.IsDir() || !idPattern.MatchString(e.Name()) {
			continue
		}
		html, err := s.GetHTML(e.Name())
		if err != nil {
			continue // 坏了一个 deck 不影响整体列表
		}
		out = append(out, Meta{ID: e.Name(), Title: extractTitle(html)})
	}
	return out, nil
}

// FilePath 校验 id 后返回 deck.html 的绝对路径。
// 只暴露路径给内部使用，LLM 的工具签名里永远只出现 id（见阶段2工具设计）。
func (s *Service) FilePath(id string) (string, error) {
	if !idPattern.MatchString(id) {
		return "", fmt.Errorf("invalid deck id: %q", id)
	}
	return filepath.Join(s.decksDir, id, "deck.html"), nil
}

func (s *Service) GetHTML(id string) (string, error) {
	p, err := s.FilePath(id)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// extractTitle 取 <title> 文本当标题。
// 阶段2会换成 goquery 逐页解析（切 section、算指纹），这里先够用。
func extractTitle(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(doc.Find("title").First().Text())
}
