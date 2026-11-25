package deepseek

/*
	这里声明deepseek api的接口请求地址
*/

const (
	// 对话	同时支持流式&非流式
	chatCompletionURL = "https://api.deepseek.com/chat/completions"

	// 获取模型列表
	modelsURL = "https://api.deepseek.com/models"

	// 查询账号余额
	balanceURL = "https://api.deepseek.com/user/balance"
)
