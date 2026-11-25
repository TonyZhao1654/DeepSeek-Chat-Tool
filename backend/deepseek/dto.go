package deepseek

// 后端请求api的（后端请求DeepSeek的）
// ===================================================================
type APIBalanceResponse struct {
	IsAvailable  bool             `json:"is_available"`
	BalanceInfos []APIBalanceInfo `json:"balance_infos"`
}

type APIBalanceInfo struct {
	Currency     string `json:"currency"`          // 货币类型（CNY）
	TotalBalance string `json:"total_balance"`     // 总可用余额
	Granted      string `json:"granted_balance"`   // 赠送金额
	ToppedUp     string `json:"topped_up_balance"` // 充值金额
}

// 账号可用的模型（单个）
type APIModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

type APIModelsResponse struct {
	Object string     `json:"object"`
	Data   []APIModel `json:"data"` // 账号可用的模型（多个）
}

type APIStreamError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    any    `json:"code,omitempty"`
}

// 流式响应块，用于解析上游 SSE
type APIStreamChunk struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage ChatUsage       `json:"usage"`
	Error *APIStreamError `json:"error,omitempty"`
}

// 发给DeepSeek的流式请求体，只包含上游认识的字段
type APIStreamRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream"`
	Tools       []Tool    `json:"tools,omitempty"`
}

type Tool struct {
	Type string `json:"type"`
}

// 前端请求后端的
// ===================================================================
type BalanceResponse struct {
	Currency     string `json:"currency"`
	TotalBalance string `json:"total_balance"`
}

type ChatRequest struct {
	SessionID   int64     `json:"session_id,omitempty"`
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream"`               // 是否用流式输出
	EnableWeb   bool      `json:"enable_web,omitempty"` // 是否允许联网
}

type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type ChatResponse struct {
	ID          string           `json:"id"`
	Object      string           `json:"object"`
	Created     int64            `json:"created"`
	Model       string           `json:"model"`
	Choices     []ChatChoice     `json:"choices"`
	Usage       ChatUsage        `json:"usage"`
	Balance     *BalanceResponse `json:"balance,omitempty"`
	Cost        float64          `json:"cost"`
	Incomplete  bool             `json:"incomplete,omitempty"` // 表示流式回答异常结束（上游报错、连接中断或没有 finish_reason）
	StreamError string           `json:"stream_error,omitempty"`
}
