package session

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"deepseek-backend/store"
)

const (
	defaultModel        = "deepseek-flash"
	defaultTitle        = "新会话"
	defaultContentTitle = "未命名内容块"

	// 发给模型的上下文字符预算，超出后丢弃最早的历史消息
	maxContextChars = 40000
)

type Service struct {
	store *store.Store
}

func NewService(s *store.Store) *Service {
	return &Service{store: s}
}

func (s *Service) CreateSession(request CreateSessionRequest) (*store.Session, error) {
	now := time.Now().Unix()

	session := &store.Session{
		Title:     request.Title,
		Model:     request.Model,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 这个判断应该给前端的，这里反复验证一下
	if session.Title == "" {
		session.Title = defaultTitle
	}

	if session.Model == "" {
		session.Model = defaultModel
	}

	err := s.store.CreateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) ListSessions() ([]store.Session, error) {
	sessions, err := s.store.ListSessions()
	if err != nil {
		// store包不进行错误包装
		return nil, fmt.Errorf("session: list sessions: %w", err)
	}

	return sessions, nil
}

func (s *Service) GetSession(id int64) (*store.Session, error) {
	session, err := s.store.GetSession(id)
	if err != nil {
		// store包不进行错误包装
		return nil, fmt.Errorf("session: get session: %w", err)
	}

	return session, nil
}

// 更新会话标题或模型，空字段保持原值
func (s *Service) UpdateSession(id int64, request UpdateSessionRequest) (*store.Session, error) {
	session, err := s.store.GetSessionMeta(id)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(request.Title)
	if title != "" {
		session.Title = title
	}

	model := strings.TrimSpace(request.Model)
	if model != "" {
		session.Model = model
	}

	now := time.Now().Unix()
	err = s.store.UpdateSession(id, session.Title, session.Model, now)
	if err != nil {
		return nil, err
	}

	session.UpdatedAt = now

	return session, nil
}

func (s *Service) DeleteSession(id int64) error {
	err := s.store.DeleteSession(id)
	if err != nil {
		// store包不进行错误包装
		return fmt.Errorf("session: delete session: %w", err)
	}

	return nil
}

func (s *Service) CreateContent(sessionID int64, request ContentBlockRequest) (*store.ContentBlock, error) {
	exists, err := s.store.SessionExists(sessionID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, store.ErrSessionNotFound
	}

	position, err := s.store.NextContentPosition(sessionID)
	if err != nil {
		return nil, err
	}

	selected := true
	if request.Selected != nil {
		selected = *request.Selected
	}

	now := time.Now().Unix()
	block := &store.ContentBlock{
		SessionID: sessionID,
		Title:     normalizeTitle(request.Title),
		Content:   request.Content,
		Selected:  selected,
		Position:  position,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.store.CreateContent(block)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (s *Service) UpdateContent(sessionID, contentID int64, request ContentBlockRequest) (*store.ContentBlock, error) {
	block, err := s.store.GetContent(sessionID, contentID)
	if err != nil {
		return nil, err
	}

	block.Title = normalizeTitle(request.Title)
	block.Content = request.Content
	block.UpdatedAt = time.Now().Unix()

	if request.Selected != nil {
		block.Selected = *request.Selected
	}

	err = s.store.UpdateContent(block)
	if err != nil {
		return nil, err
	}

	return block, nil
}

// 删除内容块
func (s *Service) DeleteContentBlock(sessionID, contentID int64) error {
	err := s.store.DeleteContentBlock(sessionID, contentID)
	if err != nil {
		// store包不进行错误包装
		return fmt.Errorf("session: delete content: %w", err)
	}

	return nil
}

// 追加消息并刷新会话更新时间
func (s *Service) AppendMessages(sessionID int64, messages []store.Message) error {
	exists, err := s.store.SessionExists(sessionID)
	if err != nil {
		return err
	}

	if !exists {
		return store.ErrSessionNotFound
	}

	return s.store.AppendMessages(sessionID, messages)
}

// 组装发给模型的上下文：勾选的内容块放 system，随后是历史消息
func (s *Service) BuildChatContext(sessionID int64, incoming []store.Message) (*ChatContext, error) {
	session, err := s.store.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	messages := make([]store.Message, 0, len(session.Messages)+1)

	systemPrompt := buildSystemPrompt(session.Contents)
	if systemPrompt != "" {
		messages = append(messages, store.Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	// 内容块和本轮新消息都占预算，剩下的才留给历史
	reservedChars := utf8.RuneCountInString(systemPrompt)
	for _, message := range incoming {
		reservedChars += utf8.RuneCountInString(message.Content)
	}

	messages = append(messages, trimHistory(reservedChars, session.Messages)...)

	return &ChatContext{
		Model:    session.Model,
		Messages: messages,
	}, nil
}
