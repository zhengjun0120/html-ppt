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
	"sync"
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

	// mu 保护 index。观测台列表是高频读路径，每次全盘重扫（几百次"读首行+
	// 尾行"，成本与磁盘上所有用户的 run 总数成正比）是它打开慢的根因——
	// 索引把成本挪到"新文件/文件变化才读一次"，之后纯内存过滤。
	mu    sync.Mutex
	index map[string]cachedRun // key: <会话目录名>/<runID>（两侧都过白名单/来自磁盘枚举）
}

// cachedRun 一条索引：列表口径的元信息 + 用于变化检测的文件指纹。
type cachedRun struct {
	meta  RunMeta
	size  int64
	mtime time.Time
}

func NewStore(dir string) *Store {
	return &Store{dir: dir, now: time.Now, index: map[string]cachedRun{}}
}

// WithClock 只给测试用：把"当前时间"钉住，好断言"运行中的 run 耗时 = now - started"。
func (s *Store) WithClock(now func() time.Time) *Store {
	s.now = now
	return s
}

// ListFilter 列表筛选条件（全空 = 不筛）。kind/status 由 handler 白名单校验。
type ListFilter struct {
	SessionID uint   // 0 = 不限会话
	Kind      string // "deck"（含旧数据的空 run_kind）/ "customize" / "tplsugg"
	Status    string // ok / paused / error / running
	Query     string // run_id / user_content / deck_id 的不区分大小写子串
}

// ListRuns 列出某个用户可见的 run，新→旧，带过滤后总数（分页页数靠它）。
//
// 元数据走内存索引：首次列表全盘扫一遍（只读每个文件的首行与尾行——一个开了
// 全量上下文的 run 可能几十 MB，全解析会让列表变成几十秒的等待），之后每次
// 只做目录枚举 + Stat 比对，文件没变直接用缓存。
func (s *Store) ListRuns(userID uint, f ListFilter, limit, offset int) ([]RunMeta, int, error) {
	all, err := s.listMetas(f.SessionID)
	if err != nil {
		return nil, 0, err
	}
	q := strings.ToLower(f.Query)
	filtered := make([]RunMeta, 0, len(all))
	for _, m := range all {
		if m.UserID != userID {
			continue
		}
		if f.Kind != "" && !kindMatches(m.RunKind, f.Kind) {
			continue
		}
		if f.Status != "" && m.Status != f.Status {
			continue
		}
		if q != "" && !queryMatch(m, q) {
			continue
		}
		filtered = append(filtered, m)
	}
	sortRunMetasNewestFirst(filtered)
	window := applyWindow(filtered, offset, limit)
	if window == nil {
		window = []RunMeta{} // 空历史要序列化成 []，不是 null
	}
	return window, len(filtered), nil
}

// kindMatches run_kind 筛选口径。"deck" 是兜底：旧数据没写 run_kind（空），
// 它们本来就是文稿对话；其余按精确匹配。
func kindMatches(runKind, want string) bool {
	if want == "deck" {
		return runKind == "" || runKind == "deck"
	}
	return runKind == want
}

// queryMatch q 子串匹配（q 已转小写）：run_id / 用户输入 / deck。
func queryMatch(m RunMeta, q string) bool {
	return strings.Contains(strings.ToLower(m.RunID), q) ||
		strings.Contains(strings.ToLower(m.UserContent), q) ||
		strings.Contains(strings.ToLower(m.DeckID), q)
}

// listMetas 索引读：先按磁盘现状校准索引，再返回（未过滤的）候选元数据。
// sessionID>0 只看那一个会话，0 = 全部。
func (s *Store) listMetas(sessionID uint) ([]RunMeta, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	dirs, err := s.sessionDirs(sessionID)
	if err != nil {
		return nil, err
	}
	s.refreshIndex(dirs, sessionID)
	prefix := ""
	if sessionID > 0 {
		prefix = strconv.FormatUint(uint64(sessionID), 10) + "/"
	}
	out := make([]RunMeta, 0, len(s.index))
	for key, cr := range s.index {
		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
		}
		out = append(out, cr.meta)
	}
	return out, nil
}

// sessionDirs 待校准的会话目录名（相对 s.dir）。目录不存在是正常状态：
// 一个 trace 都还没产生，不是错误。
func (s *Store) sessionDirs(sessionID uint) ([]string, error) {
	if sessionID > 0 {
		return []string{strconv.FormatUint(uint64(sessionID), 10)}, nil
	}
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("trace: 读取 traces 目录失败 %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// refreshIndex 把索引校准到磁盘现状：新文件、size+mtime 变化的文件重读首尾行，
// 消失的被裁剪/删除的）从索引剔除，其余直接用缓存。仍在运行中的 run（没有
// run_end）每次都重读——它的"到目前为止耗时"依赖当前时间，缓存会让列表上
// 的耗时停摆。
func (s *Store) refreshIndex(dirs []string, sessionID uint) {
	seen := make(map[string]struct{}, len(s.index))
	for _, name := range dirs {
		sessDir := filepath.Join(s.dir, name)
		entries, err := os.ReadDir(sessDir)
		if err != nil {
			continue // 目录被删/暂不可读：索引清尾负责剔掉它的旧条目
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
				continue // 图片目录 <run_id>/ 会被这里跳过
			}
			runID := strings.TrimSuffix(e.Name(), ".jsonl")
			if !ValidID(runID) {
				continue
			}
			key := name + "/" + runID
			info, err := e.Info()
			if err != nil {
				continue
			}
			if cr, ok := s.index[key]; ok &&
				cr.meta.Status != StatusRunning &&
				cr.size == info.Size() && cr.mtime.Equal(info.ModTime()) {
				seen[key] = struct{}{}
				continue
			}
			m := readRunMeta(sessDir, runID, info.Size(), s.now)
			if m == nil {
				continue // 坏文件不进索引；旧条目交给下面的清尾逻辑
			}
			s.index[key] = cachedRun{meta: *m, size: info.Size(), mtime: info.ModTime()}
			seen[key] = struct{}{}
		}
	}
	for key := range s.index {
		if _, ok := seen[key]; ok {
			continue
		}
		if sessionID == 0 || strings.HasPrefix(key, strconv.FormatUint(uint64(sessionID), 10)+"/") {
			delete(s.index, key)
		}
	}
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

// readRunMeta 读单个 run 文件的列表口径元信息（首行 run_start + 尾行 run_end）。
// 文件缺失或首行不是 run_start 返回 nil——调用方按"坏文件跳过"处理：
// 观测工具自己的记录坏掉时，你更需要看到其余记录，而不是整页 500。
// 索引刷新与 listRunMetasIn（写入侧裁剪共用）都走这里，口径不会漂移。
func readRunMeta(sessDir, runID string, size int64, nowFn func() time.Time) *RunMeta {
	path := filepath.Join(sessDir, runID+".jsonl")
	line, err := readFirstLine(path)
	if err != nil {
		return nil
	}
	var head Event
	if err := json.Unmarshal(line, &head); err != nil || head.Kind != KindRunStart {
		return nil
	}
	var end *Event
	if tail, terr := readLastLine(path); terr == nil {
		var te Event
		if json.Unmarshal(tail, &te) == nil && te.Kind == KindRunEnd {
			end = &te
		}
	}
	m := runMetaFrom(head, end, size, runID, nowFn)
	return &m
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
		var size int64
		if st, serr := e.Info(); serr == nil {
			size = st.Size()
		}
		if m := readRunMeta(sessDir, runID, size, now); m != nil {
			metas = append(metas, *m)
		}
	}
	return metas, nil
}

// runMetaFrom 把 run_start（head，必需）与 run_end（end，nil = 文件里没有，
// 还在跑或进程被杀）拼成列表口径的 RunMeta。列表扫描与导出共用这一个实现，
// 两边的"耗时/轮数/token"口径才不会漂移——漂移的表现是导出文件里的汇总
// 和列表页显示的对不上，查起来极费劲。
func runMetaFrom(head Event, end *Event, fileBytes int64, fallbackID string, nowFn func() time.Time) RunMeta {
	m := RunMeta{
		RunID:       head.RunID,
		ParentRunID: head.ParentRunID,
		SessionID:   head.SessionID,
		UserID:      head.UserID,
		DeckID:      head.DeckID,
		RunKind:     head.RunKind,
		UserContent: head.UserContent,
		Model:       head.Model,
		StartedAt:   head.TS,
		Status:      StatusRunning,
	}
	if m.RunID == "" {
		m.RunID = fallbackID
	}
	m.Bytes = fileBytes

	if end == nil {
		// 没有 run_end：还在跑（或进程被杀了）。给个"到目前为止"的耗时，
		// 否则列表上这一行的时间永远是空的
		m.DurationMS = nowFn().Sub(m.StartedAt).Milliseconds()
		return m
	}
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
	return m
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
