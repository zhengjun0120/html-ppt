package deck

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template" // 只为了 HTMLEscapeString；模板渲染在下面用 text/template
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	texttemplate "text/template"

	"github.com/PuerkitoBio/goquery"

	"html-ppt/backend/internal/store"
)

var deckIDPattern = regexp.MustCompile(`^deck-(\d{4,})$`)

type CreateResult struct {
	DeckID string `json:"deck_id"`
	Slides int    `json:"slides"`
	// Warning 给模型的回报，空串表示无违规。两类内容：
	// 消毒（"已移除 N 处内联事件属性"）与样式体检（见 stylelint.go）
	Warning string `json:"warning,omitempty"`
}

// Create 新建 deck。preset 非空时，骨架直接以该预设渲染——
// **"用户选了什么风格"必须在出生那一刻就落进主题**，而不是等 agent 之后补一次 update_theme：
// 后者的路径里没有任何机制防止"agent 答应了风格却忘调工具"，实测（deck-0022）结果是
// 用户选的 kraft 被静默丢掉、deck 停在默认纸感上——同质化最直接的一个来源。
func (s *Service) Create(userID uint, title, sectionHTML, preset string) (CreateResult, error) {
	if s.st == nil {
		return CreateResult{}, errStorage
	}

	normalized, count, warning, err := s.normalizeSections(sectionHTML)
	if err != nil {
		return CreateResult{}, err
	}

	theme := defaultTheme()
	if strings.TrimSpace(preset) != "" {
		p, ok := presets[strings.TrimSpace(preset)]
		if !ok {
			return CreateResult{}, fmt.Errorf("没有名为 %q 的预设。可用预设：%s",
				preset, strings.Join(PresetNames(), " / "))
		}
		p.applyTo(&theme)
	}

	// 骨架先渲染出来：它是纯函数（无副作用），而底下的回声校验要用它。
	rendered, err := renderSkeleton(template.HTMLEscapeString(strings.TrimSpace(title)), normalized, theme)
	if err != nil {
		return CreateResult{}, err
	}

	// 回声校验：内联 style 里引用的变量必须真有元素消费、或在这份 deck 自己的样式块里
	// 定义过。校验对象是**即将落盘的那份 HTML**（含新建时就渲染好的主题 override 块）——
	// 原来这里传的是空串，等于"新建的 deck 没有任何样式块"，于是预设预置的
	// --accent-2/--positive/--warn 会被判成无效，而提示词恰恰叫模型直接用它们。
	// 校验放在建目录之前：不合格的提交不该留下任何痕迹（连占号都不该占）。
	if err := s.checkInlineStyleVars(normalized, rendered); err != nil {
		return CreateResult{}, err
	}

	//确保父目录存在
	if err := os.MkdirAll(s.decksDir, 0o755); err != nil {
		return CreateResult{}, fmt.Errorf("创建文件夹失败 %w", err)
	}
	//创建唯一的子目录 并返回唯一的id
	id, err := s.claimDeckDir()
	if err != nil {
		return CreateResult{}, err
	}

	// 先写文件再登记归属：失败都回收整个目录，保证"列表看得到"和
	// "文件真实存在"同时成立，也避免失败残留的空目录永久占号
	if err := s.atomicWriteDeck(id, rendered); err != nil {
		os.RemoveAll(filepath.Join(s.decksDir, id))
		return CreateResult{}, fmt.Errorf("写入数据失败 :%w", err)
	}
	row := store.Deck{ID: id, UserID: userID, Title: strings.TrimSpace(title)}
	if err := s.st.DB.Create(&row).Error; err != nil {
		os.RemoveAll(filepath.Join(s.decksDir, id))
		return CreateResult{}, fmt.Errorf("登记 deck 归属: %w", err)
	}

	return CreateResult{DeckID: id, Slides: count, Warning: warning}, nil
}

// 校验llm提交的代码。返回归一化后的 sections、页数、消毒警告（可能为空串）。
func(s *Service) normalizeSections(sectionHTML string)(string,int,string,error){
	var err error
	//空值检查
	raw := strings.TrimSpace(sectionHTML)
	if raw == ""{
		return "",0,"",errors.New("sections_html 不能为空")
	}

	//转小写检查是否有不符合的标签
	lower := strings.ToLower(raw)
	if strings.Contains(lower,"<!doctype") || strings.Contains(lower,"<html"){
		return "",0,"",errors.New("sections_html 只能是<section>元素的拼接，不要输出完整的HTML文档（<!DOCTYPE>/<html>/<head>/<body> 都不需要）")
	}

	// 解析html片段
	doc,err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err!=nil{
		return "",0,"",fmt.Errorf("解析 sections_html 失败 :%w",err)
	}

	if doc.Find("section section").Length() > 0{
		return "",0,"",errors.New("不支持嵌套 <section> (垂直子页),每页只能是一个独立的顶层 <section>")
	}

	// 寻找body的子元素 section
	sections := doc.Find("body").Children().Filter("section")
	if sections.Length() == 0{
		return "",0,"",errors.New("未找到顶层<section>,请直接拼接<section>元素，不要用容器包裹")
	}

	//强制编号data-id
	// 消毒与编号在同一趟里做：先清理再序列化，剥掉的属性不会进入 parts
	var firstErr error
	var parts []string
	var merged sanitizeReport
	// 样式体检用同一个累加器跨页统计：页码、讲次提头这类"页面家具"每页只写一次，
	// 逐页看不出重复，只有整份一起数才知道"同一串 177 字符抄了 8 遍"
	lint := newStyleLinter()
	count := 0
	sections.Each(func(_ int, sec *goquery.Selection) {
		count++;
		merged.merge(sanitizeSlide(sec))
		lint.add(sec)
		sec.SetAttr("data-id",fmt.Sprintf("s%d",count))

		html,err := goquery.OuterHtml(sec);
		if err !=nil{
			if firstErr==nil{
				firstErr = err	
			}
			return
		}

		parts = append(parts, html)

	})
	if firstErr != nil{
		return "",0,"",fmt.Errorf("提取html时 OuterHtml函数出错:%w",firstErr)
	}
	// 结构性违规（危险标签）整份拒绝；属性级违规已剥除，警告随结果回给模型
	if err := merged.RejectErr(); err != nil{
		return "",0,"",err
	}
	return strings.Join(parts,"\n"),count,joinWarnings(merged.Warning(), lint.report().Warning()),nil
}

//计算出候选编号
func(s *Service) claimDeckDir() (string,error){
	for i:=0;i<100;i++{
		next,err := s.nextDeckNumber()
		if err !=nil{
			return "",err
		}

		id := fmt.Sprintf("deck-%04d",next)
		err = os.Mkdir(filepath.Join(s.decksDir,id),0o755)
		//名称被占用则重试
		if os.IsExist(err){
			continue
		}
		if err !=nil{
			return "",err
		}
		return id,nil
	}
	return "",errors.New("分配deck id 重试超限")
}

//扫描现有目录 返回最大编号+1
func(s *Service) nextDeckNumber()(int,error){
	entries,err := os.ReadDir(s.decksDir)
	if os.IsNotExist(err) {
		return 1,nil //一个文件都没有
	}

	if err!=nil{
		return 0,err
	}

	max := 0
	for _,e := range entries{
		m := deckIDPattern.FindStringSubmatch(e.Name())
		if m == nil{
			continue
		}
		if n,err := strconv.Atoi(m[1]);err == nil && n>max{
			max = n
		}
	}
	return max +1,nil
}
//go:embed deck_skeleton.html
var skeletonSrc string

var skeletonTmpl = texttemplate.Must(texttemplate.New("deck_skeleton").Parse(skeletonSrc))

// renderSkeleton 渲染骨架。ThemeJSON / ThemeCSS 都从传入的 theme 现渲染，
// 而不是在骨架里手抄一份：手抄的那份没有任何机制保证它与 defaultTheme() 一致，
// 而骨架真的会用到它们（新建的 deck 在第一次 update_theme 之前的主题就是它）——
// 这类"两份真相"迟早漂移。
func renderSkeleton(escapedTitle, sections string, theme Theme) (string, error) {
	themeJSON, err := json.Marshal(theme)
	if err != nil {
		return "", fmt.Errorf("序列化默认主题: %w", err)
	}
	var sb strings.Builder
	err = skeletonTmpl.Execute(&sb, struct{ Title, Sections, ThemeJSON, ThemeCSS string }{
		escapedTitle, sections, string(themeJSON), renderThemeCSS(theme),
	})

	if err != nil {
		return "", fmt.Errorf("render skeleton: %w", err)
	}
	return sb.String(), nil
}

func (s *Service) atomicWriteDeck(id,content string) error{
	//拼接路径
	finalPath := filepath.Join(s.decksDir,id,"deck.html")
	tmpPath := finalPath + ".tmp"
	if err := os.WriteFile(tmpPath,[]byte(content),0o644); err !=nil{
		return err
	}
	return os.Rename(tmpPath,finalPath)
	//之所以要先创建临时文件再重命名 是为了保证数据一致性，防止出现写入一半的情况
}
