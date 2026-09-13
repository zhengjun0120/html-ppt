package deck

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

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
	// assetsDir 共享静态资源目录（web/assets）：回声校验从这里现扫变量契约。
	// 只读——组件库对 AI 和对后端都是只读的（护栏）。
	assetsDir string
	st        *store.Store // nil = 数据库降级模式：所有操作返回不可用
	deckLocks sync.Map
}

func New(dataDir, assetsDir string, st *store.Store) *Service {
	return &Service{decksDir: filepath.Join(dataDir, "decks"), assetsDir: assetsDir, st: st}
}

// variableContract 取"被框架 CSS 消费过"的变量清单（详见 variable_contract.go）。
// 拿不到契约时返回 nil 并跳过回声校验（fail-open）：这项检查防的是"静默无效的样式"，
// 不是安全问题——宁可漏判，也不能因为资源目录没配好就把所有写入都卡死。
func (s *Service) variableContract() map[string]bool {
	known, err := loadVariableContract(s.assetsDir)
	if err != nil {
		log.Printf("[warn] 读取 CSS 变量契约失败，本次跳过回声校验: %v", err)
		return nil
	}
	return known
}

// checkInlineStyleVars 校验页面内容里内联 style 用到的变量是否都有效。
// known 之外还要加上该 deck 自己在样式块里定义过的变量——否则"自定义槽里定义变量、
// 页面里 var() 使用"这种正常写法会被误判成无效。
func (s *Service) checkInlineStyleVars(sectionHTML, deckHTML string) error {
	known := s.variableContract()
	if known == nil {
		return nil
	}
	defined := definedVariables(deckStyleBlocks(deckHTML))
	if unknown := unknownVariables(collectInlineStyleVars(sectionHTML), known, defined); len(unknown) > 0 {
		return unknownVarErr(unknown, known)
	}
	return nil
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

// List 当前用户的全部 deck（标题在创建时落库，列表是纯 DB 查询，不再读文件）。
func (s *Service) List(userID uint) ([]Meta, error) {
	if s.st == nil {
		return nil, errStorage
	}
	var rows []store.Deck
	if err := s.st.DB.Where("user_id = ?", userID).Order("id").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("查询 deck 列表: %w", err)
	}
	out := make([]Meta, 0, len(rows))
	for _, r := range rows {
		out = append(out, Meta{ID: r.ID, Title: r.Title})
	}
	return out, nil
}

// FilePath 校验 id 白名单后返回 deck.html 的绝对路径。
// 只暴露路径给内部使用；外部输入先过 authorize。
func (s *Service) FilePath(id string) (string, error) {
	if !idPattern.MatchString(id) {
		return "", fmt.Errorf("invalid deck id: %q", id)
	}
	return filepath.Join(s.decksDir, id, "deck.html"), nil
}

// GetHTML 读取整份 deck.html（归属校验后的公开读入口）。
func (s *Service) GetHTML(userID uint, id string) (string, error) {
	if err := s.authorize(userID, id); err != nil {
		return "", err
	}
	return s.readRaw(id)
}

// readRaw 读原文，不做归属校验——只允许 authorize 之后的内部调用。
func (s *Service) readRaw(id string) (string, error) {
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

//  鉴权+锁内读文件
func (s *Service) readOwned(userID uint,id string)(string,error){
	err := s.authorize(userID,id)
	if err !=nil{
		return "",err
	}

	return s.readRaw(id)
}

func(s *Service) lockDeck(deckID string) func(){
	v,_ := s.deckLocks.LoadOrStore(deckID,&sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}
