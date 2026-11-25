package apperr

// 包一层，防止误解
type ErrorCode string

// 统一错误映射体，用于HTTPErrorMap
type HTTPError struct {
	Code    ErrorCode // 最终进响应体，给前端程序看
	Message string    // 最终进响应体，给前端用户看
	Status  int       // 最终进响应头，HTTP状态码
}

// HTTPErrorMap 错误到返回信息的映射
type HTTPErrorMap map[error]HTTPError
