package deck

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	maxDiffLinesPerSlide  = 80
	maxDiffSlidesDetailed = 6
	// maxDiffLineChars 单行字符上限。只按行数截断挡不住"整页一行"的写法：
	// AI 把 section 全写在一行时（很常见），一次改动就是一条几十 KB 的长行，
	// 行数上限等于没截到。超长行按字符截断并留 …（要看全文用 read_slide）。
	maxDiffLineChars = 240
)


type SlideDiff struct {
	SlideID  string `json:"slide_id"`
	Title    string `json:"title"`
	Change   string `json:"change"`
	Diff     string `json:"diff,omitempty"`
	Truncated bool   `json:"truncated,omitempty"`
}

type HistoryDiff struct {
	DeckID    string      `json:"deck_id"`
	From      string      `json:"from"` //旧版本号
	To        string      `json:"to"`// 新版本号 或 “current”
	ToDetail  string      `json:"to_detail"` //to 是版本号时的操作汇总；current时为空
	Changed bool `json:"changed"` //三个列表全空时为false 模型可直接判断 没差别
	Added     []SlideDiff `json:"added"`
	Removed   []SlideDiff `json:"removed"`
	Modified  []SlideDiff `json:"modified"`
	Unchanged int         `json:"unchanged"`
}

func (s *Service) ReadVersionDiff(userID uint, deckID, fromVersion,toVersion string) (*HistoryDiff, error) {
	if err := s.authorize(userID,deckID);err !=nil{
		return nil,err
	}

	return s.versionDiff(deckID,fromVersion,toVersion)
}

func (s *Service) versionDiff(deckID,fromVersion,toVersion string)(*HistoryDiff,error){
	unlock := s.lockDeck(deckID)
	defer unlock()

	if fromVersion == toVersion{
		if fromVersion!=""&&toVersion!=""{
			return nil,fmt.Errorf("基线和目标是同一个版本 %s",fromVersion)
		}
	}

	//读索引
	idx := s.readHistoryIndex(deckID)
	if len(idx.Versions) == 0 {
		return nil, fmt.Errorf("还没有历史版本，暂时无法比较")
	}

	from := idx.Versions[0]
	if fromVersion != ""{
		m,err := lookupVersion(idx,fromVersion)
		if err != nil{
			return nil,err
		}
		from = m
	}

	fromHTML, err := s.readSnapshotHTML(deckID, from.Version)
	if err != nil {
		return nil, fmt.Errorf("读取快照 %s 失败 err:%w", from.Version, err)
	}

	toLabel := "current"
	var toDetail string
	var toHTML string
	if toVersion == ""{
		toHTML ,err = s.readIndex(deckID)
		if err != nil{
			return nil,err
		}
	}else {
		m,err := lookupVersion(idx,toVersion)
		if err !=nil{
			return nil,err
		}
		b,err := s.readSnapshotHTML(deckID, m.Version)
		if err !=nil{
			return nil,fmt.Errorf("读取快照 %s 失败 err: %w",m.Version,err)
		}
		toLabel,toDetail,toHTML = m.Version,m.Detail,b
	}

	if toLabel == from.Version{
		return nil,fmt.Errorf("基线和目标是同一个版本 %s",from.Version)
	}

	return diffDecks(deckID,from.Version,toLabel,toDetail,string(fromHTML),toHTML),nil

	
}

// readSnapshotHTML 读取某版本快照里的 index.html。v2 快照是 .json bundle
// （index_html 内嵌其中），v1 是裸 .html——按磁盘上实际存在的格式读。
// 此前硬读 .html，v2 deck 的历史目录里根本没有这个文件，diff 工具
// 每次必失败（"读取快照 v000002 失败"）。
func (s *Service) readSnapshotHTML(deckID, version string) (string, error) {
	dir := s.historyDir(deckID)
	if b, err := os.ReadFile(filepath.Join(dir, version+".json")); err == nil {
		var snap snapshotV2
		if err := json.Unmarshal(b, &snap); err != nil {
			return "", fmt.Errorf("快照 %s 解析失败 err:%w", version, err)
		}
		return snap.IndexHTML, nil
	}
	b, err := os.ReadFile(filepath.Join(dir, version+".html"))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// lookupVersion 校验版本号格式并在索引中查找
func lookupVersion(idx historyIndex,version string)(VersionMeta,error){
	if !versionPattern.MatchString(version){
		return VersionMeta{},fmt.Errorf("版本号 %q 不合法",version)
	}
	for _,m := range idx.Versions{
		if m.Version == version{
			return m,nil
		}
	}
	return VersionMeta{},fmt.Errorf("版本 %s 不存在(可能已被删除或裁剪)",version)
}

func diffDecks(deckID,fromLabel,toLabel,toDetail,fromHTML,toHTML string) *HistoryDiff{
	out := &HistoryDiff{
		DeckID: deckID,
		From: fromLabel,
		To: toLabel,
		ToDetail: toDetail,
		Added: []SlideDiff{},
		Removed: []SlideDiff{},
		Modified: []SlideDiff{},
	}

	oldMap,oldOrder := mapSlides(fromHTML)
	newMap,newOrder := mapSlides(toHTML)

	detailed := 0
	for _,id := range newOrder{
		ns := newMap[id]
		os,ok := oldMap[id]
		// 旧的没有新的有
		if !ok{
			out.Added = append(out.Added, SlideDiff{SlideID: id,Title: ns.title,Change: "added"})
			continue
		}

		if os.html == ns.html {
			out.Unchanged ++
			continue
		}

		sd := SlideDiff{SlideID: id,Title: ns.title,Change: "modified"}
		if detailed < maxDiffSlidesDetailed {
			d,trunc := unifiedDiff(os.html,ns.html,maxDiffLinesPerSlide)
			sd.Diff = d
			sd.Truncated = trunc
			detailed++
		}
		out.Modified = append(out.Modified, sd)
	}

	for _,id := range oldOrder {
		if _,ok := newMap[id]; !ok{
			out.Removed =append(out.Removed, SlideDiff{SlideID: id,Title: oldMap[id].title,Change: "removed"})
		}
	}

	out.Changed = len(out.Added)  >0||len(out.Removed)>0||len(out.Modified)>0

	return out
}

func splitKeepNewline(s string) []string{
	if s == ""{
		return nil
	}

	parts := strings.SplitAfter(s,"\n")
	if parts[len(parts)-1] == ""{
		parts = parts[:len(parts)-1]
	}
	return parts
}

func unifiedDiff(before,after string,maxLines int)(string,bool){
	dmp := diffmatchpatch.New()
	c1,c2,lineTexts := dmp.DiffLinesToChars(before,after)
	diffs := dmp.DiffCharsToLines(dmp.DiffMain(c1,c2,false),lineTexts)

	type line struct{prefix,text string}
	var all []line
	for _,d := range diffs {
		for _,l := range splitKeepNewline(d.Text){
			var prefix string
			switch d.Type {
			case diffmatchpatch.DiffEqual:
				prefix = " "
			case diffmatchpatch.DiffDelete:
				prefix="-"
			case diffmatchpatch.DiffInsert:
				prefix="+"
			}
			all = append(all, line{prefix: prefix,text: strings.TrimSuffix(l,"\n")})

		}
	}

	truncated := false
	if len(all) > maxLines {
		all = all[:maxLines]
		truncated = true
	}

	aCount,bCount := 0,0
	for _,l := range all {
		switch l.prefix{
		case "-":
			aCount++
		case "+":
			bCount++
		default:
			aCount++
			bCount++
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb,"@@ -1,%d +1,%d @@",aCount,bCount)
	for _,l := range all{
		sb.WriteString("\n")
		sb.WriteString(l.prefix)
		sb.WriteString(clipLine(l.text,maxDiffLineChars))
	}

	if truncated {
		sb.WriteString("\n...")
	}

	return sb.String(),truncated
}

// clipLine 单行超长时按字符截断（按 rune 切，避免切坏多字节字符），留 … 提示还有内容。
func clipLine(s string, n int) string {
	if len(s) <= n {
		return s
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

type slideSnapshot struct {
	id string
	title string
	html string
}

func mapSlides(html string)(map[string]slideSnapshot,[]string){
	m := map[string]slideSnapshot{}
	order := []string{}
	doc,err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err !=nil{
		log.Printf("解析html失败 err:%v",err)
		return m,order
	}
	doc.Find(".deck > section, .slides > section").Each(func(i int,sec *goquery.Selection){
		// v2 页面挂在 .deck 下，v1 是 .slides；缺 data-id 时用序号兜底，
		// 避免"整页重写没带 id"的页静默消失在 diff 之外
		id := sec.AttrOr("data-id","")
		if id == ""{
			id = fmt.Sprintf("#%d", i+1)
		}
		outer,err := goquery.OuterHtml(sec)
		if err !=nil{
			return 
		}
		m[id] = slideSnapshot{id: id,title:slideTitle(sec),html:outer}
		order = append(order, id)
	})

	return m,order
}