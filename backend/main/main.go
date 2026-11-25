package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"deepseek-backend/utils"
)

const (
	addr       = ":8080"
	browserURL = "http://localhost:8080"

	apiKeyBasePath = "../../"
	apiKeyFileName = "api.txt"
	dbPath         = "../../data/chat.db"
)

func main() {
	apiKey, err := utils.GetApiKey(apiKeyBasePath, apiKeyFileName)
	if err != nil {
		log.Fatalf("load api key failed: %v", err)
	}

	// 可能失败的用"成功"，不可能失败的用"完成"

	database, err := buildDB(dbPath)
	if err != nil {
		fmt.Println("create db.DB failed")
		log.Fatalf("init db failed: %v", err)
	}
	fmt.Println("created db.DB successfully")
	defer database.Close()

	sessionStore := buildSessionStore(database)
	fmt.Println("store.Store created")

	sessionService := buildSessionService(sessionStore)
	fmt.Println("session.Service created")

	sessionHandler := buildSessionHandler(sessionService)
	fmt.Println("session.Handler created")

	deepseekService := buildDeepSeekService(apiKey)
	fmt.Println("deepseek.Service created")

	deepseekHandler := buildDeepSeekHandler(deepseekService, sessionService)
	fmt.Println("deepseek.Handler created")

	router := buildRouter(deepseekHandler, sessionHandler)
	fmt.Println("app router created")

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	browserErr := utils.OpenBrowser(browserURL)
	if browserErr != nil {
		log.Printf("open browser failed: %v", browserErr)
	}

	log.Printf("server listening on %s", addr)

	err = http.Serve(listener, router)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
