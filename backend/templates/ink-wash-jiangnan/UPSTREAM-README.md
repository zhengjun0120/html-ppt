# ink-wash-jiangnan（水墨江南）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/ink-wash-jiangnan`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：宣纸米白 #F5F5DC / 淡雅象牙 #FAF7F0 底，浓墨 #1A1A1A 主字，柳叶嫩绿 #8FBC8F 与
  青苔绿 #6B8E6B 点缀，朱砂红 #DC143C 只用于印章、题款、关键数据，远山灰做中间调。
- 意象：拱桥、乌篷船、垂柳、白墙黑瓦、远山黛色、薄雾；装饰用水墨晕染、飞白、卷轴边框、印章。
- 排版：楷书/书法风大标题，宋体正文，行距宽松；留白即构图。
- 禁忌：高饱和荧光、赛博霓虹、金属质感、玻璃拟物、拥挤排版、现代无衬线粗体混搭。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 朱砂 #DC143C 收暗为 #C03A2E：投屏与截图管线下的观感更沉稳（原色偏信号红）。
- 远山与柳枝做成 CSS data-URI SVG 的 ambient 层（.slide::before/::after，z-index:-1），
  每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：标题 'STKaiti'/'KaiTi'（系统楷体），正文 Georgia/宋体栈，小字 Noto Sans SC；
  不新增 webfont 文件。
