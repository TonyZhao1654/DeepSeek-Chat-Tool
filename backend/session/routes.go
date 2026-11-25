package session

import "net/http"

// 注册会话与内容块相关路由
func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /api/sessions", handler.Create)
	mux.HandleFunc("GET /api/sessions", handler.List)
	mux.HandleFunc("GET /api/sessions/{id}", handler.Get)
	mux.HandleFunc("PUT /api/sessions/{id}", handler.Update)
	mux.HandleFunc("DELETE /api/sessions/{id}", handler.Delete)

	mux.HandleFunc("POST /api/sessions/{id}/contents", handler.CreateContent)
	mux.HandleFunc("PUT /api/sessions/{id}/contents/{contentID}", handler.UpdateContent)
	mux.HandleFunc("DELETE /api/sessions/{id}/contents/{contentID}", handler.DeleteContentBlock)
}
