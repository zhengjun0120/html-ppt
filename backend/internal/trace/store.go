package trace

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 白名单。**所有进 filepath.Join 的外部字符串都要先过这里**——run_id 和图片名
// 都来自 URL，不校验就能用 ../../ 绕出 data 目录读写任意文件。
// 与 deck 包的 idPattern 是同一条规矩（那边注释写得更细：白名单是唯一的安全闸门）。
//
// 图片名比 run_id 宽松一位：允许 v1-p003.png 这种"第几轮审查的第几页"命名，
// 因为一个 run 里可能审查多次（write_deck 自动审一次，之后 review_slides 再审），
// 只按页号命名会让后一次覆盖前一次。
var (
	idPattern    = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)
	imagePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,32}\.png$`)
)

func ValidID(s string) bool { return idPattern.MatchString(s) }

func ValidImageName(s string) bool { return imagePattern.MatchString(s) }

// ErrNotFound 一切"读不到"的统一出口：不存在、不属于你、文件名不合法，
// 全走这一个错误 → handler 一律 404。刻意不区分：区分了就等于告诉调用方
// "这个 id 存在但不是你的"（与 deck.EnsureOwner 不泄露存在性一致）。
var ErrNotFound = errors.New("trace: 记录不存在")

// errNoLastLine 表示文件尾部那一行比读取窗口还长，拿不到完整行。
// 不是 IO 故障，而是"这个 run 的 run_end 读不到" → 当成还在运行。
var errNoLastLine = errors.New("trace: 尾部没有完整的行")

// tailWindow 读最后一行时从文件尾部回看的字节数。
// run_end 本身很小（汇总 + 分项用量），回看 64KB 足够覆盖它以及它前面那条事件的一部分。
const tailWindow = 64 << 10

// Store 读取端。只认一个目录，不碰数据库——所以测试里给个 t.TempDir() 就能跑，
// 与 deck 那边"白盒绕开 authorize"的老办法相比，这边根本不需要绕。
type Store struct {
	dir string
	now func() time.Time
}

func NewStore(dir string) *Store {
	return &Store{dir: dir, now: time.Now}
}

// WithClock 只给测试用：把"当前时间"钉住，好断言"运行中的 run 耗时 = now - started"。
func (s *Store) WithClock(now func() time.Time) *Store {
	s.now = now
	return s
}

// ListRuns 列出某个用户可见的 run，新→旧。sessionID 为 0 表示不限会话。
//
// 只读每个文件的首行（run_start）与尾行（run_end），不解析全文：
// 一个开了全量上下文的 run 可能几十 MB，列表页要扫几十上百个，
// 全解析会让"打开观测页"变成一次几十秒的等待。
func (s *Store) ListRuns(userID uint, sessionID uint, limit, offset int) ([]RunMeta, error) {
	var dirs []string
	if sessionID > 0 {
		dirs = []string{filepath.Join(s.dir, strconv.FormatUint(uint64(sessionID), 10))}
	} else {
		entries, err := os.ReadDir(s.dir)
		if err != nil {
			if os.IsNotExist(err) {
				// 一个 trace 都还没产生：这是正常状态，不是错误
				return []RunMeta{}, nil
			}
			return nil, fmt.Errorf("trace: 读取 traces 目录失败 %w", err)
		}
		for _, e := range entries {
			if e.IsDir() {
				dirs = append(dirs, filepath.Join(s.dir, e.Name()))
			}
		}
	}

	all := make([]RunMeta, 0, 16)
	for _, d := range dirs {
		metas, err := s.listRunMetas(d)
		if err != nil {
			// 单个会话目录坏掉不该让整个列表打不开——其余会话的记录仍然有用
			log.Printf("[warn] trace: 列出 %s 里的 run 失败 err: %v", d, err)
			continue
		}
		for _, m := range metas {
			if m.UserID == userID {
				all = append(all, m)
			}
		}
	}

	sortRunMetasNewestFirst(all)
	all = applyWindow(all, offset, limit)
	if all == nil {
		all = []RunMeta{} // 空历史要序列化成 []，不是 null
	}
	return all, nil
}

// SumBySession 把 run 按会话累加。抽成纯函数（无文件系统、无 DB）是因为
// "一次对话花了多少 token"正是用户最关心的那个数，它错了不会报错、只会悄悄偏小。
func SumBySession(metas []RunMeta) []SessionMeta {
	byID := map[uint]*SessionMeta{}
	var order []uint
	for _, m := range metas { // metas 已按新→旧，所以每个会话第一次出现的就是它最新的 run
		sm, ok := byID[m.SessionID]
		if !ok {
			sm = &SessionMeta{SessionID: m.SessionID, Latest: m.RunID}
			byID[m.SessionID] = sm
			order = append(order, m.SessionID)
		}
		sm.Runs++
		if m.StartedAt.After(sm.LastAt) {
			sm.LastAt = m.StartedAt
		}
		if sm.DeckID == "" && m.DeckID != "" {
			sm.DeckID = m.DeckID
		}
		if m.Status == StatusRunning {
			sm.RunningRuns++
		}
		if m.Usage == nil {
			continue
		}
		sm.Usage.DurationMS += m.Usage.DurationMS
		sm.Usage.Turns += m.Usage.Turns
		sm.Usage.ToolCalls += m.Usage.ToolCalls
		if sm.Usage.Usage == nil {
			sm.Usage.Usage = make(map[string]UsagePart, len(m.Usage.Usage))
		}
		for k, v := range m.Usage.Usage {
			sm.Usage.Usage[k] = sm.Usage.Usage[k].Add(v)
		}
		sm.Usage.Total = sm.Usage.Total.Add(m.Usage.Total)
	}

	out := make([]SessionMeta, 0, len(order))
	for _, id := range order {
		out = append(out, *byID[id])
	}
	return out
}

// ReadOptions 读事件的坐标。**增量拉取用 FromOffset，不要用 FromSeq**：
// seq 是逻辑序号，按它过滤仍然要把整个文件从头扫一遍（一个 run 十几 MB 时，
// 实时跟随每 800ms 重扫一次是纯浪费）；offset 是字节位置，从它 seek 之后
// 只读新增的那一段，代价与新增量成正比。
type ReadOptions struct {
	FromSeq         int
	FromOffset      int64
	IncludeMessages bool // 默认剔除 messages：它是文件里唯一的大字段
	Limit           int  // 0 = 不限
}

// ReadResult 读事件的结果。NextOffset 交给调用方作为下一次增量拉取的游标；
// Running 表示文件里没有 run_end（还在跑，或进程被杀），实时页据此决定继续轮询。
type ReadResult struct {
	Events     []Event `json:"events"`
	NextOffset int64   `json:"next_offset"`
	Running    bool    `json:"running"`
}

// ReadEvents 读一个 run 的事件。
func (s *Store) ReadEvents(userID uint, sessionID uint, runID string, opt ReadOptions) (*ReadResult, error) {
	path, err := s.runFileFor(userID, sessionID, runID)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, ErrNotFound
	}
	defer f.Close()

	if opt.FromOffset > 0 {
		if _, err := f.Seek(opt.FromOffset, io.SeekStart); err != nil {
			return nil, fmt.Errorf("trace: 定位读取位置失败 %w", err)
		}
	}

	res := &ReadResult{Events: []Event{}}
	br := bufio.NewReaderSize(f, 64<<10)
	var read int64 = opt.FromOffset
	for {
		line, err := br.ReadBytes('\n')
		// ReadBytes 会把超长行完整读出来（没有 Scanner 那个 64KB 上限）——
		// 这一点很关键：带全量上下文的 llm_request 单行就是几 MB
		if len(line) > 0 {
			read += int64(len(line))
			trimmed := bytes.TrimRight(line, "\r\n")
			if len(trimmed) > 0 {
				var e Event
				if uerr := json.Unmarshal(trimmed, &e); uerr != nil {
					// 一行坏了不该让整次读取失败：实时读的正好可能是"写了一半"的那行边缘情况
					log.Printf("[warn] trace: 解析事件失败（跳过该行）err: %v", uerr)
				} else if e.Seq > opt.FromSeq {
					if !opt.IncludeMessages {
						e.Messages = nil
					}
					res.Events = append(res.Events, e)
					if opt.Limit > 0 && len(res.Events) >= opt.Limit {
						// 打到上限就停在这里。游标正好落在这一行之后，
						// 下一次从这里接着读，一条不漏也不重
						break
					}
				}
			}
		}
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("trace: 读取事件失败 %w", err)
			}
			break
		}
	}

	if opt.Limit == 0 || len(res.Events) < opt.Limit {
		// 游标 = 文件末尾的完整行边界。写到一半的那一行不推进游标，
		// 等下一次轮询时它已经写完，能被完整读到（这正是实时跟随要的行为）
		res.NextOffset = lastCompleteLineEnd(path)
	} else {
		res.NextOffset = read
	}

	// Running 一律从文件尾部判断，**不从刚读到的这一段判断**：
	// 增量轮询时 from_offset 已经越过了 run_end 之前的所有事件，
	// 靠"这一段里没见到 run_end"会永远显示成"运行中"。
	res.Running = !runHasEnded(path)
	return res, nil
}

// runHasEnded 看尾部那一行是不是 run_end。读不到（还在跑、或尾部行超长）都算没结束。
func runHasEnded(path string) bool {
	tail, err := readLastLine(path)
	if err != nil {
		return false
	}
	var end Event
	return json.Unmarshal(tail, &end) == nil && end.Kind == KindRunEnd
}

// ReadEvent 读单个事件（含 messages 全文）。观测页问"这一轮的完整上下文是什么"
// 时才调用——默认的列表读法剔掉了 messages，这是按需把它捞回来的唯一入口。
func (s *Store) ReadEvent(userID uint, sessionID uint, runID string, seq int) (*Event, error) {
	path, err := s.runFileFor(userID, sessionID, runID)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, ErrNotFound
	}
	defer f.Close()

	br := bufio.NewReaderSize(f, 64<<10)
	for {
		line, rerr := br.ReadBytes('\n')
		if len(line) > 0 {
			trimmed := bytes.TrimRight(line, "\r\n")
			if len(trimmed) > 0 {
				var e Event
				if json.Unmarshal(trimmed, &e) == nil && e.Seq == seq {
					return &e, nil
				}
			}
		}
		if rerr != nil {
			break
		}
	}
	return nil, ErrNotFound
}

// ImagePath 返回一张截图的绝对路径。三个参数都过白名单，归属也要复核——
// 图片路由是唯一直接用文件系统响应的地方，越权读的后果是拿到别人的 deck 截图。
func (s *Store) ImagePath(userID uint, sessionID uint, runID, name string) (string, error) {
	if !ValidImageName(name) {
		return "", ErrNotFound
	}
	if _, err := s.runFileFor(userID, sessionID, runID); err != nil {
		return "", err
	}
	return filepath.Join(s.dir, strconv.FormatUint(uint64(sessionID), 10), runID, "img", name), nil
}

// runFileFor 校验 + 归属检查 + 返回 jsonl 路径。
// 任何一步不过都返回 ErrNotFound，让上层的"不存在"和"不是你的"长得一模一样。
func (s *Store) runFileFor(userID uint, sessionID uint, runID string) (string, error) {
	if !ValidID(runID) {
		return "", ErrNotFound
	}
	sessDir := filepath.Join(s.dir, strconv.FormatUint(uint64(sessionID), 10))
	path := filepath.Join(sessDir, runID+".jsonl")

	line, err := readFirstLine(path)
	if err != nil {
		return "", ErrNotFound
	}
	var head Event
	if err := json.Unmarshal(line, &head); err != nil || head.Kind != KindRunStart {
		return "", ErrNotFound
	}
	if head.UserID != userID {
		return "", ErrNotFound
	}
	return path, nil
}

// listRunMetas 扫一个会话目录，只读每个 run 文件的首行与尾行。
func (s *Store) listRunMetas(sessDir string) ([]RunMeta, error) {
	return listRunMetasIn(sessDir, s.now)
}

// listRunMetasIn 是包级实现，为了写入侧（recorder 的裁剪）也能复用同一套元信息解析。
// 两边各写一份的话，裁剪用的顺序和列表显示的顺序迟早会不一致——
// 那意味着"列表上还看得见的 run 被裁掉了"。
func listRunMetasIn(sessDir string, now func() time.Time) ([]RunMeta, error) {
	entries, err := os.ReadDir(sessDir)
	if err != nil {
		return nil, err
	}
	metas := make([]RunMeta, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue // 图片目录 <run_id>/ 会被这里跳过
		}
		runID := strings.TrimSuffix(e.Name(), ".jsonl")
		if !ValidID(runID) {
			continue
		}
		path := filepath.Join(sessDir, e.Name())

		line, err := readFirstLine(path)
		if err != nil {
			continue
		}
		var head Event
		if err := json.Unmarshal(line, &head); err != nil || head.Kind != KindRunStart {
			// 首行不是 run_start：文件坏了或者版本对不上。跳过它，
			// 别让一个坏文件把整个列表拖成 500（观测工具自己的记录坏掉时，
			// 你更需要看到其余记录）
			continue
		}

		m := RunMeta{
			RunID:       head.RunID,
			ParentRunID: head.ParentRunID,
			SessionID:   head.SessionID,
			UserID:      head.UserID,
			DeckID:      head.DeckID,
			UserContent: head.UserContent,
			Model:       head.Model,
			StartedAt:   head.TS,
			Status:      StatusRunning,
		}
		if m.RunID == "" {
			m.RunID = runID
		}
		if st, serr := e.Info(); serr == nil {
			m.Bytes = st.Size()
		}

		if tail, terr := readLastLine(path); terr == nil {
			var end Event
			if json.Unmarshal(tail, &end) == nil && end.Kind == KindRunEnd {
				m.Status = end.Status
				if m.Status == "" {
					m.Status = StatusOK
				}
				if end.Summary != nil {
					sm := *end.Summary
					m.Usage = &sm
					m.DurationMS = sm.DurationMS
					m.Turns = sm.Turns
					m.ToolCalls = sm.ToolCalls
				}
			}
		}
		if m.Status == StatusRunning {
			// 没有 run_end：还在跑（或进程被杀了）。给个"到目前为止"的耗时，
			// 否则列表上这一行的时间永远是空的
			m.DurationMS = now().Sub(m.StartedAt).Milliseconds()
		}
		metas = append(metas, m)
	}
	return metas, nil
}

// ---------- 文件小工具 ----------

// readFirstLine 读首行。用 bufio.Reader 而不是 Scanner：Scanner 有 64KB 的行上限，
// 而带全量上下文的 llm_request 单行就是这个量级，会直接报 "token too long"。
func readFirstLine(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	line, err := bufio.NewReaderSize(f, 64<<10).ReadBytes('\n')
	if len(line) == 0 && err != nil {
		return nil, err
	}
	return bytes.TrimRight(line, "\r\n"), nil
}

// readLastLine 从文件尾部回看一个窗口，取最后一个完整的行。
//
// 直接从头扫到最后当然更简单，但那等于为了读一个小小的 run_end 把几十 MB 全读一遍。
// 反向读有一个必须处理的失败：最后一行比窗口还长（游标落在某条超长事件中间），
// 这时返回 errNoLastLine —— 调用方把它当成"没有 run_end"，也就是"还在运行"。
func readLastLine(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := st.Size()
	if size == 0 {
		return nil, errNoLastLine
	}
	start := int64(0)
	if size > tailWindow {
		start = size - tailWindow
	}
	buf := make([]byte, size-start)
	if _, err := f.ReadAt(buf, start); err != nil && err != io.EOF {
		return nil, err
	}
	buf = bytes.TrimRight(buf, "\r\n")
	if i := bytes.LastIndexByte(buf, '\n'); i >= 0 {
		return buf[i+1:], nil
	}
	if start > 0 {
		// 整个窗口里没有一个换行 = 最后一行比窗口还长，拿不到完整行
		return nil, errNoLastLine
	}
	return buf, nil
}

// lastCompleteLineEnd 返回"最后一个完整行的结尾"字节位置，用于增量游标。
// 末尾那种写到一半的行（没有换行）不计入，这样下一次轮询能完整读到它。
func lastCompleteLineEnd(path string) int64 {
	st, err := os.Stat(path)
	if err != nil {
		return 0
	}
	size := st.Size()
	if size == 0 {
		return 0
	}
	f, err := os.Open(path)
	if err != nil {
		return size
	}
	defer f.Close()

	start := size - 1
	if start > int64(tailWindow) {
		start = size - int64(tailWindow)
	}
	buf := make([]byte, size-start)
	if _, err := f.ReadAt(buf, start); err != nil && err != io.EOF {
		return size
	}
	if buf[len(buf)-1] != '\n' {
		// 尾部没有换行 = 有一行没写完。回退到它前面的那个换行
		buf = buf[:len(buf)-1]
	}
	if i := bytes.LastIndexByte(buf, '\n'); i >= 0 {
		return start + int64(i) + 1
	}
	return 0
}

// sortRunMetasNewestFirst 按新→旧。run_id 以 13 位定宽时间戳开头，字典序即时间序；
// 同一毫秒内并发产生的再按 StartedAt 兜底。
func sortRunMetasNewestFirst(metas []RunMeta) {
	sort.SliceStable(metas, func(i, j int) bool {
		if metas[i].RunID != metas[j].RunID {
			return metas[i].RunID > metas[j].RunID
		}
		return metas[i].StartedAt.After(metas[j].StartedAt)
	})
}

func applyWindow(metas []RunMeta, offset, limit int) []RunMeta {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(metas) {
		return []RunMeta{}
	}
	metas = metas[offset:]
	if limit > 0 && limit < len(metas) {
		metas = metas[:limit]
	}
	return metas
}

// atomicWriteFile 先写 .tmp 再 rename。借用 deck 那边的做法：
// 直接写目标文件的话，读的人有可能看到半张图（PNG 半截是打不开的白框）。
func atomicWriteFile(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
