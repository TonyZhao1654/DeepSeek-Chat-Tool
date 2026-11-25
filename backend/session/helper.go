package session

import (
	"strings"
	"unicode/utf8"

	"deepseek-backend/store"
)

// 内容块标题为空时使用默认标题
func normalizeTitle(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return defaultContentTitle
	}

	return title
}

// 按字符预算保留最近的消息，至少保留最新一条；裁完跳过开头的assistant
func trimHistory(reservedChars int, messages []store.Message) []store.Message {
	budget := maxContextChars - reservedChars
	if budget <= 0 {
		return nil
	}

	total := 0
	start := len(messages)

	for i := len(messages) - 1; i >= 0; i-- {
		total += utf8.RuneCountInString(messages[i].Content)
		if total > budget && start < len(messages) {
			break
		}

		start = i
	}

	history := messages[start:]
	for len(history) > 0 && history[0].Role == "assistant" {
		history = history[1:]
	}

	return history
}

// 把勾选的内容块拼成 system prompt，没有选中项时返回空串
func buildSystemPrompt(contents []store.ContentBlock) string {
	var builder strings.Builder

	for _, content := range contents {
		if !content.Selected {
			continue
		}

		if builder.Len() == 0 {
			builder.WriteString("你是我的编程助手。下面是我当前项目的内容片段，回答时请结合这些内容，不要假设未提供的实现。\n")
		}

		builder.WriteString("\n## ")
		builder.WriteString(content.Title)

		fence := codeFence(content.Content)
		builder.WriteString("\n")
		builder.WriteString(fence)
		builder.WriteString("text")
		builder.WriteString("\n")
		builder.WriteString(content.Content)

		if !strings.HasSuffix(content.Content, "\n") {
			builder.WriteString("\n")
		}

		builder.WriteString(fence)
		builder.WriteString("\n")
	}

	return builder.String()
}

// 根据内容里最长的连续反引号生成安全围栏，避免内容提前闭合
func codeFence(content string) string {
	maxRun := 0
	currentRun := 0

	for _, r := range content {
		if r == '`' {
			currentRun++
			if currentRun > maxRun {
				maxRun = currentRun
			}
			continue
		}

		currentRun = 0
	}

	length := 3
	if maxRun >= length {
		length = maxRun + 1
	}

	return strings.Repeat("`", length)
}
