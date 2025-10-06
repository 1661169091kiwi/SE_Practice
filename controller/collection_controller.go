package controller

import (
	"multimodal-notes/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CollectionController struct {
	colService service.CollectionService
}

func NewCollectionController(colService service.CollectionService) *CollectionController {
	return &CollectionController{
		colService: colService,
	}
}

// 前端传 userID 和 noteID
type AddCollectionRequest struct {
	UserID uint `json:"userId" binding:"required"`
	NoteID uint `json:"noteId" binding:"required"`
}

func (c *CollectionController) AddCollection(ctx *gin.Context) {
	var req AddCollectionRequest
	//userID和noteID一定传
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	err := c.colService.AddCollection(req.UserID, req.NoteID)
	if err != nil {
		if err.Error() == "笔记不存在" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "添加收藏失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "收藏成功",
	})
}

type RemoveCollectionRequest struct {
	UserID uint `json:"userId" binding:"required"`
	NoteID uint `json:"noteId" binding:"required"`
}

func (c *CollectionController) RemoveCollection(ctx *gin.Context) {
	var req RemoveCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	err := c.colService.RemoveCollection(req.UserID, req.NoteID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "取消收藏成功",
	})
}

func (c *CollectionController) GetUserCollections(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户ID格式错误"})
		return
	}

	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	page := 1
	pageSize := 10
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "页码需大于0"})
			return
		}
	}
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 || pageSize > 50 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "每页数量需在1-50之间"})
			return
		}
	}

	collections, total, err := c.colService.GetUserCollections(uint(userID), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询收藏列表失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"collections": collections, // 收藏列表（含笔记详情）
			"pagination": gin.H{
				"total":     total,                                           // 总收藏数
				"page":      page,                                            // 当前页码
				"pageSize":  pageSize,                                        // 每页数量
				"totalPage": (total + int64(pageSize) - 1) / int64(pageSize), // 总页数
			},
		},
	})
}

func (c *CollectionController) CheckCollection(ctx *gin.Context) {
	// 从URL参数获取 userID 和 noteID
	userIDStr := ctx.Query("user_id")
	noteIDStr := ctx.Query("note_id")
	userID, err1 := strconv.Atoi(userIDStr)
	noteID, err2 := strconv.Atoi(noteIDStr)
	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户ID或笔记ID格式错误"})
		return
	}

	isCollected, err := c.colService.IsCollected(uint(userID), uint(noteID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "检查收藏状态失败: " + err.Error()})
		return
	}

	var msg string
	if isCollected {
		msg = "当前用户已收藏该笔记"
	} else {
		msg = "当前用户未收藏该笔记"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":     true,
		"isCollected": isCollected, // true=已收藏，false=未收藏
		"msg":         msg,         // 使用 if-else 赋值的 msg
	})
}
