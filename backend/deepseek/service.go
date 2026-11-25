package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	APIKey       string
	HTTPClient   *http.Client
	StreamClient *http.Client // 用于流式请求，不设整体超时，靠请求context取消
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewService(apiKey string) *Service {
	deepseekService := Service{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		StreamClient: &http.Client{},
	}

	return &deepseekService
}

// 发起流式对话请求，返回上游SSE响应体，由调用方负责Close
func (s *Service) StreamChat(ctx context.Context, chatRequest ChatRequest) (io.ReadCloser, error) {
	// 填必要字段
	streamRequest := APIStreamRequest{
		Model:       chatRequest.Model,
		Messages:    chatRequest.Messages,
		Temperature: chatRequest.Temperature,
		MaxTokens:   chatRequest.MaxTokens,
		Stream:      true, // 强制流式输出
	}

	switch chatRequest.Model {
	case "deepseek-flash":
		streamRequest.Tools = []Tool{
			{Type: "deepseek_web_search"}, // 开启联网权限
		}
	case "deepseek-v4-pro":
		streamRequest.Tools = nil // deepseek-v4-pro不支持工具链
	default:
		return nil, fmt.Errorf("deepseek: stream chat: unknown model %s\n", chatRequest.Model)
	}

	// 将streamRequest序列化为JSON，准备作为请求体发送给DeepSeek服务端
	body, err := json.Marshal(streamRequest)
	if err != nil {
		return nil, ErrDeepSeekRequestFailed
	}

	// body作为请求体
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, chatCompletionURL, bytes.NewReader(body))
	if err != nil {
		return nil, ErrDeepSeekRequestFailed
	}

	// 加请求头
	request.Header.Add("Accept", "text/event-stream")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+s.APIKey)

	// 发出HTTP请求并取得HTTP响应
	apiResponse, err := s.StreamClient.Do(request)
	if err != nil {
		return nil, ErrDeepSeekRequestFailed
	}

	if apiResponse.StatusCode != http.StatusOK {
		errorBody, readErr := io.ReadAll(io.LimitReader(apiResponse.Body, 4096))
		apiResponse.Body.Close()

		if readErr != nil {
			return nil, statusError(apiResponse.StatusCode, nil)
		}

		return nil, statusError(apiResponse.StatusCode, errorBody)
	}

	return apiResponse.Body, nil
}

// 获取余额
func (s *Service) GetBalance() (*BalanceResponse, error) {
	// 获得响应体
	apiResponse, err := doRequest(s.HTTPClient, s.APIKey, http.MethodGet, balanceURL, nil)
	if err != nil {
		return nil, err
	}

	var apiBalanceResponse APIBalanceResponse
	err = json.Unmarshal(apiResponse, &apiBalanceResponse)
	if err != nil {
		return nil, ErrDeepSeekDecodeFailed
	}

	if len(apiBalanceResponse.BalanceInfos) == 0 {
		return nil, ErrDeepSeekBalanceUnavailable
	}

	balanceInfo := apiBalanceResponse.BalanceInfos[0]
	// BalanceInfos[0]下有四个值，但只用到了两个
	balanceResponse := BalanceResponse{
		Currency:     balanceInfo.Currency,
		TotalBalance: balanceInfo.TotalBalance,
	}

	return &balanceResponse, nil
}

// 获取当前账号可用的模型列表
func (s *Service) GetModels() ([]APIModel, error) {
	apiResponse, err := doRequest(s.HTTPClient, s.APIKey, http.MethodGet, modelsURL, nil)
	if err != nil {
		return nil, err
	}

	var modelsResponse APIModelsResponse

	err = json.Unmarshal(apiResponse, &modelsResponse)
	if err != nil {
		return nil, ErrDeepSeekResponseParse
	}

	if modelsResponse.Data == nil {
		return []APIModel{}, nil
	}

	return modelsResponse.Data, nil
}

// 统一封装对DeepSeek API的请求，返回原始响应体		body为nil表示不携带请求体
func doRequest(client *http.Client, apiKey, method, url string, body []byte) ([]byte, error) {
	var requestBody io.Reader

	if body != nil {
		requestBody = bytes.NewReader(body)
	}

	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, ErrDeepSeekRequestFailed
	}

	// 加请求头
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+apiKey)
	if body != nil {
		request.Header.Add("Content-Type", "application/json")
	}

	// 发送HTTP请求并取得HTTP响应
	apiResponse, err := client.Do(request)
	if err != nil {
		return nil, ErrDeepSeekRequestFailed
	}
	defer apiResponse.Body.Close()

	// 读取DeepSeek返回的响应体
	responseBody, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		return nil, ErrDeepSeekReadBodyFailed
	}

	if apiResponse.StatusCode == http.StatusOK {
		// 全部正常时，从这里出去
		return responseBody, nil
	}

	return nil, statusError(apiResponse.StatusCode, responseBody)
}

// 把HTTP状态码映射成本包的哨兵错误
func statusError(statusCode int, body []byte) error { // body是"响应体"
	var err error

	switch statusCode {
	case http.StatusUnauthorized:
		err = wrapError(ErrDeepSeekUnauthorized, body)
	case http.StatusBadRequest:
		err = wrapError(ErrDeepSeekBadRequest, body)
	case http.StatusPaymentRequired:
		err = wrapError(ErrDeepSeekInsufficientBalance, body)
	case http.StatusTooManyRequests:
		err = wrapError(ErrDeepSeekRateLimited, body)
	default:
		err = wrapError(ErrDeepSeekRequestFailed, body)
	}

	return err
}

// 对错误进行包装
func wrapError(err error, body []byte) error {
	message := strings.TrimSpace(string(body))
	if message == "" {
		return err
	}

	if len(message) > 512 {
		// 只保留前512个字符
		message = message[:512]
	}

	return fmt.Errorf("%w: %s", err, message)
}
