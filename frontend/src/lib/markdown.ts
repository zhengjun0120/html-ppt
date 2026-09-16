import MarkdownIt from 'markdown-it'

// html:false —— 模型产出的原始 HTML 一律转义，杜绝注入（不引 DOMPurify 的前提）
const md = new MarkdownIt({ html: false, linkify: true, breaks: true })

/** agent 文本（markdown）→ 安全 HTML */
export function renderMarkdown(src: string): string {
  return md.render(src)
}
