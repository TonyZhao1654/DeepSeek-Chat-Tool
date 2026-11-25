package session

import (
	"net/http"

	"deepseek-backend/apperr"
	"deepseek-backend/store"
)

var httpErrorMap = apperr.HTTPErrorMap{
	store.ErrSessionNotFound: {Code: "SESSION_NOT_FOUND", Message: "会话不存在", Status: http.StatusNotFound},
	store.ErrContentNotFound: {Code: "CONTENT_NOT_FOUND", Message: "内容块不存在", Status: http.StatusNotFound},
}
