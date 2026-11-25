package response

import (
	"encoding/json"
	"log"
	"net/http"

	"deepseek-backend/middleware"
)

// ErrorInfo 错误响应内容
type ErrorInfo struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// Response 统一的响应体
type Response[T any] struct {
	Success bool       `json:"success"`
	Data    *T         `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

// Success 返回成功响应
func Success[T any](w http.ResponseWriter, data T) {
	writeJSON(w, http.StatusOK, Response[T]{
		Success: true,
		Data:    &data,
	})
}

// Fail 返回失败响应
func Fail(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	writeJSON(w, status, Response[any]{
		Success: false,
		Error: &ErrorInfo{
			Code:      code,
			Message:   message,
			RequestID: middleware.GetRequestID(r),
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("write response failed: %v", err)
	}
}
