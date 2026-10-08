package handler

import "strings"

// 编辑器注入（docs/deck-editor-plan.md §4.3）。
//
// 仅在 ?edit=1 的读取路径上调用（deck 文件端点 / 用户模板 editor 端点）——
// 正常预览、导出、量测、无头审查拿到的一律是纯净 HTML。CSP 的
// script-src 'self' 放行同源 /assets/deck-v2/editor.js；connect-src 'none'
// 保证编辑器内不能发请求，保存数据只能 postMessage 出 iframe。
const editorScriptTag = `<script src="/assets/deck-v2/editor.js" defer></script>`

// injectEditorScript 在 </body> 前注入编辑器脚本。找不到 </body> 时整个追加
// 在尾部（fail-open：脚本 defer 执行，位置不严谨也能活）。
func injectEditorScript(html string) string {
	if i := strings.LastIndex(html, "</body>"); i >= 0 {
		return html[:i] + editorScriptTag + html[i:]
	}
	return html + editorScriptTag
}
