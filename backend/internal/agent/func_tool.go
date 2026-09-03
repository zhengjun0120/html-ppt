package agent

import (
	"context"
	"errors"
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

type WriteDeckArgs struct {
	Title       string `json:"title" jsonschema:"required,type=string,description=演示文稿标题，也是浏览器标签页标题"`
	SectionHtml string `json:"section_html" jsonschema:"required,type=string,description=所有页面的 <section> HTML，按顺序拼接成一整个字符串"`
}

func (a *AgentService)toolWriteDeck(ctx context.Context,arguments string)(string,error){
	var args WriteDeckArgs
	if err := json.Unmarshal([]byte(arguments),&args); err !=nil{
		return "",fmt.Errorf("write_deck 参数不是合法json:%v",err)
	}

	if strings.TrimSpace(args.Title) == ""{
		return "",errors.New("title 不可为空")
	}

	res,err := a.DeckService.Create(args.Title,args.SectionHtml)
	if err !=nil{
		return "",err
	}

	return fmt.Sprintf(`{"deck_id":%q,"slides":%d,"url":"/api/decks/%s/file"}`,
		res.DeckID, res.Slides, res.DeckID), nil
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
	}
}