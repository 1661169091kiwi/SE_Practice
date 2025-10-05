package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MessageContent struct {
	Type     string           `json:"type"`
	Text     string           `json:"text,omitempty"`
	ImageURL *ImageURLContent `json:"image_url,omitempty"`
}

type ImageURLContent struct {
	URL string `json:"url"`
}

type GPTClient interface {
	CallGPT(ctx context.Context, content interface{}) (string, error)
}

// openAIClient 实现 GPTClient 接口
type openAIClient struct {
	apiKey string
	client *http.Client
}

func NewGPTClient(apiKey string) GPTClient {
	client := &http.Client{}
	return &openAIClient{
		apiKey: apiKey,
		client: client,
	}
}

func (c *openAIClient) CallGPT(ctx context.Context, content interface{}) (string, error) {
	// 1. 构造 GPT 请求体（逻辑不变）
	requestBody := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": content,
			},
		},
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("请求体序列化失败: %v", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.chatanywhere.tech/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}
	// 调试日志
	fmt.Println("GPT 原始响应 JSON:", string(bodyBytes))

	var responseData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return "", fmt.Errorf("解析响应失败: %v（原始响应：%s）", err, string(bodyBytes))
	}

	if errObj, ok := responseData["error"].(map[string]interface{}); ok {
		errCode, _ := errObj["code"].(string)
		return "", fmt.Errorf("GPT API 错误（%s）: %s", errCode, errObj["message"])
	}

	choices, ok := responseData["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("GPT 未返回有效结果（原始响应：%s）", string(bodyBytes))
	}
	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("GPT 结果格式错误（原始响应：%s）", string(bodyBytes))
	}
	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("GPT 消息格式错误（原始响应：%s）", string(bodyBytes))
	}
	contentVal, ok := message["content"].(interface{})
	if !ok {
		return "", fmt.Errorf("GPT 内容为空（原始响应：%s）", string(bodyBytes))
	}
	summary := fmt.Sprintf("%v", contentVal)

	return summary, nil
}
