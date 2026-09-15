package vision

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type Grants struct {
	mu sync.Mutex
	open map[string]grant
}

type grant struct{
	uid uint
	deckID string
	expires time.Time
}

//过期时间
const grantTTL = 2 *time.Minute

//颁发门票
func(g *Grants) Issue(uid uint,deckID string)(string,error){
	buf := make([]byte,24)
	if _,err := rand.Read(buf);err != nil{
		return "",err
	}
	nonce := hex.EncodeToString(buf)
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.open== nil{
		g.open = map[string]grant{}
	}

	for k,v := range g.open{
		if time.Now().After(v.expires){
			delete(g.open,k)
		}
	}
	g.open[nonce] = grant{uid:uid,deckID: deckID,expires: time.Now().Add(grantTTL)}
	return nonce,nil
}

// 消费门票，如果有效则返回uid和deckid
func (g *Grants) Take(nonce string) (uid uint,deckID string,ok bool){
	g.mu.Lock()
	defer g.mu.Unlock()
	v,found := g.open[nonce]
	if !found{
		return 0,"",false
	}
	delete(g.open,nonce)
	//门票过期
	if time.Now().After(v.expires){
		return 0,"",false
	}
	return v.uid,v.deckID,true
}

