package session

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"deepseek-backend/apperr"
	"deepseek-backend/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	handler := Handler{
		service: service,
	}

	return &handler
}

// 创建会话
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil && !errors.Is(err, io.EOF) {
		handleError(w, r, apperr.ErrInvalidRequest)

		return
	}

	session, err := h.service.CreateSession(req)
	if err != nil {
		log.Printf("create session failed: %v", err)

		handleError(w, r, err)

		return
	}

	response.Success(w, session)
}

// 列出会话列表
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.service.ListSessions()
	if err != nil {
		log.Printf("list sessions failed: %v", err)

		handleError(w, r, err)

		return
	}

	response.Success(w, sessions)
}

// 获取会话详情（含内容块与消息）
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}

	session, err := h.service.GetSession(id)
	if err != nil {
		log.Printf("get session failed: %v", err)

		handleError(w, r, err)

		return
	}

	response.Success(w, session)
}

// 更新会话标题或模型，空字段保持原值
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req UpdateSessionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil && !errors.Is(err, io.EOF) {
		handleError(w, r, apperr.ErrInvalidRequest)

		return
	}

	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}

	session, err := h.service.UpdateSession(id, req)
	if err != nil {
		log.Printf("update session failed: %v", err)

		handleError(w, r, err)

		return
	}

	response.Success(w, session)
}

// 删除会话
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}

	err := h.service.DeleteSession(id)
	if err != nil {
		log.Printf("delete session failed: %v", err)

		handleError(w, r, err)

		return
	}

	response.Success(w, struct{}{})
}

// 新增内容块
func (h *Handler) CreateContent(w http.ResponseWriter, r *http.Request) {
	var req ContentBlockRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, r, apperr.ErrInvalidRequest)

		return
	}

	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}

	block, err := h.service.CreateContent(id, req)
	if err != nil {
		log.Printf("create content block failed: %v", err)

		handleError(w, r, err)

		return
	}

	response.Success(w, block)
}

// 整体更新内容块（标题、内容、是否注入）
func (h *Handler) UpdateContent(w http.ResponseWriter, r *http.Request) {
	var req ContentBlockRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, r, apperr.ErrInvalidRequest)

		return
	}

	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}

	contentID, ok := parseID(w, r, "contentID")
	if !ok {
		return
	}

	block, err := h.service.UpdateContent(id, contentID, req)
	if err != nil {
		log.Printf("update content block failed: %v", err)

		handleError(w, r, err)

		return
	}

	response.Success(w, block)
}

// 删除内容块
func (h *Handler) DeleteContentBlock(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}

	contentID, ok := parseID(w, r, "contentID")
	if !ok {
		return
	}

	err := h.service.DeleteContentBlock(id, contentID)
	if err != nil {
		log.Printf("delete content block failed: %v", err)

		handleError(w, r, err)

		return
	}

	response.Success(w, struct{}{})
}

// parseID 解析路径参数中的数字 ID，非法时按参数错误返回
func parseID(w http.ResponseWriter, r *http.Request, name string) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue(name), 10, 64)
	if err != nil || id <= 0 {
		handleError(w, r, apperr.ErrInvalidRequest)

		return 0, false
	}

	return id, true
}

// handleError 用本包的错误映射返回错误
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	// 交给apperr包的HandleError统一处理，以本包的httpErrorMap作为business
	apperr.HandleError(w, r, err, httpErrorMap)
}
