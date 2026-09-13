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
}

func newRunRecorder() *runRecorder{
	return &runRecorder{
		decks: make(map[string][]string),
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