package deck

// 页片段处理的共享 helper（v2 写路径使用；源自 v1 的 slide_ops.go / update.go /
// slides.go，v1 的服务方法随旧栈摘除，这几个纯函数保留）。

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"html-ppt/backend/internal/store"
)

// parseSlideFragment 校验 LLM 提交的单页片段（insert 与 update 共用的闸门）：
// 唯一根 <section>、禁嵌套、禁 doctype/完整文档，再过一遍消毒闸门
// （结构性危险标签拒绝；on* 事件与危险协议 URL 剥除并回报，见 sanitizeSlide）。
func parseSlideFragment(newHTML string) (*goquery.Selection, string, error) {
	raw := strings.TrimSpace(newHTML)
	if raw == "" {
		return nil, "", errors.New("new_html 不能为空")
	}

	lower := strings.ToLower(raw)
	if strings.Contains(lower, "<!doctype") || strings.Contains(lower, "<html") {
		return nil, "", errors.New("new_html 只能是一个 <section> 元素，不要输出完整的 HTML 文档")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err != nil {
		return nil, "", fmt.Errorf("解析 new_html 失败 err:%w", err)
	}

	roots := doc.Find("body").Children().Filter("section")
	if roots.Length() == 0 {
		return nil, "", errors.New("new_html 必须以 <section> 为根元素")
	}
	if roots.Length() > 1 || doc.Find("body").Children().Length() != roots.Length() {
		return nil, "", errors.New("new_html 只能包含一个 <section> 元素，多页修改请分多次调用")
	}
	sec := roots.First()
	if sec.Find("section").Length() > 0 {
		return nil, "", errors.New("不支持嵌套 <section>，页内不要嵌套 section")
	}

	rep := sanitizeSlide(sec)
	if err := rep.RejectErr(); err != nil {
		return nil, "", err
	}

	lint := newStyleLinter()
	lint.add(sec)

	return sec, joinWarnings(rep.Warning(), lint.report().Warning()), nil
}

var slideIDPattern = regexp.MustCompile(`^s(\d+)$`)

// slideFingerprint 页内容指纹（read_slide / update_slide 的乐观锁）。
func slideFingerprint(outer string) string {
	sum := sha256.Sum256([]byte(outer))
	return hex.EncodeToString(sum[:])[:12]
}

// slideTitle 页标题（h1/h2/h3 优先，兜底截前 30 字）。
func slideTitle(sec *goquery.Selection) string {
	for _, tag := range []string{"h1", "h2", "h3"} {
		if t := strings.TrimSpace(sec.Find(tag).First().Text()); t != "" {
			return t
		}
	}
	text := strings.Join(strings.Fields(sec.Text()), " ")
	if r := []rune(text); len(r) > 30 {
		text = string(r[:30]) + "..."
	}
	return text
}

// sanitizeAttrCopy 复制属性切片（消毒遍历时边改边遍历不安全，这里抽出防误用）。
func sanitizeAttrCopy(attrs []html.Attribute) []html.Attribute {
	return append([]html.Attribute(nil), attrs...)
}

// claimPage 归属校验的 DB 依赖封装：无库时报统一错误（v2 页路径都经过 authorize，
// 这里只是给直接调用 helper 的测试路径一个明确失败面）。
func claimPage(store_ *store.Store) error {
	if store_ == nil {
		return errStorage
	}
	_ = strconv.Itoa
	return nil
}
