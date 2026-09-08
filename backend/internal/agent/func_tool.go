package agent

import (
	"bytes"
	"context"
	"errors"
	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/service/deck"
	"strings"

	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
)

type ToolFunc func(ctx context.Context,arguments string)(string,error)
type Tool struct {
	Definition openai.ChatCompletionToolUnionParam //发给模型的schema
	Execute ToolFunc //执行函数
}

// toolUID 从 ctx 取当前登录用户。工具是"以某个用户的身份"执行的：
// 没有身份就拒绝——这是工具层的第一道归属闸门。LLM 的参数里
// 永远没有 user 概念（身份只来自认证注入，工具签名即权限边界）。
func toolUID(ctx context.Context) (uint, error) {
	uid, ok := authctx.UserID(ctx)
	if !ok {
		return 0, errors.New("未登录状态，无法操作 deck")
	}
	return uid, nil
}

type WriteDeckArgs struct {
	Title       string `json:"title" jsonschema:"required,type=string,description=演示文稿标题，也是浏览器标签页标题"`
	SectionHtml string `json:"section_html" jsonschema:"required,type=string,description=所有页面的 <section> HTML，按顺序拼接成一整个字符串"`
}

//创建新的deck
func (a *AgentService)toolWriteDeck(ctx context.Context,arguments string)(string,error){
	uid, err := toolUID(ctx)
	if err != nil {
		return "", err
	}

	var args WriteDeckArgs
	if err := json.Unmarshal([]byte(arguments),&args); err !=nil{
		return "",fmt.Errorf("write_deck 参数不是合法json:%v",err)
	}

	if strings.TrimSpace(args.Title) == ""{
		return "",errors.New("title 不可为空")
	}

	res,err := a.DeckService.Create(uid, args.Title, args.SectionHtml)
	if err !=nil{
		return "",err
	}

	return fmt.Sprintf(`{"deck_id":%q,"slides":%d,"url":"/api/decks/%s/file"}`,
		res.DeckID, res.Slides, res.DeckID), nil
}

type ListSlidesArgs struct {
	DeckID string `json:"deck_id" jsonschema:"required,type=string,description=要查看目录的 deck ID，例如 deck-0002"`
}

func (a *AgentService)toolListSlides(ctx context.Context,arguments string)(string,error){
	uid, err := toolUID(ctx)
	if err != nil {
		return "", err
	}

	var args ListSlidesArgs
	if err := json.Unmarshal([]byte(arguments),&args);err !=nil{
		return "",fmt.Errorf("list_slides 参数不是合法json:%v",err)
	}

	slides ,err := a.DeckService.ListSlides(uid, args.DeckID)
	if err !=nil{
		return "",err
	}

	if len(slides) == 0{
		return `{"slides":[]}`,nil
	}

	res,err := json.Marshal(slides)
	if err !=nil{
		return "",fmt.Errorf("json 序列化失败 err:%w",err)
	}

	return fmt.Sprintf(`{"slides":%s}`,res),nil
	
}

type ReadSlideArgs struct {
	DeckID string `json:"deck_id" jsonschema:"required,type=string,description=目标deck的ID，例如 deck-0002"`
	SlideID string `json:"slide_id" jsonschema:"required,type=string,description=要读取页的 slide_id (页面data-id)，必须来自 list_slides 的返回，例如 s3"`
}

type readSlideResult struct{
	DeckID string `json:"deck_id"`
	SlideID string `json:"slide_id"`
	Fingerprint string `json:"fingerprint"`
	HTML string `json:"html"`
}

func (a *AgentService)toolReadSlide(ctx context.Context,arguments string)(string,error){
	uid, err := toolUID(ctx)
	if err != nil {
		return "", err
	}

	var args ReadSlideArgs
	if err := json.Unmarshal([]byte(arguments),&args);err !=nil{
		return "",fmt.Errorf("read_slide 参数不是合法json,err:%v",err)
	}

	if !deck.IsValidID(args.DeckID){
		return "",fmt.Errorf("deck_id %q 不合法",args.DeckID)
	}
	if !deck.IsValidID(args.SlideID){
		return "",fmt.Errorf("slide_id %q 不合法",args.SlideID)
	}

	html,fp,err := a.DeckService.ReadSlide(uid, args.DeckID, args.SlideID)
	if err != nil{
		return "",err
	}

	res ,err := marshalNoEscape(readSlideResult{
		DeckID: args.DeckID,
		SlideID: args.SlideID,
		Fingerprint: fp,
		HTML: html,
	})

	if err !=nil{
		return "",fmt.Errorf("结果序列化json失败 err:%w",err)
	}

	return res,nil
}

func marshalNoEscape(v any) (string,error){
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(v);err!=nil{
		return "",err
	}

	return strings.TrimSuffix(buf.String(),"\n"),nil
}

type UpdateSlideArgs struct {
	DeckID string	`json:"deck_id" jsonschema:"required,type=string,description=目标 deck 的ID，例如 deck-0002"`
	SlideID string	`json:"slide_id" jsonschema:"required,type=string,description=要修改页的 slide_id， 来自 list_slides 或 read_slide"`
	NewHTML string	`json:"new_html" jsonschema:"required,type=string,description=替换后的完整 <section> 元素: 基于 read_slide 返回的内容修改，保留原 data-id"`
	Fingerprint string	`json:"fingerprint" jsonschema:"required,type=string,description=read_slide 返回的内容指纹，原样传入"`
}

func (a *AgentService) toolUpdateSlide(ctx context.Context,arguments string)(string,error){
	uid, err := toolUID(ctx)
	if err != nil {
		return "", err
	}

	var args UpdateSlideArgs
	if err := json.Unmarshal([]byte(arguments),&args); err !=nil{
		return "",fmt.Errorf("update_slide 参数不是合法json:%v",err)
	}

	if !deck.IsValidID(args.DeckID) {
		return "",fmt.Errorf("deck_id %q 不合法",args.DeckID)
	}
	if !deck.IsValidID(args.SlideID){
		return "",fmt.Errorf("slide_id %q 不合法",args.SlideID)
	}

	if err := a.DeckService.UpdateSlide(uid, args.DeckID, args.SlideID, args.NewHTML, args.Fingerprint);err !=nil{
		return "",err
	}

	return fmt.Sprintf(`{"slide_id":%q}`,args.SlideID),nil
}

type InsertSlideArgs struct {
	DeckID string `json:"deck_id" jsonschema:"required,type=string,description=目标 deck 的ID，例如 deck-0002"`
	AfterSlideID string `json:"after_slide_id" jsonschema:"required,type=string,description=插入位置：新页放到该 slide_id 的后面，或传 end 表示追加到末尾"`
	NewHTML string `json:"new_html" jsonschema:"required,type=string,description=新页的完整 <section> 元素，基于组件库 class 搭建。不要写 data-id（系统自动分配新编号）"`
}

func (a *AgentService) toolInsertSlide(ctx context.Context,arguments string)(string,error){
	uid, err := toolUID(ctx)
	if err != nil {
		return "", err
	}

	var args InsertSlideArgs
	if err := json.Unmarshal([]byte(arguments),&args); err !=nil{
		return "",fmt.Errorf("insert_slide 参数不是合法json:%v",err)
	}

	if !deck.IsValidID(args.DeckID){
		return "",fmt.Errorf("deck_id %q 不合法",args.DeckID)
	}
	if !deck.IsValidID(args.AfterSlideID){
		return "",fmt.Errorf("after_slide_id %q 不合法",args.AfterSlideID)
	}

	slideID,err := a.DeckService.InsertSlide(uid,args.DeckID,args.AfterSlideID,args.NewHTML)
	if err !=nil{
		return "",err
	}

	// 返回新页的 slide_id：LLM 后续 read_slide/update_slide 这一页时用得上
	return fmt.Sprintf(`{"slide_id":%q}`,slideID),nil
}

type DeleteSlideArgs struct {
	DeckID string `json:"deck_id" jsonschema:"required,type=string,description=目标 deck 的ID，例如 deck-0002"`
	SlideID string `json:"slide_id" jsonschema:"required,type=string,description=要删除的页的 slide_id，必须来自 list_slides 的返回"`
}

func (a *AgentService) toolDeleteSlide(ctx context.Context,arguments string)(string,error){
	uid, err := toolUID(ctx)
	if err != nil {
		return "", err
	}

	var args DeleteSlideArgs
	if err := json.Unmarshal([]byte(arguments),&args); err !=nil{
		return "",fmt.Errorf("delete_slide 参数不是合法json:%v",err)
	}

	if !deck.IsValidID(args.DeckID){
		return "",fmt.Errorf("deck_id %q 不合法",args.DeckID)
	}
	if !deck.IsValidID(args.SlideID){
		return "",fmt.Errorf("slide_id %q 不合法",args.SlideID)
	}

	if err := a.DeckService.DeleteSlide(uid,args.DeckID,args.SlideID);err !=nil{
		return "",err
	}

	return fmt.Sprintf(`{"slide_id":%q}`,args.SlideID),nil
}

type AskQuestion struct {
	Question string `json:"question" jsonschema:"required,type=string,description=问题本身，一句话说清要确认什么"`
	Options []string `json:"options,omitempty" jsonschema:"type=array,description=候选项，最多4个，把你最推荐的回答放在第一个"`
}

type AskUserArgs struct {
	Questions []AskQuestion `json:"questions" jsonschema:"required,type=array,description=要问用户的问题，1~6个，尽可能一次问完不要连环调用"`
}

func (a *AgentService) toolAskUser(ctx context.Context,arguments string)(string,error){
	return "",errors.New("ask_user 应由循环拦截暂停，不应执行到这里")
}

// generateSchema 把 Go struct 反射成 OpenAI 工具参数 schema。
// 启动时调用一次并缓存即可，不要放在请求路径上反射。
func generateSchema[T any]() openai.FunctionParameters {
	reflector := jsonschema.Reflector{
		RequiredFromJSONSchemaTags: true,  // 必填由 jsonschema:"required" 标签决定（默认规则是"无 omitempty 即必填"，太隐晦）
		AllowAdditionalProperties:  false, // 生成 additionalProperties:false，配合 strict 拒绝模型幻觉出的字段
		DoNotReference:             true,  // 嵌套 struct 内联展开，而不是拆到 $defs 再 $ref（模型对 $ref 支持参差）
	}
	s := reflector.Reflect(new(T))

	// 只取三件套重建顶层对象，丢掉 $schema/$id 等模型不关心的键
	// （openai-go 官方 README 的结构化输出示例就是同款做法）
	schema := map[string]any{
		"type":       "object",
		"properties": s.Properties, // 有序 map：模型看到的字段顺序和 struct 声明顺序一致
		// 拒绝 schema 之外的字段；开 strict 时这是硬性要求
		"additionalProperties": false,
	}
	if len(s.Required) > 0 {
		schema["required"] = s.Required
	}
	raw, err := json.Marshal(schema)
	if err != nil {
		panic(fmt.Sprintf("生成工具 schema 失败: %v", err))
	}
	var fp openai.FunctionParameters
	if err := json.Unmarshal(raw, &fp); err != nil {
		panic(err)
	}
	return fp
}

func (a *AgentService)buildTools () map[string]Tool{
	return map[string]Tool{
		"write_deck":{
			Definition: openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name: "write_deck",
				Description: openai.String("从零创建一整份新演示文稿。sections_html 中提供全部 <section> 内容，按页序拼接，不要写 data-id（系统自动编号）。若用户要修改已有 deck，禁止使用本工具（会覆盖重做），应使用页级编辑工具。"),
				Parameters: generateSchema[WriteDeckArgs](),
			}),
			Execute: a.toolWriteDeck,
		},
		"list_slides":{
			Definition: openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name: "list_slides",
				Description: openai.String("查看某份 deck 的页面目录（position=当前页序，id=该页的唯一标识，按页序排列）。position 会随插入/删除而变化，slide_id 永远不变。修改任何一页之前必须先用本工具确认准确的 slide_id，禁止凭记忆或按 id 数字大小猜测。"),
				Parameters: generateSchema[ListSlidesArgs](),
			}),
			Execute: a.toolListSlides,
		},
		"read_slide":{
			Definition: openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name: "read_slide",
				Description: openai.String("读取某一页的完整 HTML（整块 <section>）和内容指纹 fingerprint。用户要求查看或修改某一页的具体内容时使用。修改任何一页之前必须先用本工具拿到当前内容，禁止凭记忆修改；拿到的 fingerprint 之后要原样传给 update_slide。先用 list_slides 获取页清单和 slide_id。"),
				Parameters: generateSchema[ReadSlideArgs](),
			}),
			Execute: a.toolReadSlide,
		},
		"update_slide":{
			Definition: openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name: "update_slide",
				Description: openai.String("整块替换单页内容,工作流:先用 read_slide 拿到该页当前的 HTML 和 fingerprint,基于它修改后提交 new_html(完整<section>,保留原data-id),fingerprint 原样传入,一次只改一页,多页修改分多次调用;同一页改完如需再改,必须重新 read_slide,禁止使用本工具新建页"),
				Parameters: generateSchema[UpdateSlideArgs](),
			}),
			Execute: a.toolUpdateSlide,
		},
		"insert_slide":{
			Definition: openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name: "insert_slide",
				Description: openai.String("在某页之后插入一页新内容。after_slide_id 传某页的 slide_id（新页排在它后面）或 \"end\"（追加到末尾）。new_html 是完整 <section>，不要写 data-id（后端自动分配新编号，永不复用）。新增单页用本工具，禁止用 write_deck 整份重做。"),
				Parameters: generateSchema[InsertSlideArgs](),
			}),
			Execute: a.toolInsertSlide,
		},
		"delete_slide":{
			Definition: openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name: "delete_slide",
				Description: openai.String("删除一页。删除前必须先用 list_slides 确认要删的 slide_id。deck 至少保留一页，无法删空；被删页的编号不会复用，其余页的 data-id 保持不变。"),
				Parameters: generateSchema[DeleteSlideArgs](),
			}),
			Execute: a.toolDeleteSlide,
		},
		"ask_user":{
			Definition: openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name:"ask_user",
				Description:openai.String("向用户提问，缺少关键决策信息、或执行影响较大/不可逆的操作前使用。一次提问 1~6 个问题，可带候选项（将你最推荐的答案写在候选项第一个）；必须单独调用，不要与其他工具同时调用。提问后本轮会暂停，用户的回答会作为你的工具结果返回。"),
				Parameters: generateSchema[AskUserArgs](),
			}),
			Execute: a.toolAskUser,
		},
	}
}