package handler

import (
	"errors"
	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
	"html-ppt/backend/internal/trace"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 观测记录的读取接口（观测页 /trace 用）。
//
// 与图片路由的分工刻意不同：这里全部走 response 包（JSON API），
// 只有取图那一条用原生写法（二进制流不进 JSON）。
//
// 权限：挂在 guarded 组里，靠 authctx 拿用户、再靠 trace.Store 复核 run 归属。
// 归属不匹配一律 404 且**和"不存在"用同一个错误**——区分它们就等于告诉调用方
// "这个 run 是存在的，只是不是你的"，而 run_id 里带着会话号，
// 那等于泄露"这个用户聊过几次"（与 deck.EnsureOwner 的约定一致）。
//
// 读接口**不受 features.trace 开关影响**：关掉的是"继续记录"，
// 已经落盘的记录应该照样能看——否则调完开关，之前跑出来的东西就再也读不到了。

// maxRunsPerQuery 单次列表的硬上限。
// 会话汇总必须建立在**完整**的 run 列表上（否则"这次对话花了多少 token"会因为
// 分页而少算），所以这里先取一个够大的窗口，再在内存里分页。
// 500 个 run 的元信息读取是几百次"只读首尾行"，可以接受。
const maxRunsPerQuery = 500

// ListTraces GET /api/traces
//
// 查询参数：session_id（可选，只看某个会话）、limit / offset（可选，给 runs 分页）
// 返回：{ sessions: [...], runs: [...] }
//   - sessions 是按会话累加的汇总（跨 run），回答"这次对话花了多少 token"
//   - runs 是扁平的 run 列表，新→旧，页面左栏用它
func (h *Handler) ListTraces(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}

	sessionID, err := queryUint(c, "session_id")
	if err != nil {
		response.Err(c, http.StatusBadRequest, "session_id 必须是数字")
		return
	}
	limit, err := queryInt(c, "limit")
	if err != nil || limit < 0 {
		response.Err(c, http.StatusBadRequest, "limit 必须是非负整数")
		return
	}
	offset, err := queryInt(c, "offset")
	if err != nil || offset < 0 {
		response.Err(c, http.StatusBadRequest, "offset 必须是非负整数")
		return
	}

	metas, err := h.traces.ListRuns(uid, sessionID, maxRunsPerQuery, 0)
	if err != nil {
		// 具体原因（哪个目录读不了）进日志，响应里只给人看的话
		response.Err(c, http.StatusInternalServerError, "读取观测记录失败")
		return
	}

	// 会话汇总建立在完整列表上；runs 才按 offset/limit 截
	sessions := trace.SumBySession(metas)
	if offset > 0 {
		if offset >= len(metas) {
			metas = []trace.RunMeta{}
		} else {
			metas = metas[offset:]
		}
	}
	if limit > 0 && limit < len(metas) {
		metas = metas[:limit]
	}

	response.OK(c, gin.H{"sessions": sessions, "runs": metas})
}

// GetTraceRun GET /api/traces/:sessionID/:runID
//
// 查询参数：
//
//	from_offset  字节游标，**实时跟随用这个**（只读新增的那一段，代价与新增量成正比）
//	from_seq     只要 seq 大于它的（一次性查询用，会从头扫）
//	include      传 messages 时带上全量上下文（默认剔除，它是唯一的大字段）
//	limit        最多返回几条
//
// 返回 { events, next_offset, running }：next_offset 交给下一次增量拉取；
// running 由文件尾部有没有 run_end 判断（缺了它，页面分不清"还在跑"和"已经崩了"）。
func (h *Handler) GetTraceRun(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	sessionID, err := paramUint(c, "sessionID")
	if err != nil {
		response.Err(c, http.StatusNotFound, "观测记录不存在")
		return
	}

	fromOffset, err := queryInt64(c, "from_offset")
	if err != nil || fromOffset < 0 {
		response.Err(c, http.StatusBadRequest, "from_offset 必须是非负整数")
		return
	}
	fromSeq, err := queryInt(c, "from_seq")
	if err != nil || fromSeq < 0 {
		response.Err(c, http.StatusBadRequest, "from_seq 必须是非负整数")
		return
	}
	limit, err := queryInt(c, "limit")
	if err != nil || limit < 0 {
		response.Err(c, http.StatusBadRequest, "limit 必须是非负整数")
		return
	}

	res, err := h.traces.ReadEvents(uid, sessionID, c.Param("runID"), trace.ReadOptions{
		FromSeq:    fromSeq,
		FromOffset: fromOffset,
		// 显式传 include=messages 才带上上下文：默认带的话，一次拉取就是几十 MB，
		// 而绝大多数查看动作（看工具调了什么、搜了什么）并不需要看上下文原文
		IncludeMessages: strings.EqualFold(strings.TrimSpace(c.Query("include")), "messages"),
		Limit:           limit,
	})
	if err != nil {
		if errors.Is(err, trace.ErrNotFound) {
			response.Err(c, http.StatusNotFound, "观测记录不存在")
			return
		}
		response.Err(c, http.StatusInternalServerError, "读取观测记录失败")
		return
	}
	response.OK(c, res)
}

// GetTraceEvent GET /api/traces/:sessionID/:runID/events/:seq
//
// 单个事件的**全量**版本（含 messages）。观测页上"展开完整上下文"按钮走它——
// 列表接口默认把上下文剔掉了，这是按需把它捞回来的唯一入口。
func (h *Handler) GetTraceEvent(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	sessionID, err := paramUint(c, "sessionID")
	if err != nil {
		response.Err(c, http.StatusNotFound, "观测记录不存在")
		return
	}
	seq, err := strconv.Atoi(c.Param("seq"))
	if err != nil || seq <= 0 {
		response.Err(c, http.StatusNotFound, "观测记录不存在")
		return
	}

	ev, err := h.traces.ReadEvent(uid, sessionID, c.Param("runID"), seq)
	if err != nil {
		if errors.Is(err, trace.ErrNotFound) {
			response.Err(c, http.StatusNotFound, "观测记录不存在")
			return
		}
		response.Err(c, http.StatusInternalServerError, "读取观测记录失败")
		return
	}
	response.OK(c, ev)
}

// GetTraceImage GET /api/traces/:sessionID/:runID/img/:name
//
// 视觉审查的截图。这条走原生写法（二进制不进 JSON），失败也返回**裸状态码**：
// 它主要被 <img src> 加载，塞一个 {"error":...} 的 JSON 进去除了让控制台多一行
// 解析失败之外没有任何用处（与 render.go 的处理一致）。
//
// 令牌走 ?token= —— <img> 带不了 Authorization 头，而鉴权中间件原生支持这个回退
// （见 middleware/auth.go）。已知妥协：token 会进访问日志，与 deck 预览同一条。
func (h *Handler) GetTraceImage(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	sessionID, err := paramUint(c, "sessionID")
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	path, err := h.traces.ImagePath(uid, sessionID, c.Param("runID"), c.Param("name"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	png, err := os.ReadFile(path)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	// 图片名里带审查轮次与页号，写到磁盘之后内容不再变，所以可以放心缓一会儿。
	// private：这是用户自己的 deck 截图，不允许中间层缓存
	c.Header("Cache-Control", "private, max-age=300")
	c.Data(http.StatusOK, "image/png", png)
}

// ---------- 参数解析小工具 ----------
//
// 全部返回错误而不是"错了就当 0"：观测页的查询参数拼错时，
// 静默当成 0 会表现为"数据不对"（比如 offset 无效导致每次都从第一条开始，
// 页面看起来只是在原地刷新），而明确的 400 一眼就知道是参数问题。

func queryUint(c *gin.Context, key string) (uint, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0, nil
	}
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func paramUint(c *gin.Context, key string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func queryInt(c *gin.Context, key string) (int, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0, nil
	}
	return strconv.Atoi(raw)
}

func queryInt64(c *gin.Context, key string) (int64, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return 0, nil
	}
	return strconv.ParseInt(raw, 10, 64)
}
