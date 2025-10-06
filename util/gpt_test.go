package util

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"testing"
)

// 把本地图片文件转成 base64 字符串
func imageToBase64(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(file), nil
}

// TestGPTClient_CallGPT_LocalImage 测试文本+本地图片输入
func TestGPTClient_CallGPT_LocalImage(t *testing.T) {
	apiKey := "sk-gLJpO9I10cMfjn0PYz80SSwELl84fmTyKjhYlMUwkyANTfpf"
	if apiKey == "" {
		t.Fatal("API Key 不能为空！")
	}

	gptClient := NewGPTClient(apiKey)
	if gptClient == nil {
		t.Fatal("创建 GPT 客户端失败")
	}

	imagePath := "C:\\Users\\34517\\Documents\\xwechat_files\\wxid_8c3i3fb6j5sr22_56e0\\msg\\file\\2025-09\\plot\\plot\\heatmap-label2.png" // 假设项目根目录下有这张图片

	// 4. 转成 base64
	imgBase64, err := imageToBase64(imagePath)
	if err != nil {
		t.Fatalf("读取图片失败: %v", err)
	}

	// 5. 构造多模态输入
	content := []MessageContent{
		{
			Type: "text",
			Text: "详细描述这张图片的内容",
		},
		{
			Type: "image_url",
			ImageURL: &ImageURLContent{
				URL: fmt.Sprintf("data:image/png;base64,%s", imgBase64),
			},
		},
	}

	fmt.Println("正在调用 GPT API，输入：文本 + 本地图片(base64)")

	// 6. 调用 API
	ctx := context.Background()
	summary, err := gptClient.CallGPT(ctx, content)
	if err != nil {
		t.Fatalf("GPT API 调用失败！错误：%v", err)
	}

	// 7. 输出结果
	fmt.Println("GPT API 调用成功！")
	fmt.Println("返回结果：", summary)
}
