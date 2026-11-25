package deepseek

import "net/http"

// RegisterRoutes 注册 DeepSeek 相关路由
func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /api/balance", handler.GetBalance)
	mux.HandleFunc("GET /api/models", handler.GetModels)
	mux.HandleFunc("POST /api/chat", handler.Chat)
}
