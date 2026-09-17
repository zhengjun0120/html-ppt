package agent

// v2 工具参数的健壮解码。
//
// 实测（trace 1789662713586-33b8，12 轮 17 万 token 的教训）：DeepSeek 输出的
// 工具参数有两种发病率——
//  1. JSON 里混入隐形字符（BOM/零宽空格），unmarshal 报 invalid character 'ï'；
//  2. 一次失败后开始把参数**再包一层** {"arguments": ...}（模仿传输层形状），
//     并在错误信息的误导下疯狂加字段，直到预算耗尽。
//
// 对策都在解码这一层：隐形字符清洗 + arguments 拆包重试，把这两种病 absorbs 掉。
// 报错文案也要给模型可执行的指令，而不是让它猜。

import (
	"encoding/json"
	"fmt"
	"strings"
)

// invisibleJSONChars 模型偶尔夹带、且 JSON 解析必炸的隐形字符。
var invisibleJSONChars = strings.NewReplacer(
	"\uFEFF", "", // BOM
	"\u200B", "", // zero-width space
	"\u200C", "", // ZWNJ
	"\u200D", "", // ZWJ
	"\u2060", "", // word joiner
	"\u00AD", "", // soft hyphen
)

// sanitizeModelJSON 清掉隐形字符（它们在任何位置都非法且不可见，删了只赚不赔）。
func sanitizeModelJSON(s string) string {
	return invisibleJSONChars.Replace(strings.TrimSpace(s))
}

// decodeToolArgs 把模型提交的 arguments 解码到 out。
//
// 顺序：
//  1. 单键 {"arguments": ...} 包裹探测（退化形态的解析会"成功"，只是字段全空——
//     所以必须在正常解码之前先探测拆包，而不是等解码失败）；
//  2. 清洗隐形字符后直接解码；
//  3. 全部失败返回带修复指令的错误——错误文案是给模型看的下一轮提示词。
func decodeToolArgs(arguments string, out any) error {
	raw := sanitizeModelJSON(arguments)

	// 包裹探测：外层只有一个 "arguments" 键就是退化形态（正常的参数里
	// 叫 arguments 的字段不存在于任何工具 schema）。
	var probe map[string]json.RawMessage
	if json.Unmarshal([]byte(raw), &probe) == nil {
		if len(probe) == 1 {
			if inner, ok := probe["arguments"]; ok {
				var innerStr string
				data := inner
				if json.Unmarshal(inner, &innerStr) == nil {
					data = []byte(innerStr) // 字符串化 JSON：再解一层
				}
				if err := json.Unmarshal([]byte(sanitizeModelJSON(string(data))), out); err == nil {
					return nil
				}
			}
		}
	}

	if err := json.Unmarshal([]byte(raw), out); err != nil {
		return fmt.Errorf("参数不是合法 json：%v。请把工具参数直接作为顶层 JSON 对象提交（不要包在 {\"arguments\": ...} 里，不要夹带注释或不可见字符）", err)
	}
	return nil
}
