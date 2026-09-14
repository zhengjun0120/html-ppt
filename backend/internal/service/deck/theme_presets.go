package deck

import "strings"

// 风格预设：**有名字的整套美学**，每个名字背后是一整组具体的设计决策
// （配色 + 面板色 + 字体配对 + 圆角 + 纹理 + 语义色），不是"换个强调色"。
//
// 为什么需要它——这个文件存在的全部理由：
// 让模型"自由发挥"配色，它必然回到训练分布的中位数：深蓝底 + 青绿/紫的发光强调色 +
// 圆角卡片 + 无衬线字体。这是统计规律，不是态度问题；它见过的绝大多数"现代网页"
// 就长这样，"求稳"就等于选它。所以提高自由度治不了同质化——**给出具体的参照物才可以**。
// "纸感""双色印刷""终端"这些名字本身就带着一整套有据可依的视觉决定。
//
// 为什么预设是"写路径上的糖"而不是活引用：套用一次就把整套字段落进主题 JSON，
// 渲染路径完全不知道预设的存在（只认存储的字段）。反过来做（存 preset 名、渲染时现查表）
// 意味着改一次预设定义会静默改掉所有用过它的 deck——用户没碰过的页面自己变了样，
// 这正是这个代码库最不想制造的那类"不可见的变化"。
//
// 顺带的收益：预设里预置了 --accent-2 / --positive / --warn 三个语义色。
// 一份 deck 只有"一个强调色"时，需要区分正负、或者画一条双色对比的页面，模型就只能
// 就地写死颜色（那会让 update_theme 之后失效）。预置好这三个，它直接 var() 引用即可。
const (
	PresetPaper     = "paper"
	PresetEditorial = "editorial"
	PresetNoir      = "noir"
	PresetDuotone   = "duotone"
	PresetTerminal  = "terminal"
	PresetTech      = "tech"
)

// Preset 是一套预设的公开描述。Name 是工具参数里的取值，
// Label/About 面向模型与用户——About 要在**一句话里说清"什么时候该选它"**，
// 因为模型选预设靠的就是这一句。
type Preset struct {
	Name  string
	Label string
	About string
	Theme Theme // 只填观感字段；Transition/Canvas 不参与（那是演示参数，不是视觉语言）
}

// presets 是预设表。tech 刻意排在最后并在 About 里点明"最容易撞衫"：
// 深底发光色不是被删掉了，而是从"默认那一套"降级成"用户点名要科技感时才用的一套"。
var presets = map[string]Preset{
	PresetPaper: {
		Name: PresetPaper, Label: "纸感",
		About: "暖白纸面 + 墨色正文 + 印泥红，衬线标题，近乎直角，带细网格纹理。讲稿、人文、复盘，以及任何需要「像是人写出来的」的场合",
		Theme: Theme{
			Accent: "#b23a2e", Background: "#f6f2e9",
			HeadingColor: "#1c1a17", TextColor: "#2b2723",
			Surface: "#fffdf8", BorderColor: "#ded4c3",
			Font: FontEditorial, Radius: "2px", Texture: TextureGrid,
			Vars: map[string]string{
				"--accent-2": "#4a6fa5", // 褪色的靛青，做对比/图表第二色
				"--positive": "#2f6b4f",
				"--warn":     "#a8721f",
			},
		},
	},
	PresetEditorial: {
		Name: PresetEditorial, Label: "编辑风",
		About: "纯白纸面 + 黑主色 + 一处橙红点缀，零装饰零圆角，强对比大字。观点表达、评论、产品发布，需要「设计过」而不是「生成过」的场合",
		Theme: Theme{
			Accent: "#111111", Background: "#ffffff",
			HeadingColor: "#0d0d0d", TextColor: "#1f1f1f",
			Surface: "#f7f7f7", BorderColor: "#e2e2e2",
			Font: FontEditorial, Radius: "0px", Texture: TextureNone,
			Vars: map[string]string{
				"--accent-2": "#e8590c", // 唯一的彩色，留给真正要强调的那一处
				"--positive": "#1a7f4b",
				"--warn":     "#b8860b",
			},
		},
	},
	PresetNoir: {
		Name: PresetNoir, Label: "暗夜电影",
		About: "近黑 + 暖金强调 + 全衬线，大留白。需要氛围与仪式感的场合：年度总结、品牌故事、开场与结尾",
		Theme: Theme{
			Accent: "#d9a441", Background: "#0e0e11",
			HeadingColor: "#f4efe6", TextColor: "#cfcbc2",
			Surface: "#17171c", BorderColor: "#2e2e36",
			Font: FontSerif, Radius: "0px", Texture: TextureNone,
			Vars: map[string]string{
				"--accent-2": "#8a93a6",
				"--positive": "#7fa650",
				"--warn":     "#c9762c",
			},
		},
	},
	PresetDuotone: {
		Name: PresetDuotone, Label: "双色印刷",
		About: "米黄纸 + 深靛墨 + 荧光橙，深色硬描边 + 印刷网点。设计、创意、年轻化的表达，要「印出来」的质感时用",
		Theme: Theme{
			Accent: "#ff5a36", Background: "#f4efe3",
			HeadingColor: "#12233d", TextColor: "#2b3a4d",
			Surface: "#fbf7ec", BorderColor: "#12233d", // 描边用墨色而不是强调色，才有版画感
			Font: FontModern, Radius: "0px", Texture: TextureDots,
			Vars: map[string]string{
				"--accent-2": "#12233d",
				"--positive": "#2f6b4f",
				"--warn":     "#e0a020",
			},
		},
	},
	PresetTerminal: {
		Name: PresetTerminal, Label: "终端",
		About: "纯黑 + 磷光绿 + 等宽字体 + 扫描线纹理。技术分享、架构剖析、代码讲解，给工程师看的场合",
		Theme: Theme{
			Accent: "#39ff7a", Background: "#080c0a",
			HeadingColor: "#5cff9d", TextColor: "#a9d6b0",
			Surface: "#0f1512", BorderColor: "#1f3527",
			Font: FontMono, Radius: "0px", Texture: TextureRule,
			Vars: map[string]string{
				"--accent-2": "#62d0ff",
				"--positive": "#7dffb0",
				"--warn":     "#ffd166",
			},
		},
	},
	PresetTech: {
		Name: PresetTech, Label: "科技暗色",
		About: "深蓝底 + 青绿发光 + 圆角卡片。这一套最容易撞衫（它是「现代网页」的最大公约数），只在用户明确要科技感/现代感时才用",
		Theme: Theme{
			Accent: "#5eead4", Background: "#0b132b",
			HeadingColor: "#5eead4", TextColor: "#e2e8f0",
			// Surface/BorderColor 留空 = 由 accent 派生（原来的观感原样保留）
			Font: FontSans, Radius: "10px", Texture: TextureNone,
			Vars: map[string]string{
				"--accent-2": "#818cf8",
				"--positive": "#4ade80",
				"--warn":     "#fbbf24",
			},
		},
	},
}

// presetOrder 是稳定顺序：map 遍历随机，而它出现在工具说明与报错信息里，
// 顺序变了 diff 与文档就跟着抖。
var presetOrder = []string{PresetPaper, PresetEditorial, PresetNoir, PresetDuotone, PresetTerminal, PresetTech}

// Presets 按稳定顺序返回全部预设（工具说明、文档、测试用）。
func Presets() []Preset {
	out := make([]Preset, 0, len(presetOrder))
	for _, name := range presetOrder {
		out = append(out, presets[name])
	}
	return out
}

// PresetNames 返回稳定顺序的预设名清单（schema enum 契约测试用）。
func PresetNames() []string {
	return append([]string(nil), presetOrder...)
}

// PresetSummary 生成给模型看的预设清单（一行一个："名字（中文名）：一句话性格"）。
// 由 Go 从预设表直接生成而不是在 prompt 里再抄一遍：抄一遍就多一处会漂移的真相。
func PresetSummary() string {
	var sb strings.Builder
	for _, p := range Presets() {
		sb.WriteString("- ")
		sb.WriteString(p.Name)
		sb.WriteString("（")
		sb.WriteString(p.Label)
		sb.WriteString("）：")
		sb.WriteString(p.About)
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// applyTo 把整套观感铺到主题上。**不动 transition/canvas**：
// 它们决定的是"怎么翻页、画布多宽"，属于演示参数；用户说"用纸感，翻页改成淡入"
// 之后这份 deck 仍然是纸感，不该因此丢掉预设名。
func (p Preset) applyTo(t *Theme) {
	t.Preset = p.Name
	t.Accent = p.Theme.Accent
	t.Background = p.Theme.Background
	t.HeadingColor = p.Theme.HeadingColor
	t.TextColor = p.Theme.TextColor
	t.Surface = p.Theme.Surface
	t.BorderColor = p.Theme.BorderColor
	t.Font = p.Theme.Font
	t.Radius = p.Theme.Radius
	t.Texture = p.Theme.Texture
	t.Vars = copyVars(p.Theme.Vars)
}

// copyVars 拷一份再存：预设表是包级共享变量，裸赋值会让"改了一个 deck 的 vars"
// 穿透到所有用过这个预设的 deck（和 ThemePatch 里那条同一个理由）。
func copyVars(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
