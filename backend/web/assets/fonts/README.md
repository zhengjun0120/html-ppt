# web/assets/fonts —— 随附字体

这个目录里放的是**我们自己服务的字体**（不依赖 CDN、不依赖用户机器上装了什么）。
字体栈里的其余四对仍然走系统字体栈，理由见 `theme.css` 顶部与
`internal/service/deck/theme.go` 的 `fontPairs` 注释：只有"中日文等宽"这一档
必须自带才成立，其余几档靠系统字体就能得到可预期的结果。

## JetBrains Maple Mono

| | |
|---|---|
| 文件 | `jetbrains-maple-mono/JetBrainsMapleMono-{Regular,Bold,Italic}.woff2` |
| 来源 | [SpaceTimee/Fusion-JetBrainsMapleMono](https://github.com/SpaceTimee/Fusion-JetBrainsMapleMono)（JetBrains Mono 的拉丁 + Maple Mono 的中日韩，自动跟踪上游构建） |
| 授权 | SIL Open Font License 1.1，全文见 `jetbrains-maple-mono/OFL.txt` |
| 版权 | JetBrains Mono Project Authors (2020)、Maple Mono Project Authors (2022, subframe7536，保留字体名 Maple Mono)、Space Time (2025) |

OFL-1.1 允许商用与嵌入（OFL 的 FAQ 明确把 web 字体列为允许用途），要求随附授权文本
（所以 `OFL.txt` 必须和字体放在一起，不能删）。若将来要**修改**字体文件再分发，
OFL 的保留字体名条款要求改名——不改就不用管。

三个文件各约 6.5MB，一共约 19MB。它们大是因为**含完整的中日韩字型**
（只含拉丁的等宽字体 woff2 通常只有几十 KB）——这份体积正是我们在这一档
自带字体所换来的东西。

## 体积与将来的子集化

现在不做子集化（subsetting），原因是**内容在运行时才生成**：deck 里的字由模型写，
一份静态子集必然对没预料到的字缺字（豆腐块），而"哪些字会出现"在读到一个 deck 之前
根本不知道。按 deck 实际用到的字符在**导出时**做子集化才是对的做法，
那件事属于"导出单文件"功能（尚未实现）。

运行时这条路径不需要子集化：字体从本机 /assets 静态服务，浏览器**只在排版真的用到
某个 `@font-face` 时才下载对应文件**——没用到等宽、也没写 `<code>` 的 deck
一个字节都不会下载。
