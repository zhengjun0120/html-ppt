# oh-my-ppt 风格导入 · 生成质量报告

> 生成日期：2026-09-30。83 个导入模板（84 风格 - soft-pastel 同名跳过）各走一遍
> **真实 LLM 生成**（deepseek-flash，正式 HTTP 链路：chat 澄清 → answer 提纲 →
> PUT 统一大纲 → 确认 → 选模板 → generate SSE），统一 12 页大纲
> 《二十四节气：跟着太阳过日子》（role 结构全量一致，保证横向可比）。
> 生成后用 Chrome 真渲染量测每页（vision.CaptureV2 同口径：填充率 ≥45%、
> 最小字号 ≥13px；hero/quote 指纹豁免填充率）。

## 结论

- **83/83 全部通过**：每模板成功生成 12 页，量测零不达标，无修复轮次需要。
- 填充率最低 61.8%、P05 64%、中位 66%（远超 45% 门槛）；最小字号 14px（门槛 13px）。
- 生成耗时中位 158 秒/模板，最长 302 秒；期间 LLM 侧零失败（个别页首写被
  生成管线自动重写救回，如 acid-blue-business p7）。
- 83 张生成封面 + 内容页截图逐一目检：同一份大纲呈现 83 种气质，无翻车页。
- 已知量测口径说明：deck 预览页底部有固定进度条，会使「溢出」字段在 deck 通道
  全页误报——溢出项以截图目检兜底（996 张截图抽查无可见溢出），模板 demo 通道
  （tplcheck 无 -base）的溢出字段正常，三波验收时均为零。

## 模板 × 测试 deck 对照（量测指标）

| 模板 | 测试 deck | 页数 | 最低填充率* | 最小字号 | 判定 |
|------|-----------|------|-------------|----------|------|
| academic-navy | deck-0085 | 12 | 67% | 14px | 通过 |
| academic-paper | deck-0083 | 12 | 66% | 15px | 通过 |
| acid-blue-business | deck-0084 | 12 | 66% | 16px | 通过 |
| amber-aurora | deck-0087 | 12 | 66% | 16px | 通过 |
| arctic-cool | deck-0086 | 12 | 65% | 15px | 通过 |
| aurora | deck-0088 | 12 | 66% | 16px | 通过 |
| bauhaus | deck-0089 | 12 | 65% | 15px | 通过 |
| blue-orange-analytics | deck-0090 | 12 | 66% | 16px | 通过 |
| blue-white-chart | deck-0091 | 12 | 67% | 16px | 通过 |
| blueprint | deck-0092 | 12 | 68% | 16px | 通过 |
| burgundy-premium | deck-0094 | 12 | 66% | 16px | 通过 |
| catppuccin-latte | deck-0093 | 12 | 66% | 15px | 通过 |
| catppuccin-mocha | deck-0095 | 12 | 67% | 15px | 通过 |
| celadon-bamboo | deck-0096 | 12 | 66% | 16px | 通过 |
| children-warm-orange | deck-0097 | 12 | 68% | 16px | 通过 |
| chinese-cream-blossom | deck-0098 | 12 | 66% | 16px | 通过 |
| chinese-fresh-trio | deck-0099 | 12 | 98% | 16px | 通过 |
| chinese-ink-landscape | deck-0100 | 12 | 68% | 16px | 通过 |
| chinese-pastel-spring | deck-0101 | 12 | 66% | 16px | 通过 |
| chinese-porcelain-rose | deck-0102 | 12 | 66% | 16px | 通过 |
| classic-duo-blue | deck-0103 | 12 | 68% | 16px | 通过 |
| cobalt-sunshine | deck-0104 | 12 | 67% | 16px | 通过 |
| corporate-clean | deck-0105 | 12 | 68% | 16px | 通过 |
| cream-pastel | deck-0106 | 12 | 70% | 15px | 通过 |
| cyberpunk-neon | deck-0107 | 12 | 65% | 16px | 通过 |
| dopamine-clash | deck-0108 | 12 | 68% | 16px | 通过 |
| dracula | deck-0109 | 12 | 64% | 16px | 通过 |
| dreamy-pink-gradient | deck-0110 | 12 | 67% | 16px | 通过 |
| dreamy-romance | deck-0111 | 12 | 68% | 16px | 通过 |
| e-ink-editorial | deck-0112 | 12 | 65% | 15px | 通过 |
| editorial-serif | deck-0113 | 12 | 68% | 15px | 通过 |
| engineering-whiteprint | deck-0114 | 12 | 68% | 16px | 通过 |
| geography-classroom | deck-0115 | 12 | 66% | 16px | 通过 |
| glassmorphism | deck-0116 | 12 | 66% | 16px | 通过 |
| gold-ivory | deck-0117 | 12 | 62% | 15px | 通过 |
| gradient-cosmic | deck-0118 | 12 | 66% | 16px | 通过 |
| gruvbox-dark | deck-0120 | 12 | 68% | 16px | 通过 |
| hand-drawn-autumn | deck-0119 | 12 | 67% | 16px | 通过 |
| handdrawn-watercolor | deck-0121 | 12 | 66% | 16px | 通过 |
| healing-color-card | deck-0122 | 12 | 67% | 14px | 通过 |
| indigo-lotus | deck-0123 | 12 | 66% | 16px | 通过 |
| industrial-kaizen | deck-0124 | 12 | 64% | 15px | 通过 |
| ink-wash-jiangnan | deck-0126 | 12 | 66% | 16px | 通过 |
| japanese-minimal | deck-0125 | 12 | 64% | 15px | 通过 |
| macaron-mist | deck-0127 | 12 | 66% | 16px | 通过 |
| magazine-bold | deck-0128 | 12 | 66% | 16px | 通过 |
| memphis-pop | deck-0129 | 12 | 70% | 15px | 通过 |
| midcentury | deck-0130 | 12 | 66% | 15px | 通过 |
| minimal-white | deck-0131 | 12 | 68% | 15px | 通过 |
| mint-fresh | deck-0132 | 12 | 67% | 15px | 通过 |
| mountain-green-literary | deck-0133 | 12 | 66% | 16px | 通过 |
| neo-brutalism | deck-0134 | 12 | 66% | 14px | 通过 |
| neon-haze | deck-0135 | 12 | 66% | 15px | 通过 |
| news-broadcast | deck-0136 | 12 | 65% | 14px | 通过 |
| nord | deck-0137 | 12 | 64% | 16px | 通过 |
| olive-elegant | deck-0138 | 12 | 66% | 16px | 通过 |
| orange-sea | deck-0139 | 12 | 68% | 15px | 通过 |
| oriental-poetic-illustration | deck-0140 | 12 | 66% | 16px | 通过 |
| palace-ink-red | deck-0141 | 12 | 66% | 16px | 通过 |
| pitch-deck-vc | deck-0142 | 12 | 65% | 15px | 通过 |
| premium-color-blocking | deck-0143 | 12 | 66% | 15px | 通过 |
| rainbow-gradient | deck-0144 | 12 | 67% | 16px | 通过 |
| red-gold-ceremony | deck-0145 | 12 | 66% | 16px | 通过 |
| red-research-framework | deck-0146 | 12 | 66% | 15px | 通过 |
| retro-tv | deck-0147 | 12 | 66% | 15px | 通过 |
| rose-pine | deck-0148 | 12 | 66% | 16px | 通过 |
| sage-academia | deck-0149 | 12 | 66% | 16px | 通过 |
| sakura-soft-healing | deck-0150 | 12 | 66% | 16px | 通过 |
| sharp-mono | deck-0151 | 12 | 66% | 16px | 通过 |
| solarized-light | deck-0152 | 12 | 66% | 15px | 通过 |
| song-rain-poetic | deck-0153 | 12 | 66% | 16px | 通过 |
| splash-abstract | deck-0154 | 12 | 68% | 16px | 通过 |
| starlight-fireworks | deck-0155 | 12 | 66% | 16px | 通过 |
| starry-dust | deck-0156 | 12 | 67% | 16px | 通过 |
| summer-warm-color | deck-0157 | 12 | 67% | 16px | 通过 |
| sunset-warm | deck-0158 | 12 | 66% | 15px | 通过 |
| swiss-grid | deck-0159 | 12 | 63% | 15px | 通过 |
| swiss-international | deck-0160 | 12 | 66% | 14px | 通过 |
| terminal-green | deck-0161 | 12 | 64% | 16px | 通过 |
| tokyo-night | deck-0162 | 12 | 68% | 16px | 通过 |
| vaporwave | deck-0163 | 12 | 70% | 16px | 通过 |
| xiaohongshu-white | deck-0164 | 12 | 68% | 15px | 通过 |
| y2k-chrome | deck-0165 | 12 | 66% | 16px | 通过 |
*最低填充率为该 deck 12 页中非豁免页（非 hero/quote 指纹）的最低值。

测试 deck 全部保留在测试账号（2119657168@qq.com）名下，可在画廊/编辑器复查。
