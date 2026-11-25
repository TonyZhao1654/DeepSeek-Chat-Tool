// 前端接口契约（仅供类型检查，浏览器不会加载、也不需要编译）
//
// 命名规则：
//   - 对应后端 Request 的类型命名为 XxxRequest
//   - 对应后端 Response 的类型命名为 XxxPayload（即后端 XxxResponse 去掉 Response 换成 Payload）
//   - 只声明前端实际消费的字段，后端多返回的字段可以忽略
//
// 其他 JS 模块用 JSDoc 引用，例如 js/api 下：
//   /** @typedef {import("../types/api.js").SessionPayload} SessionPayload */

// 对应后端 response.Response[T]
export interface ResponsePayload<T> {
	success: boolean;
	data?: T;
	error?: ErrorInfoPayload;
}

// 对应后端 response.ErrorInfo
export interface ErrorInfoPayload {
	code: string;
	message: string;
	request_id?: string;
}

// 对应后端 store.Message
export interface MessagePayload {
	role: "system" | "user" | "assistant";
	content: string;
}

// 对应后端 store.ContentBlock
export interface ContentBlockPayload {
	id: number;
	title: string;
	content: string;
	selected: boolean;
}

// 对应后端 store.Session（列表接口不返回 contents / messages）
export interface SessionPayload {
	id: number;
	title: string;
	model: string;
	contents?: ContentBlockPayload[];
	messages?: MessagePayload[];
}

// 对应后端 session.CreateSessionRequest（空字段后端取默认值）
export interface CreateSessionRequest {
	title?: string;
	model?: string;
}

// 对应后端 session.UpdateSessionRequest（空字段后端保持原值）
export interface UpdateSessionRequest {
	title?: string;
	model?: string;
}

// 对应后端 session.ContentBlockRequest（selected 省略时：新建默认选中，更新保持原值）
export interface ContentBlockRequest {
	title: string;
	content: string;
	selected?: boolean;
}

// 对应后端 deepseek.ChatRequest（前端只发这三个字段，历史与内容块由后端补齐）
export interface ChatRequest {
	session_id: number;
	model: string;
	messages: MessagePayload[];
}

// 对应后端 deepseek.BalanceResponse
export interface BalancePayload {
	currency: string;
	total_balance: string;
}

// 对应后端 deepseek.APIModel（GET /api/models 的数组项）
export interface ModelPayload {
	id: string;
}

// 对应后端 deepseek.ChatUsage
export interface ChatUsagePayload {
	prompt_tokens: number;
	completion_tokens: number;
	total_tokens: number;
}

// 对应后端 deepseek.ChatChoice
export interface ChatChoicePayload {
	index: number;
	message: MessagePayload;
	finish_reason: string;
}

// 对应后端 deepseek.ChatResponse（SSE done 事件的载荷）
export interface ChatPayload {
	choices?: ChatChoicePayload[];
	usage?: ChatUsagePayload;
	balance?: BalancePayload;
	cost?: number;
	incomplete?: boolean;
	stream_error?: string;
}

// 对应后端透传的 deepseek.APIStreamChunk.Choices[].Delta
export interface ChatDeltaPayload {
	content?: string;
	reasoning_content?: string;
}

// 对应后端透传的 deepseek.APIStreamChunk.Choices[]
export interface ChatStreamChoicePayload {
	index?: number;
	delta?: ChatDeltaPayload;
	finish_reason?: string;
}

// 对应后端透传的 deepseek.APIStreamChunk（上游 SSE 的 data 块）
export interface ChatChunkPayload {
	choices?: ChatStreamChoicePayload[];
}
