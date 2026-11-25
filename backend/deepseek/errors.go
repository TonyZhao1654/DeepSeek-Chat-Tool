package deepseek

import (
	"errors"
	"net/http"

	"deepseek-backend/apperr"
	"deepseek-backend/store"
)

var (
	ErrDeepSeekRequestFailed       = errors.New("deepseek api request failed")
	ErrDeepSeekUnauthorized        = errors.New("deepseek api unauthorized, check apikey")
	ErrDeepSeekBadRequest          = errors.New("deepseek api bad request")
	ErrDeepSeekInsufficientBalance = errors.New("deepseek api insufficient balance")
	ErrDeepSeekRateLimited         = errors.New("deepseek api rate limited")
	ErrDeepSeekResponseParse       = errors.New("parse deepseek api response failed")
	ErrDeepSeekBalanceUnavailable  = errors.New("deepseek balance unavailable")
	ErrDeepSeekReadBodyFailed      = errors.New("read deepseek api response failed")
	ErrDeepSeekDecodeFailed        = errors.New("decode deepseek api response failed")
)

var httpErrorMap = apperr.HTTPErrorMap{
	ErrDeepSeekRequestFailed:       {Code: "DEEPSEEK_REQUEST_FAILED", Message: "请求 DeepSeek 服务失败，请稍后重试", Status: http.StatusBadGateway},
	ErrDeepSeekUnauthorized:        {Code: "DEEPSEEK_UNAUTHORIZED", Message: "API Key 无效或已过期，请检查配置", Status: http.StatusUnauthorized},
	ErrDeepSeekBadRequest:          {Code: "DEEPSEEK_BAD_REQUEST", Message: "DeepSeek 拒绝了请求参数，请检查模型或消息内容", Status: http.StatusBadRequest},
	ErrDeepSeekInsufficientBalance: {Code: "DEEPSEEK_INSUFFICIENT_BALANCE", Message: "DeepSeek 账户余额不足，请先充值", Status: http.StatusPaymentRequired},
	ErrDeepSeekRateLimited:         {Code: "DEEPSEEK_RATE_LIMITED", Message: "请求过于频繁，请稍后重试", Status: http.StatusTooManyRequests},
	ErrDeepSeekResponseParse:       {Code: "DEEPSEEK_RESPONSE_PARSE_FAILED", Message: "解析 DeepSeek 响应失败", Status: http.StatusBadGateway},
	ErrDeepSeekBalanceUnavailable:  {Code: "DEEPSEEK_BALANCE_UNAVAILABLE", Message: "无法获取账户余额", Status: http.StatusServiceUnavailable},
	ErrDeepSeekReadBodyFailed:      {Code: "DEEPSEEK_READ_BODY_FAILED", Message: "读取 DeepSeek 响应失败", Status: http.StatusBadGateway},
	ErrDeepSeekDecodeFailed:        {Code: "DEEPSEEK_DECODE_FAILED", Message: "解析 DeepSeek 响应数据失败", Status: http.StatusBadGateway},

	// 带 session_id 对话时会话不存在
	store.ErrSessionNotFound: {Code: "SESSION_NOT_FOUND", Message: "会话不存在", Status: http.StatusNotFound},
}
