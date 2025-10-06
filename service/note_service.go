package service

import (
	"context"
	"fmt"
	"multimodal-notes/util"
)

// NoteService 处理笔记相关的业务逻辑
type NoteService interface {
	SummarizeNote(textContent, imageBase64 string) (string, error)
}

// noteService 实现 NoteService 接口
type noteService struct {
	gptClient util.GPTClient
}

// NewNoteService 创建 NoteService 实例
func NewNoteService(gptClient util.GPTClient) NoteService {
	return &noteService{gptClient: gptClient}
}

func (ns *noteService) SummarizeNote(textContent, imageBase64 string) (string, error) {
	var gptContent interface{}
	if imageBase64 != "" {
		gptContent = []map[string]interface{}{
			{
				"type": "text",
				"text": "请结合以下文字和图片内容，进行总结，形成一份结构完整、条理清晰的笔记整理结果。同时根据全文类容生成标签，并且以#tag 的形式附在全文开头。比如#机器学习#决策树。然后才是整理内容。：\n" + textContent,
			},
			{
				"type": "image_url",
				"image_url": map[string]string{
					"url": imageBase64,
				},
			},
		}
	} else {
		gptContent = "请总结以下内容：\n" + textContent
	}

	fmt.Println("内容：", gptContent)

	summary, err := ns.gptClient.CallGPT(context.Background(), gptContent)
	if err != nil {
		fmt.Println("调用失败：", err)
		return "", err
	}
	return summary, nil
}
