package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoTemplatesDir 仓库自带的模板目录（templates/）；assetsDir 是 web/assets。
// 相对路径以包目录（backend/internal/service/template）为基准，
// `go test` 的 cwd 就是包目录，所以这样解析是稳定的。
func repoTemplatesDir(t *testing.T) (string, string) {
	t.Helper()
	return filepath.Clean(filepath.Join("..", "..", "..", "templates")),
		filepath.Clean(filepath.Join("..", "..", "..", "web", "assets"))
}

// 真实模板必须能通过注册校验：这是"入库校验器（mjs）与 Go 校验规则一致"的守门测试。
// 入库时跑 mjs 版，CI/启动跑 Go 版，两边规则漂移会在这里炸出来。
func TestLoadRepoTemplates(t *testing.T) {
	dir, assets := repoTemplatesDir(t)
	reg, err := NewRegistry(dir, assets)
	if err != nil {
		t.Fatalf("真实模板注册失败: %v", err)
	}
	if reg.Count() == 0 {
		t.Fatal("一个模板都没注册到：目录或路径配错了")
	}
	ts, err := reg.Get("tech-sharing")
	if err != nil {
		t.Fatalf("tech-sharing 未注册: %v", err)
	}
	if len(ts.Layouts) < 12 {
		t.Errorf("tech-sharing 版式数 %d < 12（上架标准）", len(ts.Layouts))
	}
	skel, ok := ts.Layout("cover")
	if !ok || !strings.Contains(skel, `data-layout="cover"`) {
		t.Errorf("cover 骨架缺失或没有 data-layout 标记:\n%s", skel)
	}
	if !strings.Contains(ts.Rules(), "禁 emoji") {
		t.Error("rules.md 应包含模板纪律")
	}
	// AllowedClasses：base 原语 + 模板专属类都该在
	allowed := ts.AllowedClasses("split-terminal")
	for _, c := range []string{"slide", "kicker", "terminal", "agenda-row"} {
		if !allowed[c] {
			t.Errorf("split-terminal 的合法类名缺 %q", c)
		}
	}
	if allowed["not-a-real-class"] {
		t.Error("AllowedClasses 返回了不存在的类")
	}
}

// 默认变体/未知变体/命名变体三条路。
func TestInstantiateVariants(t *testing.T) {
	dir, assets := repoTemplatesDir(t)
	reg, err := NewRegistry(dir, assets)
	if err != nil {
		t.Fatalf("注册失败: %v", err)
	}
	ts, _ := reg.Get("tech-sharing")

	ins, err := ts.Instantiate("", "Go 语言入门")
	if err != nil {
		t.Fatalf("默认变体实例化失败: %v", err)
	}
	if strings.Contains(ins.IndexHTML, "Rust 异步运行时") {
		t.Error("demo 页面没有被剥离")
	}
	if !strings.Contains(ins.IndexHTML, slidesStartMarker) || !strings.Contains(ins.IndexHTML, slidesEndMarker) {
		t.Error("实例化产物丢了挂载标记——deck 写路径将没有锚点")
	}
	if strings.Contains(ins.IndexHTML, "agenda-row") && strings.Count(ins.IndexHTML, "<section") > 0 {
		t.Error("剥离后不应残留任何 section")
	}
	if !strings.Contains(ins.IndexHTML, "<title>Go 语言入门</title>") {
		t.Error("标题没有被替换")
	}
	if strings.Contains(ins.IndexHTML, "v-blue") {
		t.Error("默认变体不该挂变体 class")
	}
	if ins.Canvas.W != 1920 || ins.Canvas.H != 1080 {
		t.Errorf("画布错误: %+v", ins.Canvas)
	}
	if ins.StyleCSS == "" || !strings.Contains(ins.StyleCSS, ".tpl-tech-sharing") {
		t.Error("style.css 未随实例化带出")
	}

	insBlue, err := ts.Instantiate("blue", "标题")
	if err != nil {
		t.Fatalf("blue 变体实例化失败: %v", err)
	}
	if !strings.Contains(insBlue.IndexHTML, `class="tpl-tech-sharing v-blue"`) {
		t.Errorf("blue 变体 class 没挂上 body:\n%s", bodyLine(insBlue.IndexHTML))
	}

	if _, err := ts.Instantiate("no-such-variant", "x"); err == nil {
		t.Error("未知变体应该报错")
	}
}

func bodyLine(doc string) string {
	for _, l := range strings.Split(doc, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "<body") {
			return l
		}
	}
	return "(no body tag)"
}

// ---------- 坏模板必须被拒绝 ----------

// fixture 在临时目录写一套最小但合法的模板，返回目录路径。
// base.css / runtime.js 的存在性校验用真实仓库的 assets 目录（repoTemplatesDir）。
func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	// template.json 的 id 必须与目录名一致（生产约束），所以夹具写进固定名字的子目录
	dir := filepath.Join(root, "fixture")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"template.json": `{
			"id":"fixture","name":"夹具","description":"d",
			"canvas":{"w":1920,"h":1080},
			"variants":[{"id":"default","name":"默认","class":""}],
			"layouts":[{"id":"cover","name":"封面","use":"u"}]}`,
		"index.html": `<!DOCTYPE html><html><head><title>old</title>
			<link rel="stylesheet" href="/assets/deck-v2/base.css"></head>
			<body class="tpl-fixture"><div class="deck">
			<!-- SLIDES:START -->
			<section class="slide" data-layout="cover"><h1 class="h1">demo</h1></section>
			<!-- SLIDES:END -->
			</div><script src="/assets/deck-v2/runtime.js"></script></body></html>`,
		"style.css":  `.tpl-fixture{--accent:#fff}`,
		"layouts.md": "## cover（封面）\n\n用途：x。\n\n```html\n<section class=\"slide\" data-layout=\"cover\"><h1 class=\"h1\">{{标题}}</h1></section>\n```\n\n合法类名：slide, h1\n",
		"rules.md":   "# 规则\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// overwrite 覆盖夹具的某个文件内容（制造缺陷用）。
func overwrite(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadRejectsBadTemplates(t *testing.T) {
	_, assets := repoTemplatesDir(t)

	good := fixture(t)
	if _, err := loadTemplate(good, assets, map[string]bool{"slide": true, "h1": true}); err != nil {
		t.Fatalf("基线夹具应当通过: %v", err)
	}

	cases := []struct {
		name   string
		mutate func(t *testing.T, dir string)
		want   string
	}{
		{"缺挂载标记", func(t *testing.T, d string) {
			overwrite(t, d, "index.html", strings.ReplaceAll(mustRead(t, d, "index.html"), slidesStartMarker, ""))
		}, "挂载标记"},
		{"断链资产", func(t *testing.T, d string) {
			overwrite(t, d, "index.html", strings.Replace(mustRead(t, d, "index.html"),
				"/assets/deck-v2/base.css", "/assets/deck-v2/not-exist.css", 1))
		}, "不存在的资产"},
		{"layouts.md 与 json 不一致", func(t *testing.T, d string) {
			overwrite(t, d, "layouts.md", strings.Replace(mustRead(t, d, "layouts.md"),
				"## cover", "## cover-x", 1))
		}, "无条目"},
		{"缺合法类名行", func(t *testing.T, d string) {
			overwrite(t, d, "layouts.md", strings.Replace(mustRead(t, d, "layouts.md"),
				"合法类名：slide, h1", "合法类名：", 1))
		}, "合法类名"},
		{"未知类名", func(t *testing.T, d string) {
			overwrite(t, d, "layouts.md", strings.Replace(mustRead(t, d, "layouts.md"),
				"合法类名：slide, h1", "合法类名：slide, h1, ghost-class", 1))
		}, "ghost-class"},
		{"缺骨架代码", func(t *testing.T, d string) {
			overwrite(t, d, "layouts.md", fenceRe.ReplaceAllString(mustRead(t, d, "layouts.md"), ""))
		}, "骨架"},
		{"变体 class 未定义", func(t *testing.T, d string) {
			overwrite(t, d, "template.json", strings.Replace(mustRead(t, d, "template.json"),
				`{"id":"default","name":"默认","class":""}`,
				`{"id":"default","name":"默认","class":""},{"id":"x","name":"X","class":"v-ghost"}`, 1))
		}, "无定义"},
		{"demo section 缺 data-layout", func(t *testing.T, d string) {
			overwrite(t, d, "index.html", strings.Replace(mustRead(t, d, "index.html"),
				`data-layout="cover"><h1`, `><h1`, 1))
		}, "data-layout"},
		{"id 与目录名不一致", func(t *testing.T, d string) {
			overwrite(t, d, "template.json", strings.Replace(mustRead(t, d, "template.json"),
				`"id":"fixture"`, `"id":"other"`, 1))
		}, "不一致"},
		{"数量行与骨架矛盾", func(t *testing.T, d string) {
			overwrite(t, d, "layouts.md", strings.Replace(mustRead(t, d, "layouts.md"),
				"合法类名：slide, h1", "合法类名：slide, h1\n数量：h1=2", 1))
		}, "矛盾"},
		{"数量行声明未登记的类", func(t *testing.T, d string) {
			overwrite(t, d, "layouts.md", strings.Replace(mustRead(t, d, "layouts.md"),
				"合法类名：slide, h1", "合法类名：slide, h1\n数量：ghost=1", 1))
		}, "未登记的类"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := fixture(t)
			tc.mutate(t, dir)
			_, err := loadTemplate(dir, assets, map[string]bool{"slide": true, "h1": true})
			if err == nil {
				t.Fatalf("应当被拒绝")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("错误信息缺关键字 %q:\n%v", tc.want, err)
			}
		})
	}
}

func mustRead(t *testing.T, dir, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// 「数量：」行 → 解析 → 骨架交叉校验 → Repeats 访问器的端到端检查。
// 真实模板用 product-launch（数量契约的首个用户），夹具验证 happy path。
func TestRepeatsContract(t *testing.T) {
	dir, assets := repoTemplatesDir(t)
	reg, err := NewRegistry(dir, assets)
	if err != nil {
		t.Fatalf("真实模板注册失败: %v", err)
	}
	pl, err := reg.Get("product-launch")
	if err != nil {
		t.Fatalf("product-launch 未注册: %v", err)
	}
	if got := pl.Repeats("how-it-works"); got["step"] != 3 {
		t.Errorf("how-it-works 的数量契约应为 step=3，得到 %v", got)
	}
	if got := pl.Repeats("feature-duo"); got["feature-card"] != 2 {
		t.Errorf("feature-duo 的数量契约应为 feature-card=2，得到 %v", got)
	}
	if got := pl.Repeats("pricing"); got["price-card"] != 3 {
		t.Errorf("pricing 的数量契约应为 price-card=3，得到 %v", got)
	}
	if got := pl.Repeats("cover"); len(got) != 0 {
		t.Errorf("cover 未登记数量契约，应返回空，得到 %v", got)
	}
	if got := pl.Repeats("no-such-layout"); got != nil {
		t.Errorf("未登记版式应返回 nil，得到 %v", got)
	}

	fix := fixture(t)
	tpl, err := loadTemplate(fix, assets, map[string]bool{"slide": true, "h1": true})
	if err != nil {
		t.Fatalf("夹具应通过: %v", err)
	}
	if got := tpl.Repeats("cover"); len(got) != 0 {
		t.Errorf("未写「数量：」行的版式应无契约，得到 %v", got)
	}
}
