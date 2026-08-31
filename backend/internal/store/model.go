package store

import "time"

// User 用户表。阶段5（免费额度 / 用户自带 API Key）启用，这里先占位建表。
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:64" json:"name"`
	APIKeyEnc string    `gorm:"size:512" json:"-"` // 自带 Key 的密文（AES），json:"-" 保证永不泄漏到响应
	CreatedAt time.Time `json:"created_at"`
}

// ChatSession 对话会话。阶段2 的 ask_user 暂停/恢复靠它保存完整消息列表：
// agent 循环暂停时把整个 messages 写进来，用户答题后读出来接着循环。
type ChatSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DeckID    string    `gorm:"size:32;index" json:"deck_id"`
	Title     string    `gorm:"size:128" json:"title"`
	Messages  string    `gorm:"type:longtext" json:"-"` // 完整消息列表的 JSON（阶段2启用）
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
