# hand-drawn-autumn（手绘秋日旅行手账）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/hand-drawn-autumn`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：温暖米黄 #F5E6C8 打底，暖棕 #8B5E3C 做描边与文字骨干；砖红 #C04851 与
  枫叶橙 #E87D3E 做焦点高光，天空蓝 #6BA3BE 只在需要"清凉呼吸感"时出场；
  正文用更深的棕 #4A3728 保证可读。
- 母题：小房子（烟囱炊烟）、枫叶、小树、云朵、纸船、行李箱、热气球；手账装饰用
  星星、虚线路径、手绘圆圈、胶带式标签、对话框、便利贴。所有装饰半透明（0.2-0.5）。
- 排版：圆润手写/卡通标题放色块标签里像贴上去的胶带；虚线胶囊做副标题；
  层级靠"标签 + 色块"建立而不是字号差；布局刻意不规整、留白充足。
- 禁忌：冷几何、网格化规整对齐、渐变、3D、写实阴影、严肃衬线体、高饱和荧光色。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 手绘感全部用 CSS/SVG 实现：不规则 border-radius 模拟异形贴纸与手绘圈，
  ±0.6°-3° 的微旋转模拟"贴上去的层次感"，虚线 border 与 stroke-dasharray 做路线与勾边；
  不用任何外部图片，不加 webfont。
- 左上枫叶云朵与页底小屋虚线路由做成 CSS data-URI SVG 的 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：'Yuanti SC'/'YouYuan' 圆体打头（童趣手写感），回退 PingFang SC/Microsoft YaHei；
  不新增 webfont 文件。
