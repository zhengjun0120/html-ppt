package deck

import (
	_ "embed"
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
)

var deckIDPattern = regexp.MustCompile(`^deck-(\d{4,})$`)

type CreateResult struct {
	DeckID string `json:"deck_id"`
	Slides int `json:"slides"`
}

func (s *Service) Create(title,sectionHTML string)(CreateResult,error){
	normalized,count,err := s.normalizeSections(sectionHTML)
	if err !=nil {
		return CreateResult{},err
	}

	//确保父目录存在
	if err := os.MkdirAll(s.decksDir,0o755); err !=nil{
		return CreateResult{},fmt.Errorf("创建文件夹失败 %w",err)
	}
	//创建唯一的子目录 并返回唯一的id
	id,err := s.claimDeckDir()
	if err != nil{
		return CreateResult{},err
	}

	rendered ,err := renderSkeleton(template.HTMLEscapeString(strings.TrimSpace(title)),normalized)
	if err !=nil{
		return CreateResult{},err
	}

	if err := s.atomicWriteDeck(id,rendered); err !=nil{
		return CreateResult{},fmt.Errorf("写入数据失败 :%w",err)
	}

	return CreateResult{DeckID: id,Slides: count},nil

}

// 校验llm提交的代码
func(s *Service) normalizeSections(sectionHTML string)(string,int ,error){
	var err error
	//空值检查
	raw := strings.TrimSpace(sectionHTML)
	if raw == ""{
		return "",0,errors.New("sections_html 不能为空")
	}

	//转小写检查是否有不符合的标签
	lower := strings.ToLower(raw)
	if strings.Contains(lower,"<!doctype") || strings.Contains(lower,"<html"){
		return "",0,errors.New("sections_html 只能是<section>元素的拼接，不要输出完整的HTML文档（<!DOCTYPE>/<html>/<head>/<body> 都不需要）")
	}

	// 解析html片段
	doc,err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err!=nil{
		return "",0,fmt.Errorf("解析 sections_html 失败 :%w",err)
	}

	// 拒绝脚本和样式
	if doc.Find("script").Length() > 0{
		return "",0,errors.New("幻灯片不支持 <script> ，请移除后重新提交")
	}
	if doc.Find("style").Length()>0{
		return "",0,errors.New("幻灯片不支持 <style>,请使用组件库 class,请移除后重新提交")
	}

	// 寻找body的子元素 section
	sections := doc.Find("body").Children().Filter("section")
	if sections.Length() == 0{
		return "",0,errors.New("未找到顶层<section>,请直接拼接<section>元素，不要用容器包裹")
	}

	//强制编号data-id

	var firstErr error
	var parts []string
	count := 0
	sections.Each(func(_ int, sec *goquery.Selection) {
		count++;
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
		return "",0,fmt.Errorf("提取html时 OuterHtml函数出错:%w",firstErr)
	}
	return strings.Join(parts,"\n"),count,nil
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

func renderSkeleton(escapedTitle,sections string) (string,error){
	var sb strings.Builder
	err := skeletonTmpl.Execute(&sb,struct{Title,Sections string}{escapedTitle,sections})

	if err !=nil{
		return "",fmt.Errorf("render skeleton: %w", err)
	}
	return sb.String(),nil
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
