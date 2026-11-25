package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// 把请求ID放进request的context
func WithRequestID(r *http.Request, id string) *http.Request {
	ctx := context.WithValue(r.Context(), requestIDKey, id)
	return r.WithContext(ctx)
}

// 从request的context取请求ID	相当于gi的c.GetString
func GetRequestID(r *http.Request) string {
	id, _ := r.Context().Value(requestIDKey).(string)
	return id
}

// 为每个请求生成或复用请求ID写入响应头和context
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = newRequestID()
		}

		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, WithRequestID(r, id))
	})
}

func newRequestID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}

	return hex.EncodeToString(b)
}
