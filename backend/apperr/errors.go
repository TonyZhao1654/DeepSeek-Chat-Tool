package apperr

import (
	"errors"
	"net/http"
)

var (
	ErrInternal         = errors.New("internal server error")
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrInvalidRequest   = errors.New("invalid request")
)

// httpErrorMap apperr 自己的通用错误映射
var httpErrorMap = HTTPErrorMap{
	// 对应HTTP状态码 500
	ErrInternal: {Code: "INTERNAL_ERROR", Message: "后端服务内部错误", Status: http.StatusInternalServerError},

	// 对应HTTP状态码 405
	ErrMethodNotAllowed: {Code: "METHOD_NOT_ALLOWED", Message: "请求方法不允许", Status: http.StatusMethodNotAllowed},

	// 对应HTTP状态码 400
	ErrInvalidRequest: {Code: "INVALID_REQUEST", Message: "请求参数不合法", Status: http.StatusBadRequest},
}
