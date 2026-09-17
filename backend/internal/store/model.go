package store

import "time"

// User 用户表。邮箱即身份；密码只存 bcrypt 哈希（不可逆），
// 自带 API Key 存 AES-GCM 密文（可逆——后端要解出来替用户调 LLM）。
type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Email        string `gorm:"uniqueIndex;size:255" json:"email"`
	PasswordHash string `gorm:"size:255" json:"-"`
	APIKeyEnc    string `gorm:"size:1024" json:"-"` // base64(nonce|密文)；json:"-" 保证永不进响应
	// 上次换发 token 的时刻（auth.Refresh 记账用）：NULL 或早于今天零点 = 今天还有刷新额度。
	// 用指针是刻意的：'从未刷过' 必须是 NULL 而不是 Go 零值——MySQL 严格模式拒收零值日期
	TokenRefreshedAt *time.Time `json:"-"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Deck 归属表：deck 文件在磁盘上，这张表负责 id → 归属用户 的权威映射。
// 所有"这个 deck 是谁的"判断只认这里；越权访问统一按"不存在"处理（不泄露存在性）。
//
// v2 新增四列（deck-v2 重构）：Format 区分新旧格式（v1 行在列表里隐藏，不迁移）；
// Stage 是生成流程状态机（draft→…→iterating，权威在 DB，deck.json 里是冗余副本）；
// TemplateID/Variant 记录所选模板与主题变体（实例化后不再变更）。
type Deck struct {
	ID         string    `gorm:"primaryKey;size:32" json:"id"` // deck-0006
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	Title      string    `gorm:"size:255" json:"title"`
	Format     string    `gorm:"index;size:8;default:v1" json:"format"`
	Stage      string    `gorm:"index;size:32;default:''" json:"stage"`
	TemplateID string    `gorm:"size:64;default:''" json:"template_id"`
	Variant    string    `gorm:"size:64;default:''" json:"variant"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ChatSession 对话会话。阶段2 的 ask_user 暂停/恢复靠它保存完整消息列表：
// agent 循环暂停时把整个 messages 写进来，用户答题后读出来接着循环。
type ChatSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	DeckID    string    `gorm:"size:32;index" json:"deck_id"`
	PendingAsk string `gorm:"size:1024" json:"-"`  //非空 = 循环暂停等待回答
	Title     string    `gorm:"size:128" json:"title"`
	Messages  string    `gorm:"type:longtext" json:"-"` // 完整消息列表的 JSON（阶段2启用）
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
