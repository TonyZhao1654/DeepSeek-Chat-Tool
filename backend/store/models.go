package store

// 会话	对应sessions表的一行
type Session struct {
	ID        int64          `json:"id"`
	Title     string         `json:"title"`
	Model     string         `json:"model"`
	Contents  []ContentBlock `json:"contents,omitempty"` // 全部内容块
	Messages  []Message      `json:"messages,omitempty"` // 全部消息
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
}

// 单个内容块	对应blocks表的一行
type ContentBlock struct {
	ID        int64  `json:"id"`
	SessionID int64  `json:"session_id,omitempty"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Selected  bool   `json:"selected"` // 聊天时按Selected决定是否注入上下文
	Position  int    `json:"position"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// 会话中的一条消息	对应messages表的一行
type Message struct {
	ID        int64  `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}
