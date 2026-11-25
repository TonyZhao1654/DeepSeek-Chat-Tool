package main

import (
	"deepseek-backend/db"
	"deepseek-backend/deepseek"
	"deepseek-backend/session"
	"deepseek-backend/store"
)

// 构造db.DB实例
func buildDB(path string) (*db.DB, error) {
	database, err := db.NewDB(path)
	if err != nil {
		return nil, err // 下层已经包装好错误了
	}

	return database, nil
}

// 构造store.Store实例
func buildSessionStore(database *db.DB) *store.Store {
	s := store.NewStore(database)

	return s
}

// 构造session.Service实例
func buildSessionService(sessionStore *store.Store) *session.Service {
	sessionService := session.NewService(sessionStore)

	return sessionService
}

// 构造deepseek.Service实例
func buildDeepSeekService(apiKey string) *deepseek.Service {
	deepseekService := deepseek.NewService(apiKey)

	return deepseekService
}

// 构造deepseek.Handler实例
func buildDeepSeekHandler(deepseekService *deepseek.Service, sessionService *session.Service) *deepseek.Handler {
	deepseekHandler := deepseek.NewHandler(deepseekService, sessionService)

	return deepseekHandler
}

// 构造session.Handler实例
func buildSessionHandler(sessionService *session.Service) *session.Handler {
	sessionHandler := session.NewHandler(sessionService)

	return sessionHandler
}
