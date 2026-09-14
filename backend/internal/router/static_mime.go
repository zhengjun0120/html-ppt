package router

import (
	"log"
	"mime"
)

// 字体文件的 MIME 注册。
//
// 为什么需要：Go 的 mime.TypeByExtension 在**部分平台**上查不到字体类型——
// Windows 上实测 .woff2/.woff/.ttf/.otf 四个都返回空串（它只读注册表里已有的
// Content Type，这几个通常没有）。查不到时 net/http 会退回内容嗅探，
// 而 woff2 的魔数（wOF2）不在 sniff 表里 → 响应头是 application/octet-stream。
// 浏览器对字体 MIME 的严格程度各家不同，但"用自己的字体服务返回 octet-stream"
// 是不该出现的状态：正确性不该取决于用户机器上装没装某个注册表项。
//
// 为什么用 init() 而不是在 New() 里调：mime 注册表是**进程级**的，
// 必须在任何静态响应之前完成；而 New() 的职责被限定为"URL → 方法的映射"。
// 放在 router 包里意味着：任何装配了静态路由的二进制（server、各种 demo）都自动带上。
func init() {
	for _, f := range []struct{ ext, typ string }{
		{".woff2", "font/woff2"},
		{".woff", "font/woff"},
		{".ttf", "font/ttf"},
		{".otf", "font/otf"},
	} {
		// 只在"原本查不到"时补：若某个平台的注册表已经给了更准确的值，不覆盖它
		if mime.TypeByExtension(f.ext) != "" {
			continue
		}
		if err := mime.AddExtensionType(f.ext, f.typ); err != nil {
			log.Printf("[warn] 注册 %s 的 MIME 类型失败: %v", f.ext, err)
		}
	}
}
