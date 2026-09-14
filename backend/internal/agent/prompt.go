package agent

import (
	_ "embed"
	"strings"
)

//go:embed systemPrompt.md
var systemPromptRaw string

// systemPrompt 是真正发给模型的那份提示词：行尾统一成 LF。
//
// 为什么必须做这一步：//go:embed 原样嵌入文件字节，而这个文件在 Windows 上被编辑器
// 存成 CRLF，于是模型收到的是带 \r 的版本。实测（同一份内容、只差行尾，各发一次
// max_tokens=1 的请求读 usage.prompt_tokens）：LF 7003 token、CRLF 7794 token——
// 433 个 \r 让每行多花约 1.8 个 token，白交 11% 的输入费用，而且每次请求都要交一遍。
//
// 为什么在这里归一化而不是去改文件的保存格式：这个仓库整体是 CRLF，编辑器和 git 的
// autocrlf 都会把文件改回 CRLF——改一次文件只对当次有效。在入口处归一化才不受那些
// 设置影响，而且这行代码本身就是"为什么"的证据。
var systemPrompt = strings.ReplaceAll(systemPromptRaw, "\r\n", "\n")
