package controller

import (
	"multimodal-notes/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SaveNoteRequest 请求参数结构
type SaveNoteRequest struct {
	OriginalText string   `json:"originalText" binding:"required"`
	Summary      string   `json:"summary" binding:"required"`
	ImageBase64  string   `json:"imageBase64"`
	UserID       uint     `json:"userId" binding:"required"` // 新增：发表笔记的用户ID
	Tags         []string `json:"tags"`                      // 新增标签列表

}

// SaveNoteController 专门处理保存笔记
type SaveNoteController struct {
	saveNoteService service.SaveNoteService
}

func NewSaveNoteController(saveNoteService service.SaveNoteService) *SaveNoteController {
	return &SaveNoteController{saveNoteService: saveNoteService}
}

// 保存笔记
func (c *SaveNoteController) SaveNote(ctx *gin.Context) {
	var req SaveNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	noteID, err := c.saveNoteService.SaveNote(req.UserID, req.OriginalText, req.Summary, req.ImageBase64, req.Tags)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存笔记失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"noteID":  noteID,
		"msg":     "笔记保存成功",
	})
}
