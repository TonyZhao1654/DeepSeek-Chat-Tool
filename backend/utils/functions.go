package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// fileName需要加后缀
func GetApiKey(basePath string, fileName string) (string, error) {
	apiKey, err := os.ReadFile(basePath + fileName)
	if err != nil {
		return "", fmt.Errorf("read fail: %w", err)
	}

	apiKeyStr := strings.TrimSpace(string(apiKey))
	apiKeyStr = strings.TrimSpace(strings.TrimPrefix(apiKeyStr, "\ufeff"))

	if apiKeyStr == "" {
		return "", fmt.Errorf("api empty")
	}

	return apiKeyStr, nil
}

// 用系统默认浏览器打开地址（Windows）
func OpenBrowser(url string) error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)

	return cmd.Start()
}

// 每百万token的单价（元）
type tokenPrice struct {
	Input        float64
	Output       float64
	OffPeakScale float64 // 夜间优惠的折扣系数
}

// 价格是会调整的，这里只是固定值估算
var tokenPrices = map[string]tokenPrice{
	"deepseek-flash": {
		Input:        2,
		Output:       8,
		OffPeakScale: 0.5,
	},
	"deepseek-v4-pro": {
		Input:        4,
		Output:       16,
		OffPeakScale: 0.25,
	},
}

// 判断是否处于夜间时段（00:30-08:30）
func isOffPeak(t time.Time) bool {
	local := t.Local() // 转为本地时区
	minutes := local.Hour()*60 + local.Minute()

	return minutes >= 30 && minutes < 8*60+30
}

// 按token数估算花费金额（CNY）
func EstimateCost(model string, inputTokens, outputTokens int) float64 {
	price, ok := tokenPrices[model]
	if !ok {
		price = tokenPrices["deepseek-flash"] // flash模型兜底
	}

	cost := float64(inputTokens)/1e6*price.Input + float64(outputTokens)/1e6*price.Output
	if isOffPeak(time.Now()) {
		cost *= price.OffPeakScale
	}

	return cost
}
