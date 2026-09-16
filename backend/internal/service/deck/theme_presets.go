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
	// 后四套是按"题材"补的（来历见 presets 表里那段注释），不是又加了四个配色
	PresetIndigo = "indigo"
	PresetPine   = "pine"
	PresetKraft  = "kraft"
	PresetDune   = "dune"
	// 第 11 套是按"视觉语言的完备性"补的：前 10 套里没有一套是"中性底 + 单一高饱和色 +
	// 极细大字"的瑞士国际主义——那是演示排版史上有据可依的一档，也是整页色块节奏最天然的风格
	PresetSwiss = "swiss"
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
	// 下面四套补的是"题材缺口"，不是"再多几个配色"：上面六套里浅色的三套（纸感/编辑风/双色印刷）
	// 全偏暖，于是科研数据、自然可持续、人文历史、艺术设计这四个常见题材没有对应的观感可选，
	// 模型只能把内容往"暖纸 + 印泥红"或"深蓝发光"上塞——同一个题材换个题目会得到几乎一样的页面。
	//
	// 调性分组参考了外部 html-ppt skill 按题材给主题的思路（那个项目是 AGPL-3.0，所以只取
	// "题材 → 调性"这个想法，**取值全部按我们自己的约束重定**，没有搬它的色板）：它一套主题
	// 只给 ink/paper 两色，而我们的预设要铺满 accent/surface/border/字体/圆角/纹理/语义色一整套，
	// 且每套都要过正文与强调色的对比度下限（见 theme_presets_test.go 的对比度测试）。
	PresetIndigo: {
		Name: PresetIndigo, Label: "靛蓝",
		About: "瓷白底 + 深靛蓝 + 一点朱砂，衬线标题，细网格纹理。研究报告、数据汇报、技术方案，需要「像学术期刊」的场合",
		Theme: Theme{
			Accent: "#1b4a86", Background: "#f1f3f5",
			HeadingColor: "#0a1f3d", TextColor: "#26384f",
			// 面板是纯瓷白而不是强调色的淡染：浅色主题上淡染会让每张卡片泛一层蓝
			Surface: "#ffffff", BorderColor: "#d5dde7",
			Font: FontEditorial, Radius: "2px", Texture: TextureGrid,
			Vars: map[string]string{
				"--accent-2": "#c2453a", // 朱砂：只用在对比页/图表的第二色上，像瓷器上的那一笔记号
				"--positive": "#2f6b4f",
				"--warn":     "#a8721f",
			},
		},
	},
	PresetPine: {
		Name: PresetPine, Label: "松墨",
		About: "象牙底 + 深松绿 + 赭石，全衬线，零纹理零圆角。自然、可持续、文化纪实、非虚构，需要「旧杂志内页」的场合",
		Theme: Theme{
			Accent: "#2f5d3a", Background: "#f5f1e8",
			HeadingColor: "#1a2e1f", TextColor: "#34422e",
			Surface: "#fdfbf5", BorderColor: "#dcd5c2",
			Font: FontSerif, Radius: "0px", Texture: TextureNone,
			Vars: map[string]string{
				"--accent-2": "#8a6a2f",
				"--positive": "#3f7d4f",
				"--warn":     "#b06a1f",
			},
		},
	},
	PresetKraft: {
		Name: PresetKraft, Label: "牛皮纸",
		About: "牛皮纸底 + 深赭 + 横格信纸线，衬线标题。人文、历史、文学、独立杂志，需要年代感的场合",
		Theme: Theme{
			Accent: "#8f4520", Background: "#eedfc7",
			HeadingColor: "#2a1e13", TextColor: "#4a3826",
			Surface: "#f7ecd8", BorderColor: "#d3bf9d",
			// 横线纹理 = 老笔记本/信纸：这套调性里唯一能把"纸"讲清楚的那一种
			Font: FontEditorial, Radius: "2px", Texture: TextureRule,
			Vars: map[string]string{
				"--accent-2": "#4a6fa5",
				"--positive": "#2f6b4f",
				"--warn":     "#a8721f",
			},
		},
	},
	PresetDune: {
		Name: PresetDune, Label: "沙丘",
		About: "沙色底 + 炭灰 + 几何无衬线 + 大圆角，没有彩色强调。艺术、设计、画廊手册，审美优先、不想被配色抢戏的场合",
		Theme: Theme{
			Accent: "#2b241b", Background: "#f0e6d2",
			HeadingColor: "#1f1a14", TextColor: "#4a4238",
			Surface: "#faf4e8", BorderColor: "#ddd0b6",
			// 唯一一套"强调色不带彩"的预设。编辑风也是黑强调，但那套是纯白底 + 直角 + 现代无衬线；
			// 这套是沙色底 + 8px 圆角，落在"画廊手册 / 建筑图册"那一档
			Font: FontModern, Radius: "8px", Texture: TextureNone,
			Vars: map[string]string{
				"--accent-2": "#8a7a5c",
				"--positive": "#4f6b46",
				"--warn":     "#a8721f",
			},
		},
	},
	// 瑞士国际主义（Typographic Style）：中性灰白纸 + 单一高饱和色 + 近黑网格。
	// 这套预设的"参照物"不是某个文件，是一段有公开文献的平面设计史（Mueller-Brockmann 一脉）。
	// 取值按这个代码库的既有规矩独立挑选（外部同风格项目是 AGPL-3.0，不抄色板）：
	// 克莱因蓝是公开的色彩史料值，纸/墨/灰都是中性件。
	// 为什么它和编辑风不撞：编辑风强调色是黑、零彩色、衬线配对；swiss 强调色是全场唯一
	// 的饱和色（页面几乎全部是黑字白底，那一点蓝才值钱），且自带 Inter 细字重大字——
	// 一个像杂志评论页，一个像印刷展海报。
	PresetSwiss: {
		Name: PresetSwiss, Label: "瑞士风",
		About: "灰白纸 + 克莱因蓝单色 + 极细大字 + 等宽小标签，直角、细网格。品牌发布、设计理念、数据与观点并重的表达，要「印刷展」气质时用",
		Theme: Theme{
			Accent: "#002fa7", Background: "#fafaf8",
			HeadingColor: "#0a0a0a", TextColor: "#1a1a1a",
			Surface: "#f0f0ee", BorderColor: "#d4d4d2",
			Font: FontSwiss, Radius: "0px", Texture: TextureGrid,
			Vars: map[string]string{
				"--hero-weight": "200",   // Inter 细体：瑞士风的"极大字号 + 极轻字重"全靠它
				"--accent-2": "#0a0a0a", // 双色对比时：蓝 vs 墨，都是体系内的颜色
				"--positive": "#1a7f4b",
				"--warn":     "#b8860b",
			},
		},
	},
}

// presetOrder 是稳定顺序：map 遍历随机，而它出现在工具说明与报错信息里，
// 顺序变了 diff 与文档就跟着抖。
// 新增的四套排在六套之后：顺序决定了工具说明里清单的排布，插在中间会让已有 deck 的
// "第一眼看哪个"发生变化，而它们本来就是补充项。
var presetOrder = []string{
	PresetPaper, PresetEditorial, PresetNoir, PresetDuotone, PresetTerminal, PresetTech,
	PresetIndigo, PresetPine, PresetKraft, PresetDune, PresetSwiss,
}

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
