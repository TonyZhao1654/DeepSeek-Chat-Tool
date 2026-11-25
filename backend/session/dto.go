package session

import "deepseek-backend/store"

// 创建会话
type CreateSessionRequest struct {
	Title string `json:"title"`
	Model string `json:"model"`
}

// 更新会话标题或模型，空字段保持原值
type UpdateSessionRequest struct {
	Title string `json:"title"`
	Model string `json:"model"`
}

// 新建/更新内容块
type ContentBlockRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Selected *bool  `json:"selected"` // 为空时：新建默认选中，更新保持原值
}

// 发给模型前的上下文
type ChatContext struct {
	Model    string          `json:"model"`
	Messages []store.Message `json:"messages"`
}
