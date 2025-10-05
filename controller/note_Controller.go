package controller

import (
	"encoding/base64"
	"io"
	"multimodal-notes/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 处理笔记相关的 HTTP 请求
type NoteController struct {
	noteService service.NoteService
}

// 创建 NoteController 实例
func NewNoteController(noteService service.NoteService) *NoteController {
	return &NoteController{noteService: noteService}
}

// 处理用户上传笔记
func (nc *NoteController) UploadNote(c *gin.Context) {
	// 1. 处理文字内容
	textContent := c.PostForm("content")
	if textContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文字内容不能为空"})
		return
	}

	var imageBase64 string
	file, err := c.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片上传失败"})
		return
	}

	if file != nil {
		fileStream, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取图片"})
			return
		}
		defer fileStream.Close()

		fileBytes, err := io.ReadAll(fileStream)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "图片读取失败"})
			return
		}
		imageFormat := "jpeg"
		if file.Filename[len(file.Filename)-4:] == ".png" {
			imageFormat = "png"
		}
		imageBase64 = "data:image/" + imageFormat + ";base64," + base64.StdEncoding.EncodeToString(fileBytes)
	}

	summary, err := nc.noteService.SummarizeNote(textContent, imageBase64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "总结笔记失败"})
		return
	}

	//  返回结果给前端
	c.JSON(http.StatusOK, gin.H{"summary": summary})
}
