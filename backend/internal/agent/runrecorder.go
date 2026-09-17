package agent

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type runRecorder struct {
	mu sync.Mutex
	decks map[string][]string
	// toolCalls 本 run 内各工具的执行次数（含被配额拒掉的），配额检查用。
	// runRecorder 每 run 一个，正好是配额的作用域；execTool 逐个顺序执行工具，
	// 没有并发，锁沿用这里现成的 mu，不为它单独引入原子类型。
	toolCalls map[string]int
}

func newRunRecorder() *runRecorder{
	return &runRecorder{
		decks: make(map[string][]string),
	}
}

// overQuota 记一次工具调用并判断是否超出配额。max 为 0 或负数表示不限。
// 先计数再比较：第 1~max 次放行执行，第 max+1 次起被拒。
func (rr *runRecorder) overQuota(name string, max int) bool {
	if max <= 0 {
		return false
	}
	rr.mu.Lock()
	defer rr.mu.Unlock()
	if rr.toolCalls == nil {
		rr.toolCalls = make(map[string]int)
	}
	rr.toolCalls[name]++
	return rr.toolCalls[name] > max
}

// refund 把一次已计入的调用退回去（下限钳到 0）。存在的理由：配额的本意是拦
// "审查→修复→再审"的无限循环，拦的是**烧了钱、产出了报告**的那类调用。
// review_slides 的"只看数字"复查（pages 留空）和"看图失败"的调用什么结论都没产出，
// 让它们占配额的结果是（实测）：两次免费复查 + 一次看图失败 = 配额耗尽，
// agent 想重试真审查时被拒——工具说明承诺"pages 留空不占配额"，代码必须兑现它。
func (rr *runRecorder) refund(name string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	if rr.toolCalls[name] > 0 {
		rr.toolCalls[name]--
	}
}

func(rr *runRecorder) note(deckID,op string){
	rr.mu.Lock()
	defer rr.mu.Unlock()
	rr.decks[deckID] = append(rr.decks[deckID], op)
}

func(rr *runRecorder)noteToolCall(name,argsJSON,resultJSON string){
	switch name{
	case "write_deck":
		var res struct {
			DeckID string `json:"deck_id"`
			Slides int `json:"slides"`
		}
		if json.Unmarshal([]byte(resultJSON),&res) == nil && res.DeckID!=""{
			rr.note(res.DeckID,fmt.Sprintf("创建(%d页)",res.Slides))
		}
	case "update_slide":
		var a struct {
			DeckID string `json:"deck_id"`
			SlideID string `json:"slide_id"`
		}
		if json.Unmarshal([]byte(argsJSON),&a) == nil && a.DeckID!= ""{
			rr.note(a.DeckID,"修改了 "+a.SlideID)
			
		}
	case "insert_slide":
		var a struct {
			DeckID string `json:"deck_id"`
			AfterSlideID string `json:"after_slide_id"`
		}
		if json.Unmarshal([]byte(argsJSON),&a) == nil && a.DeckID != ""{
			var res struct {
				SlideID string `json:"slide_id"`
			}
			if json.Unmarshal([]byte(resultJSON),&res) == nil && res.SlideID != ""{
				rr.note(a.DeckID,"插入了新页 "+res.SlideID)
			}
		}
	case "delete_slide":
		var a struct {
			DeckID string `json:"deck_id"`
			SlideID string `json:"slide_id"`
		}
		if json.Unmarshal([]byte(argsJSON),&a) == nil && a.DeckID !=""{
			rr.note(a.DeckID,"删除了 "+a.SlideID)
		}
	case "update_theme":
		var res struct {
			DeckID string `json:"deck_id"`
		}
		if json.Unmarshal([]byte(resultJSON),&res) == nil && res.DeckID!= ""{
			rr.note(res.DeckID,"调整主题")
		}
	}
}

func (rr *runRecorder) deckIDs() []string{
	rr.mu.Lock()
	defer rr.mu.Unlock()
	ids := make([]string,0,len(rr.decks))
	for id := range rr.decks{
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func(rr *runRecorder) detail(deckID string) string{
	rr.mu.Lock()
	defer rr.mu.Unlock()
	return strings.Join(rr.decks[deckID],"、")
}