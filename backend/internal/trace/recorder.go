package trace

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// Config 观测的落盘配置。开关单独放（Features.Trace），这里只管"存到哪、存多少"。
type Config struct {
	Enabled       bool
	Dir           string
	CaptureImages bool
	RetainRuns    int
	MaxFieldBytes int
}

// Options 建一个 Recorder 需要的东西。
type Options struct {
	Dir           string // traces 根目录（<data.dir>/traces）
	SessionID     uint
	RunID         string
	UserID        uint
	DeckID        string
	ParentRunID   string // ask_user 恢复时，指回被暂停的那个 run
	CaptureImages bool
	MaxFieldBytes int // 0 = 不截断
	RetainRuns    int // 每个会话最多留几个 run；0 = 不限
	Live          func(Event)
}

// Recorder 一次 run 的写入端。**必须由调用方 Close**，否则文件句柄会跟着会话数一起涨。
//
// 它只被 agent 循环那一个 goroutine 串行使用（工具内部的子调用也是同步跑在同一个
// goroutine 上），锁是为了防御将来的异步化，不是在争用。
type Recorder struct {
	// disabled 是 Discard 用的标记：构造后不变，所以读取它不需要加锁。
	// 单独一个字段而不是"nil 判断"，是为了让 From(ctx) 永远能返回一个可用的对象——
	// 让每个调用点都写 if r != nil 是不现实的，漏一个就是一次 panic。
	disabled bool

	runID      string
	filePath   string
	imgDir     string
	live       func(Event)
	capture    bool
	maxField   int
	retainRuns int
	started    time.Time

	mu         sync.Mutex
	file       *os.File
	seq        int
	turns      int
	toolCalls  int
	imageRound int
	usage      map[string]UsagePart
	closed     bool
}

// New 打开一个新 run 的 trace 文件。返回的错误交给调用方决定怎么处理——
// agent 那边的处理是"打条 warn 然后当没有观测继续跑"，见 recorder 包注释的 fail-open 原则。
func New(o Options) (*Recorder, error) {
	if o.Dir == "" {
		return nil, errors.New("trace: dir 不能为空")
	}
	if !ValidID(o.RunID) {
		// run_id 是自己生成的，走到这里说明生成逻辑坏了；但拼路径前必须校验
		// （与 deck 的 idPattern 同一条规矩：任何进 filepath.Join 的外部字符串都要过白名单）
		return nil, fmt.Errorf("trace: run_id %q 不合法", o.RunID)
	}

	sessDir := filepath.Join(o.Dir, strconv.FormatUint(uint64(o.SessionID), 10))
	if err := os.MkdirAll(sessDir, 0o755); err != nil {
		return nil, fmt.Errorf("trace: 创建目录失败 %w", err)
	}

	r := &Recorder{
		runID:      o.RunID,
		filePath:   filepath.Join(sessDir, o.RunID+".jsonl"),
		imgDir:     filepath.Join(sessDir, o.RunID, "img"),
		live:       o.Live,
		capture:    o.CaptureImages,
		maxField:   o.MaxFieldBytes,
		retainRuns: o.RetainRuns,
		started:    time.Now(),
	}

	// 追加模式而不是"攒到最后整文件写回"：实时查看靠的就是"写到哪看到哪"，
	// 整文件覆盖会让一个正在跑的 run 在观测页上从头到尾都是空的。
	// 也因此不能用 atomicWriteFile 那套 tmp+rename（它是给整份替换用的）。
	f, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("trace: 打开 trace 文件失败 %w", err)
	}
	r.file = f

	r.pruneOldRuns(sessDir)
	return r, nil
}

// NewRunID 生成 run_id：13 位 unix 毫秒 + 4 位随机十六进制。
//
// 定宽时间戳在前，字典序就等于时间序，所以列表和裁剪都不需要额外的时间字段排序；
// 后缀随机位是为了同一毫秒内并发两个 run 也不撞车。用 crypto/rand + hex 而不是
// 引入 uuid 依赖，是这个仓库既有的做法（见 vision/grant.go 的 nonce）。
func NewRunID(t time.Time) string {
	var b [2]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand 失败极罕见。观测工具不能因为这个阻断对话，退回无后缀
		return fmt.Sprintf("%013d-0000", t.UnixMilli())
	}
	return fmt.Sprintf("%013d-%s", t.UnixMilli(), hex.EncodeToString(b[:]))
}

// Close 收尾。可重复调用。
func (r *Recorder) Close() {
	if r == nil || r.disabled {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return
	}
	r.closed = true
	if r.file != nil {
		if err := r.file.Close(); err != nil {
			log.Printf("[warn] trace: 关闭 trace 文件失败 err: %v", err)
		}
		r.file = nil
	}
}

// writeLocked 是唯一的落盘入口：定序、补图片路径、截断、序列化、写一行。
// 返回最终事件（被丢弃时 Seq 为 0），供调用方在锁外推实时。
func (r *Recorder) writeLocked(e Event) Event {
	if r.closed || r.file == nil {
		return Event{}
	}
	r.seq++
	e.Seq = r.seq
	if e.TS.IsZero() {
		e.TS = time.Now()
	}

	// 计数在这里做，而不是让 agent 循环自己数：tool_call 是 execTool 发的、
	// llm_request 是 runLoop 发的，散在两地各数各的迟早会对不上
	switch e.Kind {
	case KindToolCall:
		r.toolCalls++
	case KindLLMRequest:
		r.turns++
	}

	r.attachImages(&e)
	r.truncateFields(&e)

	line, err := encodeLine(e)
	if err != nil {
		log.Printf("[warn] trace: 序列化事件失败 kind=%s err: %v", e.Kind, err)
		return Event{}
	}
	// 单次 Write 一次性写出整行（含换行），读者的 bufio 按行切就能拿到完整事件；
	// 分多次写会出现"读到半行"的窗口，实时轮询正好会撞上
	if _, err := r.file.Write(line); err != nil {
		log.Printf("[warn] trace: 写入 trace 失败 kind=%s err: %v", e.Kind, err)
		return Event{}
	}
	return e
}

// emitLockedThenLive 收敛"加锁写盘 → 解锁 → 推实时"这个固定动作。
//
// 实时回调**必须在锁外**调用：它要做通道发送，而 serveAgentSSE 的通道只有 16 槽，
// 满时会阻塞。持着锁阻塞 = 整个 agent 循环卡在观测工具上，这是观测工具最不该犯的错。
func (r *Recorder) emitLockedThenLive(e Event) {
	var live Event
	r.mu.Lock()
	live = r.writeLocked(e)
	r.mu.Unlock()

	if r.live != nil && live.Seq > 0 {
		r.live(live)
	}
}

// Emit 直接往这个 recorder 写一条（不读 ctx）。多半该用包级的 Emit(ctx, ...)。
func (r *Recorder) Emit(e Event) {
	if r == nil || r.disabled {
		return
	}
	r.emitLockedThenLive(e)
}

// NextImageRound 返回本 run 内递增的审查轮次，供截图命名使用。
//
// 为什么需要它：一次 run 里可能审查多次（write_deck 自动审一次，之后 review_slides 再审），
// 而它们都从第 1 页开始拍。只按页号命名的话，后一次会覆盖前一次的图，
// 于是第一次审查的 trace 卡片指向的图片内容已经变了——这种"记录指向的东西悄悄换了"
// 比图片丢失更难发现。名字形如 v1-p003.png。
func (r *Recorder) NextImageRound() int {
	if r == nil || r.disabled {
		return 1
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.imageRound++
	return r.imageRound
}

// RunID 本次 run 的标识。ask_user 暂停时要把这个值存进会话，
// 恢复时作为 parent_run_id 带回来——这样观测页能把被打断的一次对话串成一条。
func (r *Recorder) RunID() string {
	if r == nil || r.disabled {
		return ""
	}
	return r.runID
}

// Summary 当前累计的汇总，用于 run_end。
func (r *Recorder) Summary() Summary {
	if r == nil || r.disabled {
		return Summary{}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.summaryLocked()
}

func (r *Recorder) summaryLocked() Summary {
	// 复制一份：map 交出去之后调用方可能在锁外读，而后续 Usage 还会往里写
	us := make(map[string]UsagePart, len(r.usage))
	var total UsagePart
	for k, v := range r.usage {
		us[k] = v
		total = total.Add(v)
	}
	return Summary{
		DurationMS: time.Since(r.started).Milliseconds(),
		Turns:      r.turns,
		ToolCalls:  r.toolCalls,
		Usage:      us,
		Total:      total,
	}
}

// attachImages 把事件里的图片字节落盘，并把 Bytes 换成相对路径。
//
// 这一步是"视觉审查的截图终于看得见"的关键：capture.go 每页都实拍了一张
// 1244×700 的 PNG，过去发给审查模型之后就丢了（Slide.PNG 是 json:"-"），
// 所以报告里说"第 3 页文字被裁"时，人没有任何办法核对。
func (r *Recorder) attachImages(e *Event) {
	for i := range e.Images {
		img := &e.Images[i]
		if len(img.Bytes) == 0 {
			continue
		}
		if !r.capture {
			// 关了截图也要留下"本来有几张、多大"，否则观测页上完全看不出
			// 视觉审查拍过几页——那等于把"审查了什么"这件事又变回不可见
			img.Bytes = nil
			continue
		}
		if !ValidImageName(img.Name) {
			log.Printf("[warn] trace: 图片名 %q 不合法，丢弃", img.Name)
			img.Bytes = nil
			continue
		}
		if err := os.MkdirAll(r.imgDir, 0o755); err != nil {
			log.Printf("[warn] trace: 创建图片目录失败 err: %v", err)
			img.Bytes = nil
			continue
		}
		path := filepath.Join(r.imgDir, img.Name)
		if err := atomicWriteFile(path, img.Bytes); err != nil {
			log.Printf("[warn] trace: 写图片 %s 失败 err: %v", img.Name, err)
			img.Bytes = nil
			continue
		}
		// URL 只在写成功时非空。观测页据此区分"这张图没保存"和"图保存了但加载失败"——
		// 两者要采取的行动不一样（前者去看 capture_images 开关，后者去看文件权限）
		img.URL = filepath.ToSlash(filepath.Join(r.runID, "img", img.Name))
		img.Bytes = nil
	}
}

// truncateFields 按上限截断长字段。默认 maxField=0 = 不截断（全量存，便于回答
// "模型这一轮到底看到了什么"）。写满磁盘时把 max_field_bytes 打开即可瘦身。
func (r *Recorder) truncateFields(e *Event) {
	if r.maxField <= 0 {
		return
	}
	e.UserContent, _ = Truncate(e.UserContent, r.maxField)
	e.Args, _ = Truncate(e.Args, r.maxField)
	e.Result, _ = Truncate(e.Result, r.maxField)
	e.Content, _ = Truncate(e.Content, r.maxField)
	if e.Sub != nil {
		e.Sub.Text, _ = Truncate(e.Sub.Text, r.maxField)
	}
	if len(e.Messages) > r.maxField {
		// 不能直接在 JSON 数组里砍一刀：那会留下一个语法坏掉的 messages，
		// 读的人无法分辨"上下文就长这样"和"被截断了"。宁可明确说"这里没存"。
		e.Messages = json.RawMessage(fmt.Sprintf(`{"truncated":true,"bytes":%d}`, len(e.Messages)))
	}
}

// pruneOldRuns 保留最近 retainRuns 个 run，删掉更早的 jsonl 与图片目录。
// 每个 run 开始时跑一次：全量上下文存起来是几十 MB 级别的，不裁剪迟早撑爆 data 目录。
func (r *Recorder) pruneOldRuns(sessDir string) {
	if r.retainRuns <= 0 {
		return
	}
	metas, err := listRunMetasIn(sessDir, time.Now)
	if err != nil {
		log.Printf("[warn] trace: 列出历史 run 失败，跳过裁剪 err: %v", err)
		return
	}
	// 减 1：给"正在启动的这个 run"留一个名额。此时它的文件还是空的、
	// 不在 metas 里，不减的话 retain=N 实际会留下 N+1 个——
	// 而"配置写了 20 结果留了 21 个"这种事没人会去数，磁盘只会慢慢涨。
	drop := pruneRuns(metas, r.retainRuns-1)
	for _, m := range drop {
		if m.RunID == r.runID {
			continue // 不能把自己删了（理论上不会发生，防御一下）
		}
		_ = os.Remove(filepath.Join(sessDir, m.RunID+".jsonl"))
		_ = os.RemoveAll(filepath.Join(sessDir, m.RunID))
	}
}

// pruneRuns 纯函数：按新→旧排序后，保留前 keep 个，返回应当删除的。
// 抽成纯函数是为了能无文件系统测试——裁剪写错的表现是"悄悄删多了"，
// 这种错误不该靠"真跑一次再去看目录"来发现。
func pruneRuns(metas []RunMeta, keep int) []RunMeta {
	if keep <= 0 || len(metas) <= keep {
		return nil
	}
	sorted := make([]RunMeta, len(metas))
	copy(sorted, metas)
	// run_id 以定宽时间戳开头，字典序即时间序；同一毫秒的用 StartedAt 兜底
	sortRunMetasNewestFirst(sorted)
	return sorted[keep:]
}

// ---------- context 携带 ----------
//
// 一次 run 里只有"一个 recorder、一个当前轮次、一个当前工具"。
// 这些状态刻意走 ctx 而不是做成包级全局：agent 支持并发对话，
// 全局变量会让两个会话的事件串到一起，而且串得毫无规律。

var Discard = &Recorder{disabled: true}

type recorderKey struct{}
type scopeKey struct{}

// scope 当前事件该归到哪一轮、哪个工具调用下。
type scope struct {
	turn    int
	hasTurn bool // "第 0 轮"与"不适用"必须分得开，见 Event.Turn 的说明
	tool    string
	callID  string
}

// applyTo 把归属盖到事件上。Emit 与 Usage 共用，两条路径才不会一个带轮次一个不带。
func (s scope) applyTo(e *Event) {
	if s.hasTurn {
		t := s.turn
		e.Turn = &t
	}
	e.ToolName, e.ToolCallID = s.tool, s.callID
}

// With 把 recorder 挂进 ctx。r 为 nil / Discard 时原样返回，
// 于是"关掉 trace"不需要在业务代码里散落任何 if。
func With(ctx context.Context, r *Recorder) context.Context {
	if r == nil || r.disabled {
		return ctx
	}
	return context.WithValue(ctx, recorderKey{}, r)
}

// From 永不返回 nil：拿不到 recorder 时返回 Discard，调用方不必判空。
func From(ctx context.Context) *Recorder {
	if r, ok := ctx.Value(recorderKey{}).(*Recorder); ok && r != nil {
		return r
	}
	return Discard
}

// Active 报告这个 ctx 上是否真的在记录。
//
// 用于"要先付出构造代价才能发的事件"：比如把整份对话上下文序列化成 JSON
// 可能有好几 MB，关掉观测却每轮白序列化一遍就是纯亏。这类事件外面套一层
// if trace.Active(ctx) —— 其余事件不必套，Emit 自己在 Discard 上是空操作。
func Active(ctx context.Context) bool {
	r := From(ctx)
	return r != nil && !r.disabled
}

// WithTurn 标记"接下来这一轮"。runLoop 每轮开始调一次。
func WithTurn(ctx context.Context, turn int) context.Context {
	s := scopeFrom(ctx)
	s.turn, s.hasTurn = turn, true
	return context.WithValue(ctx, scopeKey{}, s)
}

// WithTool 标记"接下来这次工具调用"。从**当轮的 ctx** 派生，所以兄弟工具之间
// 不会互相污染（每个都从同一个 turnCtx 派生，各自覆盖自己的 tool/callID）。
func WithTool(ctx context.Context, name, callID string) context.Context {
	s := scopeFrom(ctx)
	s.tool, s.callID = name, callID
	return context.WithValue(ctx, scopeKey{}, s)
}

func scopeFrom(ctx context.Context) scope {
	if s, ok := ctx.Value(scopeKey{}).(scope); ok {
		return s
	}
	return scope{}
}

// Emit 记一条事件。归属（轮次 / 工具名 / tool_call_id）一律由 ctx 补齐。
//
// 刻意不信任调用方传进来的这三个字段：工具根本不知道自己在第几轮，
// 让它自己填的结果一定是空的或者错的。它只管说"发生了什么"。
func Emit(ctx context.Context, e Event) {
	r := From(ctx)
	if r == nil || r.disabled {
		return
	}
	scopeFrom(ctx).applyTo(&e)
	r.emitLockedThenLive(e)
}

// Usage 记一笔分项用量。component 取 CompMain / CompVision / CompWebSearch。
//
// 它同时干两件事：累加进 run 的汇总（run_end 用），以及**当场发一条 usage 事件**——
// 后者是为了实时看的时候能看着数字往上涨，而不是等这轮结束才知道花了多少。
func Usage(ctx context.Context, component string, part UsagePart) {
	r := From(ctx)
	if r == nil || r.disabled {
		return
	}
	if part.Calls == 0 {
		part.Calls = 1 // 没显式给次数就按"一次调用"记；否则分项上会出现 0 次却有 token
	}

	r.mu.Lock()
	if r.usage == nil {
		r.usage = make(map[string]UsagePart)
	}
	r.usage[component] = r.usage[component].Add(part)
	ev := Event{Kind: KindUsage, Component: component, Usage: &part}
	// 用量事件同样要带上轮次：页面上"第几轮花了多少"是按这些事件聚合的，
	// 漏掉归属就变成一堆没有出处的数字
	scopeFrom(ctx).applyTo(&ev)
	live := r.writeLocked(ev)
	r.mu.Unlock()

	if r.live != nil && live.Seq > 0 {
		r.live(live)
	}
}

// encodeLine 一行一个事件。不转义 HTML——事件里大量内容是 <section>、<div>
// 这类原文，转义成 \u003c 之后人读不了、diff 也看不出来，和 agent 那边的
// marshalNoEscape 是同一个理由。
func encodeLine(e Event) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(e); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil // Encode 自带结尾换行，正好当 JSONL 的分隔符
}
