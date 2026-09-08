// Package authctx 定义"当前登录用户"在 context 里的传递约定。
// 单独一个小包而不是塞进 middleware：agent 的工具函数也要读它，
// middleware 不能被 agent 依赖（层次上是 agent 的上游），反向依赖会造成环。
package authctx

import "context"

type ctxKey struct{}

// WithUser 把登录用户的 id 注入 context。JWT 中间件在入口调用一次，
// 之后整个调用链（handler → agent → 工具 → deck 服务）都能取到。
func WithUser(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, ctxKey{}, userID)
}

// UserID 从 context 取登录用户。第二个返回值为 false 表示未登录
// （中间件没跑到、或上游忘了注入）——工具层据此拒绝操作。
func UserID(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(ctxKey{}).(uint)
	return id, ok
}
