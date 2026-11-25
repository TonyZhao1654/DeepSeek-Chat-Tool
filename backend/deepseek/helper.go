package deepseek

/*
	该文件只放用于handler.go的辅助函数
*/

import (
	"bufio"
	"deepseek-backend/store"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

// 原样转发上游 SSE，同时汇总成 ChatResponse；clientBroken 表示客户端已断开
func forwardStream(w io.Writer, stream io.Reader) (ChatResponse, bool) {
	scanner := bufio.NewScanner(stream)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var chatResponse ChatResponse
	var content strings.Builder
	var finishReason string
	var streamError string
	clientBroken := false

	for scanner.Scan() {
		line := scanner.Text()

		_, err := fmt.Fprintf(w, "%s\n", line)
		if err != nil {
			log.Printf("write stream failed: %v", err)
			clientBroken = true
			break
		}

		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			continue
		}

		var chunk APIStreamChunk

		err = json.Unmarshal([]byte(data), &chunk)
		if err != nil {
			continue
		}

		if chunk.Error != nil {
			streamError = chunk.Error.Message
			break
		}

		if chatResponse.ID == "" {
			chatResponse.ID = chunk.ID
			chatResponse.Created = chunk.Created
			chatResponse.Model = chunk.Model
		}

		if chunk.Usage.TotalTokens > 0 {
			chatResponse.Usage = chunk.Usage
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		content.WriteString(chunk.Choices[0].Delta.Content)
		if chunk.Choices[0].FinishReason != "" {
			finishReason = chunk.Choices[0].FinishReason
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		log.Printf("read stream failed: %v", scanErr)
	}

	if streamError != "" {
		log.Printf("deepseek stream error: %s", streamError)
	}

	chatResponse.Object = "chat.completion"
	chatResponse.Choices = []ChatChoice{
		{
			Index:        0,
			Message:      Message{Role: "assistant", Content: content.String()},
			FinishReason: finishReason,
		},
	}
	chatResponse.Incomplete = finishReason == "" || streamError != "" || scanErr != nil
	chatResponse.StreamError = streamError

	return chatResponse, clientBroken
}

// 把对话消息转换成会话消息
func toSessionMessages(messages []Message) []store.Message {
	convertedMessage := make([]store.Message, 0, len(messages))
	for _, message := range messages {
		convertedMessage = append(convertedMessage, store.Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	return convertedMessage
}
