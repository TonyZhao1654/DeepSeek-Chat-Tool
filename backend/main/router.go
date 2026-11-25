package main

import (
	"net/http"

	"deepseek-backend/deepseek"
	"deepseek-backend/middleware"
	"deepseek-backend/session"
)

// 前端静态资源目录（相对backend/main的目录）
const frontendDir = "../../frontend"

// 组装各包暴露的路由
func buildRouter(deepseekHandler *deepseek.Handler, sessionHandler *session.Handler) http.Handler {
	mux := http.NewServeMux()

	deepseek.RegisterRoutes(mux, deepseekHandler)
	session.RegisterRoutes(mux, sessionHandler)

	// 前端静态资源
	mux.Handle("/", http.FileServer(http.Dir(frontendDir)))

	return middleware.RequestIDMiddleware(mux)
}
