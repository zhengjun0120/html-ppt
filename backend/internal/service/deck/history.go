package deck

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	maxVersionsPerDeck = 200
	keepFullDays       = 7
)

const (
	OpRun     = "run"
	OpRestore = "restore"
)

var versionPattern = regexp.MustCompile(`^v\d{6}$`)


type VersionMeta struct {
	Version string `json:"version"` //v000001 起，单调递增，永不复用
	Time int64 `json:"time"` //unix 秒
	Operation string `json:"operation"` //Op*
	Detail string `json:"detail"` //本轮操作汇总，人读的
	Slides int `json:"slides"` // 快照时的页数
}

type historyIndex struct {
	NextSeq int `json:"next_seq"`
	Versions []VersionMeta `json:"versions"`
}

func(s *Service) historyDir(deckID string) string{
	return filepath.Join(s.decksDir,deckID,"history")
}

//读索引
func(s *Service) readHistoryIndex(deckID string) historyIndex {
	data,err := os.ReadFile(filepath.Join(s.historyDir(deckID),"index.json"))
	if err !=nil{
		return historyIndex{NextSeq: 1}
	}
	var idx historyIndex
	if json.Unmarshal(data,&idx) !=nil || idx.NextSeq < 1{
		return historyIndex{NextSeq: 1}
	}

	return idx
}

//写入索引
func(s *Service)writeHistoryIndex(deckID string,idx historyIndex) error{
	if err := os.MkdirAll(s.historyDir(deckID),0o755); err !=nil{
		return err
	}

	data,err := json.Marshal(idx)
	if err != nil{
		return err
	}

	return atomicWriteFile(filepath.Join(s.historyDir(deckID),"index.json"),data)
}

// 一轮 agent run 结束后，对其中动过的每份 deck 各记一笔版本
// 由agent 层在 run成功结束或 ask_user 暂停时调用； 失败的run不调用
// detail 是本轮操作汇总。快照是三件套 bundle（index.html + style.css + outline.json），
// runtime/base.css 是共享资产不进快照（版本化路径逃生舱）。
func(s *Service) RecordRunVersion(userID uint,deckID,detail string) error{
	if err := s.authorize(userID,deckID);err !=nil{
		return err
	}
	unlock := s.lockDeck(deckID)
	defer unlock()

	return s.recordVersionV2(deckID, OpRun, detail)
}

// 分层保留 最近 keepFullDays 天全留；更早的每天只留当天最后
// 记录的一条 （Versions 新->旧，同一天先遇到的就是更晚记录的)。总数超限时截尾
// 返回保留列表和应删除文件的版本号
func pruneVersions(versions []VersionMeta)(kept []VersionMeta,dropped []string){
	cutoff := time.Now().AddDate(0,0,-keepFullDays)
	dayOf := func(t int64) string {return time.Unix(t,0).Format("2006-01-02")}

	seenOldDays := make(map[string]bool)
	for _,m := range versions{
		recent := time.Unix(m.Time,0).After(cutoff)
		if recent || !seenOldDays[dayOf(m.Time)] {
			kept = append(kept, m)
			if !recent {
				seenOldDays[dayOf(m.Time)] = true
			}
		}else{
			dropped = append(dropped, m.Version)
		}
	}
	if len(kept) > maxVersionsPerDeck {
		for _,m := range kept[maxVersionsPerDeck:] {
			dropped = append(dropped, m.Version)
		}
		kept = kept[:maxVersionsPerDeck]
	}

	return kept,dropped
}

// ListVersions 获取某份 deck 的版本历史
func(s *Service) ListVersions(userID uint,deckID string)([]VersionMeta,error){
	if err := s.authorize(userID,deckID); err !=nil{
		return nil,err
	}

	idx := s.readHistoryIndex(deckID)
	if idx.Versions == nil{
		idx.Versions = []VersionMeta{}
	}
	return idx.Versions,nil
}

// 用户手动删除一条历史
func (s *Service) DeleteVersion(userID uint, deckID, version string) error {
	if err := s.authorize(userID, deckID); err != nil {
		return err
	}
	return s.deleteVersionLocked(deckID, version)
}

// deleteVersionLocked 删除的核心逻辑（锁内）。拆出来是为了白盒测试能绕开 authorize，
// 和 restoreSnapshot 是同一个套路。
func (s *Service) deleteVersionLocked(deckID, version string) error {
	if !versionPattern.MatchString(version) {
		return fmt.Errorf("版本号 %q 不合法", version)
	}
	unlock := s.lockDeck(deckID)
	defer unlock()

	idx := s.readHistoryIndex(deckID)
	kept := make([]VersionMeta,0,len(idx.Versions))
	found := false
	for _,m := range idx.Versions{
		if m.Version == version {
			found = true
			continue
		}
		kept = append(kept, m)
	}
	if !found {
		return fmt.Errorf("版本 %s 不存在",version)
	}
	if err := os.Remove(filepath.Join(s.historyDir(deckID),version+".html"));err !=nil && !os.IsNotExist(err){
		return fmt.Errorf("删除快照文件失败：%w",err)
	}
	idx.Versions = kept
	return s.writeHistoryIndex(deckID,idx)
}

//清空某份deck的全部历史，NextSeq 不清零 编号不复用
func (s *Service) ClearHistory(userID uint, deckID string) (int, error) {
	if err := s.authorize(userID, deckID); err != nil {
		return 0, err
	}
	return s.clearHistoryLocked(deckID)
}

// clearHistoryLocked 清空的核心逻辑（锁内），同样为白盒测试拆出。
func (s *Service) clearHistoryLocked(deckID string) (int, error) {
	unlock := s.lockDeck(deckID)
	defer unlock()

	idx := s.readHistoryIndex(deckID)
	n := len(idx.Versions)
	dir := s.historyDir(deckID)
	for _,m := range idx.Versions {
		if err := os.Remove(filepath.Join(dir,m.Version+".html"));err !=nil&&!os.IsNotExist(err){
			return 0,fmt.Errorf("删除快照 %s 失败: %w",m.Version,err)
		}
	}
	idx.Versions = []VersionMeta{}
	if err := s.writeHistoryIndex(deckID,idx); err !=nil{
		return 0,err
	}
	return n,nil
}

