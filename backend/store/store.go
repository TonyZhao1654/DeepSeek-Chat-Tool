package store

/*
	store包是session包的附属包
	不合并是因为我不想让session包接触SQL语句
	store.Store结构体的方法的调用方只有session包的service.go或store包自身

	在用"`xxx`"定义sql语句时，可以理解为"换行 = 空格"

	store包不进行错误包装
*/

import (
	"database/sql"
	"errors"
	"time"

	"deepseek-backend/db"
)

type Store struct {
	db *db.DB
}

func NewStore(database *db.DB) *Store {
	store := Store{
		db: database,
	}

	return &store
}

// 插入会话
func (s *Store) CreateSession(session *Session) error {
	const sql = `INSERT INTO sessions
	(title, model, created_at, updated_at)
	VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(sql, session.Title, session.Model, session.CreatedAt, session.UpdatedAt)
	if err != nil {
		return err
	}

	session.ID, err = result.LastInsertId() // ID由数据库自增分配后回填
	if err != nil {
		return err
	}

	/*
		如果插入成功
		result.LastInsertId()的err会返回nil，也就是这里的err会是nil
		因此可以直接写成return err
		而不是if err != nil {}; return nil
		但我更习惯这样写
	*/

	return nil
}

func (s *Store) ListSessions() ([]Session, error) {
	const sql = `SELECT
	id, title, model, created_at, updated_at
	FROM sessions
	ORDER BY updated_at DESC`

	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]Session, 0)
	for rows.Next() {
		var item Session
		err = rows.Scan(&item.ID, &item.Title, &item.Model, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, item)
	}

	return sessions, rows.Err()
}

func (s *Store) GetSessionMeta(id int64) (*Session, error) {
	// 避免命名冲突而不叫sql
	const query = `SELECT
	id, title, model, created_at, updated_at
	FROM sessions
	WHERE id = ?`

	var session Session

	err := s.db.QueryRow(query, id).Scan(&session.ID, &session.Title, &session.Model, &session.CreatedAt, &session.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Store) GetSession(id int64) (*Session, error) {
	session, err := s.GetSessionMeta(id)
	if err != nil {
		return nil, err
	}

	contents, err := s.ListContents(id)
	if err != nil {
		return nil, err
	}
	session.Contents = contents

	messages, err := s.ListMessages(id)
	if err != nil {
		return nil, err
	}
	session.Messages = messages

	return session, nil
}

// 删除会话
func (s *Store) DeleteSession(id int64) error {
	const sql = `DELETE
	FROM sessions
	WHERE id = ?`

	result, err := s.db.Exec(sql, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrSessionNotFound
	}

	return nil
}

// 检查会话是否存在
func (s *Store) SessionExists(id int64) (bool, error) {
	// 避免命名冲突而不叫sql
	const query = `SELECT 1
	FROM sessions
	WHERE id = ?`

	var exists int
	err := s.db.QueryRow(query, id).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// 更新会话
func (s *Store) UpdateSession(id int64, title, model string, updatedAt int64) error {
	const sql = `UPDATE sessions
	SET title = ?, model = ?, updated_at = ?
	WHERE id = ?`

	result, err := s.db.Exec(sql, title, model, updatedAt, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrSessionNotFound
	}

	return nil
}

func (s *Store) NextContentPosition(sessionID int64) (int, error) {
	const sql = `SELECT COALESCE(MAX(position), 0) + 1
	FROM blocks
	WHERE session_id = ?`

	var position int

	err := s.db.QueryRow(sql, sessionID).Scan(&position)

	return position, err
}

// 创建内容块（插入内容块）
func (s *Store) CreateContent(block *ContentBlock) error {
	const sql = `INSERT INTO blocks
	(session_id, title, content,
	selected, position, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := s.db.Exec(
		sql,
		block.SessionID, block.Title, block.Content,
		block.Selected, block.Position, block.CreatedAt, block.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// ID由数据库自增分配后回填
	block.ID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return s.touchSession(block.SessionID)
}

// 更新内容块
func (s *Store) UpdateContent(block *ContentBlock) error {
	const sql = `
		UPDATE blocks
		SET title = ?, content = ?, selected = ?, updated_at = ?
		WHERE id = ? AND session_id = ?
	`

	result, err := s.db.Exec(
		sql,
		block.Title, block.Content, block.Selected, block.UpdatedAt,
		block.ID, block.SessionID,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrContentNotFound
	}

	err = s.touchSession(block.SessionID)
	if err != nil {
		return err
	}

	return nil
}

// 删除内容块（根据会话id和内容块id）
func (s *Store) DeleteContentBlock(sessionID, contentID int64) error {
	const sql = `DELETE FROM blocks
	WHERE id = ? AND session_id = ?`

	result, err := s.db.Exec(sql, contentID, sessionID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrContentNotFound
	}

	err = s.touchSession(sessionID)
	if err != nil {
		return err
	}

	return nil
}

// 查询某个内容块的内容（根据会话id和内容块id）
func (s *Store) GetContent(sessionID, contentID int64) (*ContentBlock, error) {
	// 避免命名冲突而不叫sql
	const query = `SELECT id, session_id, title, content,
	selected, position, created_at, updated_at
	FROM blocks WHERE id = ? AND session_id = ?`

	var block ContentBlock

	err := s.db.QueryRow(
		query, contentID, sessionID,
	).Scan(
		&block.ID, &block.SessionID, &block.Title, &block.Content,
		&block.Selected, &block.Position, &block.CreatedAt, &block.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrContentNotFound
	}
	if err != nil {
		return nil, err
	}

	return &block, nil
}

// 不在session包的service.go中被调用，而是被同为Store方法的GetSessionMeta调用
// 查询一个会话下的所有内容块
func (s *Store) ListContents(sessionID int64) ([]ContentBlock, error) {
	const sql = `SELECT id, session_id, title, content,
	selected, position, created_at, updated_at
	FROM blocks
	WHERE session_id = ?
	ORDER BY position, created_at`

	// 定义为方法而不是辅助函数书要是因为它要用Store结构体中的db
	rows, err := s.db.Query(sql, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contents := make([]ContentBlock, 0)
	for rows.Next() {
		var block ContentBlock
		err = rows.Scan(
			&block.ID, &block.SessionID, &block.Title, &block.Content,
			&block.Selected, &block.Position, &block.CreatedAt, &block.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		contents = append(contents, block)
	}

	return contents, rows.Err()
}

func (s *Store) AppendMessages(sessionID int64, messages []Message) error {
	sql := `INSERT INTO messages
	(session_id, role, content, created_at)
	VALUES (?, ?, ?, ?)`

	if len(messages) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	for _, message := range messages {
		_, err = tx.Exec(sql, sessionID, message.Role, message.Content, now)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	sql = `UPDATE sessions
	SET updated_at = ?
	WHERE id = ?`

	_, err = tx.Exec(sql, now, sessionID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// 不在session包的service.go中被调用，而是被同为Store方法的Create/Updata/DeleteContent调用
// 查询一个会话下的所有消息
func (s *Store) ListMessages(sessionID int64) ([]Message, error) {
	const sql = `SELECT id, role, content, created_at
	FROM messages
	WHERE session_id = ?
	ORDER BY id`

	// 定义为方法而不是辅助函数是因为它要用Store结构体中的db
	rows, err := s.db.Query(sql, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]Message, 0)
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ID, &message.Role, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, rows.Err()
}

// 不在session包的service.go中被调用，而是被同为Store方法的GetSession调用
// 更新当前会话的updated_at（sqlite不支持自动更新updated_at字段）
func (s *Store) touchSession(sessionID int64) error {
	const sql = `UPDATE sessions
	SET updated_at = ?
	WHERE id = ?`

	now := time.Now().Unix()

	_, err := s.db.Exec(sql, now, sessionID)

	return err
}
