package agent

// 看图审查的公共常量与配额分流判定（v1/v2 两代审查共用；实现分别在
// vision_review_v2.go —— v1 的实现已随旧栈摘除）。

import (
	"errors"
	"encoding/json"
)

const reviewChars = 1500 //限制返回给模型的报告字数长度

// 单次最多看几页。图片 token 是整个 agent 最贵的一项，
// 不设上限的话"让 agent 自己挑页"很容易退化成"每页都挑"。
const maxReviewPages = 6

// reviewQuotaPerRun 一次 run 里 review_slides 的硬配额（Tool.MaxPerRun 的值）。
// 预期用法三步：全量审可疑页（1~2 次覆盖完）→ 批量修复 → 确认（1 次，只看改过的页）。
// 它是"审查→修复→再审"死循环的唯一硬闸门：提示词拦不住它，计数器拦得住。
const reviewQuotaPerRun = 3

// 开关没开。与"跑了但失败"必须分开：后者要出声，前者本来就该安静。
var errVisionOff = errors.New("视觉审查没开")

// visionReportMarker 是 review_slides 结果里"真的看到了图"的标记：
// 看图成功的结果一定带它，"只回数字"（pages 留空）和"看图失败"的结果一定不带。
// execTool 用它决定要不要把这次调用从配额里退还——配额只数烧了钱、产出了报告的审查。
const visionReportMarker = "\n看图审查：\n"

// visionFallbackNote 是看图失败时附在量测数字后面的说明。除了告诉模型"这次没有
// 看图的结论"，还要告诉它"这次不占配额"。
const visionFallbackNote = "\n(看图部分失败，以上只有量测数字；这次没有消耗审查配额，可以直接重试)"

// reviewMeasureOnlyArgs 判断 review_slides 的调用参数是不是"数字复查"（pages 缺省或空）。
// 解析失败返回 false：宁把它当真审查计一次配额，也不能让一个坏参数绕过闸门。
// 这个判定只服务配额分流，工具真正的参数解析仍归它自己的实现。
func reviewMeasureOnlyArgs(argsJSON string) bool {
	var a struct {
		Pages []int `json:"pages"`
	}
	if json.Unmarshal([]byte(argsJSON), &a) != nil {
		return false
	}
	return len(a.Pages) == 0
}
