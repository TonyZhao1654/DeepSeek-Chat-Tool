package deepseek

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"deepseek-backend/apperr"
	"deepseek-backend/response"
	"deepseek-backend/session"
	"deepseek-backend/store"
	"deepseek-backend/utils"
)

type Handler struct {
	service  *Service
	sessions *session.Service
}

func NewHandler(service *Service, sessions *session.Service) *Handler {
	handler := Handler{
		service:  service,
		sessions: sessions,
	}

	return &handler
}

// 默认模型
const defaultChatModel = "deepseek-flash"

// 查询账户余额
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, r, apperr.ErrMethodNotAllowed)

		return
	}

	balance, err := h.service.GetBalance()
	if err != nil {
		log.Printf("get balance failed: %v", err)

		handleError(w, r, err)

		return
	}

	response.Success(w, balance)
}

// 获取当前账号可用的模型列表
func (h *Handler) GetModels(w http.ResponseWriter, r *http.Request) {
	// 获取模型列表(切片)
	models, err := h.service.GetModels()
	if err != nil {
		log.Printf("get models failed: %v", err)

		handleError(w, r, err)

		return
	}

	response.Success(w, models)
}

// 对话
func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, r, apperr.ErrMethodNotAllowed)

		return
	}

	startTime := time.Now()

	var chatRequest ChatRequest
	err := json.NewDecoder(r.Body).Decode(&chatRequest)
	if err != nil {
		handleError(w, r, apperr.ErrInvalidRequest)

		return
	}

	if len(chatRequest.Messages) == 0 {
		handleError(w, r, apperr.ErrInvalidRequest)

		return
	}

	sessionID := chatRequest.SessionID
	newMessages := chatRequest.Messages

	// 带 session 时：后端拼上勾选的内容块和历史消息，前端只发新增内容
	if sessionID != 0 {
		chatContext, ctxErr := h.sessions.BuildChatContext(sessionID, toSessionMessages(newMessages))
		if ctxErr != nil {
			log.Printf("build chat context failed: %v", ctxErr)

			handleError(w, r, ctxErr)

			return
		}

		if chatRequest.Model == "" {
			chatRequest.Model = chatContext.Model
		}

		messages := make([]Message, 0, len(chatContext.Messages)+len(newMessages))
		for _, message := range chatContext.Messages {
			messages = append(messages, Message{Role: message.Role, Content: message.Content})
		}
		messages = append(messages, newMessages...)
		chatRequest.Messages = messages
	}

	if chatRequest.Model == "" {
		chatRequest.Model = defaultChatModel
	}

	// 采用流式对话，透传上游SSE，结束后补一个done事件带回余额与花费
	stream, err := h.service.StreamChat(r.Context(), chatRequest)
	if err != nil {
		log.Printf("chat failed: %v", err)

		handleError(w, r, err)

		return
	}
	defer stream.Close()

	// 上游已接受请求后再落库，避免请求失败时留下没有回复的用户消息
	if sessionID != 0 {
		appendErr := h.sessions.AppendMessages(sessionID, toSessionMessages(newMessages))
		if appendErr != nil {
			log.Printf("append user messages failed: %v", appendErr)

			handleError(w, r, appendErr)

			return
		}
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		handleError(w, r, apperr.ErrInternal)

		return
	}

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")

	fw := flushWriter{w: w, flusher: flusher}
	chatResponse, clientBroken := forwardStream(fw, stream)

	if sessionID != 0 && len(chatResponse.Choices) > 0 {
		reply := chatResponse.Choices[0].Message.Content
		if reply != "" {
			appendErr := h.sessions.AppendMessages(sessionID, []store.Message{
				{Role: "assistant", Content: reply},
			})
			if appendErr != nil {
				log.Printf("append assistant message failed: %v", appendErr)
			}
		}
	}

	chatResponse.Cost = utils.EstimateCost(chatResponse.Model, chatResponse.Usage.PromptTokens, chatResponse.Usage.CompletionTokens)

	logModel := chatResponse.Model
	if logModel == "" {
		logModel = chatRequest.Model
	}

	log.Printf(
		"chat completed: [%s] elapsed:%s input:%dtokens output:%dtokens cost:%.8fCNY",
		logModel,
		time.Since(startTime).Round(time.Millisecond),
		chatResponse.Usage.PromptTokens,
		chatResponse.Usage.CompletionTokens,
		chatResponse.Cost,
	)

	// 客户端已断开，不用再查余额和补done事件
	if clientBroken {
		return
	}

	// 获取余额
	balance, balanceErr := h.service.GetBalance()
	if balanceErr != nil {
		log.Printf("get balance after chat failed: %v", balanceErr)
	}

	chatResponse.Balance = balance

	// 序列化为json
	payload, err := json.Marshal(chatResponse)
	if err != nil {
		log.Printf("marshal chat response failed: %v", err)

		return
	}

	_, err = fmt.Fprintf(fw, "event: done\ndata: %s\n\n", payload)
	if err != nil {
		log.Printf("write done event failed: %v", err)
	}
}

// 每写一块就 flush，保证流式数据及时到达前端
type flushWriter struct {
	w       io.Writer
	flusher http.Flusher
}

func (f flushWriter) Write(p []byte) (int, error) {
	n, err := f.w.Write(p)
	f.flusher.Flush()

	return n, err
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	// 交给apperr包的HandleError统一处理，使用本包的httpErrorMap作为business
	apperr.HandleError(w, r, err, httpErrorMap)
}
