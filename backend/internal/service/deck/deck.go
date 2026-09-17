package deck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
)

// idPattern 是 deck id 的白名单。所有来自外部（URL、LLM 工具调用）的 id
// 都是"不可信输入"：不校验就拼路径，攻击者可以传 ../../ 绕出数据目录
// （路径穿越攻击）。白名单正则是这里唯一的安全闸门。
var idPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

// IsValidID 是全项目唯一的"合法 id"判定。URL、LLM 工具参数注入前都过这里。
func IsValidID(id string) bool {
	return idPattern.MatchString(id)
}

// Meta 是列表接口返回的 deck 摘要。
type Meta struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Service 负责 deck 文件的存取 + 归属校验。
// 归属的权威在 decks 表：文件系统只是内容存储，"这个 deck 是谁的"只认 DB。
type Service struct {
	decksDir string
	assetsDir string
	st        *store.Store // nil = 数据库降级模式：所有操作返回不可用
	deckLocks sync.Map
	// templates deck-v2 模板注册表（WithTemplateRegistry 挂载；nil = v2 写路径不可用）
	templates *template.Registry
}

func New(dataDir, assetsDir string, st *store.Store) *Service {
	return &Service{decksDir: filepath.Join(dataDir, "decks"), assetsDir: assetsDir, st: st}
}

var errStorage = fmt.Errorf("数据库不可用")

// authorize 是每个公开方法的第一个动作：id 白名单 + 归属校验。
// 越权和不存在返回同一个错误——不泄露"这个 deck 存在但不是你的"。
func (s *Service) authorize(userID uint, id string) error {
	if s.st == nil {
		return errStorage
	}
	if !idPattern.MatchString(id) {
		return fmt.Errorf("deck %q 不存在", id)
	}
	var cnt int64
	if err := s.st.DB.Model(&store.Deck{}).
		Where("id = ? AND user_id = ?", id, userID).Count(&cnt).Error; err != nil {
		return fmt.Errorf("查询 deck 归属: %w", err)
	}
	if cnt == 0 {
		return fmt.Errorf("deck %q 不存在", id)
	}
	return nil
}

// EnsureOwner 归属校验的公开薄封装，供 chat 入口预校验注入的 deck_id
//（在把 deck_id 写进 system prompt 之前就拦住无主/越权的 id）。
func (s *Service) EnsureOwner(userID uint, id string) error {
	return s.authorize(userID, id)
}

// List 当前用户的 deck 列表。D7：只列 v2（format=v2），旧格式行隐藏不迁移。
func (s *Service) List(userID uint) ([]Meta, error) {
	if s.st == nil {
		return nil, errStorage
	}
	var rows []store.Deck
	if err := s.st.DB.Where("user_id = ? AND format = ?", userID, FormatV2).Order("id").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("查询 deck 列表: %w", err)
	}
	out := make([]Meta, 0, len(rows))
	for _, r := range rows {
		out = append(out, Meta{ID: r.ID, Title: r.Title})
	}
	return out, nil
}

// GetHTML 读取 deck 源文件（归属校验后的公开读入口）。
// v2 唯一格式：index.html。量测/审查/渲染通道共用这一个读入口，
// 保证"审查看到的页面"和"用户看到的页面"永远来自同一个文件。
func (s *Service) GetHTML(userID uint, id string) (string, error) {
	if err := s.authorize(userID, id); err != nil {
		return "", err
	}
	return s.readIndex(id)
}

// readIndex 读 index.html（不做归属校验——只允许 authorize 之后的内部调用）。
func (s *Service) readIndex(id string) (string, error) {
	p, err := s.IndexPathV2(id)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("deck %s 尚未实例化（index.html 缺失，先选择模板）", id)
	}
	return string(data), nil
}

// IsV2 以 deck.json 的存在为准（DB Format 列是权威，但这里只做读取分流，
// 用文件判断避免为纯读取路径打一次 DB）。v1 格式已随旧栈摘除，非 v2 一律 false。
func (s *Service) IsV2(id string) bool {
	_, err := os.Stat(s.deckFilePath(id))
	return err == nil
}

// NormalizeTitle 列表/会话共用的标题截断（超长截 30 字）。
func NormalizeTitle(title string) string {
	title = strings.TrimSpace(title)
	if r := []rune(title); len(r) > 30 {
		return string(r[:30])
	}
	return title
}

func (s *Service) lockDeck(deckID string) func() {
	v, _ := s.deckLocks.LoadOrStore(deckID, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}
