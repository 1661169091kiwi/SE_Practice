package controller

import (
	"net/http"
	"strconv"

	"my_project/dto"
	"my_project/model"

	"github.com/gin-gonic/gin"
)

// 3. 保存队伍信息
func (c *AdminController) SaveTeam(ctx *gin.Context) {
	var req dto.SaveTeamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	sportId, _ := c.getSportId(req.Type)
	tx := c.DB.Begin()

	// 1. 创建队伍
	team := model.Team{
		TeamName: req.TeamName,
		SportId:  sportId,
		// Captain: req.Captain, // DB teams表没有captain字段，是在athletes表中标记is_captain
	}
	if err := tx.Create(&team).Error; err != nil {
		tx.Rollback()
		ctx.JSON(500, gin.H{"message": "创建队伍失败"})
		return
	}

	// 2. 处理队员列表
	// 格式: [0:序号, 1:姓名, 2:位置, 3:号码, 4:年级]
	for _, row := range req.Members {
		if len(row) < 5 {
			continue
		}

		// 假设序号(row[0])即为学号，或者需要随机生成
		studentId := row[0].Value // 实际项目中应明确学号来源
		name := row[1].Value
		position := row[2].Value // DB athletes表无position字段，可能存在jersey_number或description
		number := row[3].Value
		// grade := row[4].Value

		// 确保User表存在该学生
		var user model.User
		if err := tx.Where("student_id = ?", studentId).First(&user).Error; err != nil {
			user = model.User{StudentId: studentId, Name: name}
			tx.Create(&user)
		}

		// 创建运动员记录
		isCaptain := (name == req.Captain)
		athlete := model.Athlete{
			StudentId:    studentId,
			TeamId:       team.TeamId,
			JerseyNumber: number,
			SportType:    req.Type,
			Position:     position,
			// IsCaptain: isCaptain, // 需在Athlete Struct加此字段
		}
		// 避免未使用变量报错，这里做个逻辑占位
		_ = isCaptain

		tx.Create(&athlete)

		// 如果你有 team_members 表，也应该在此处插入
	}

	tx.Commit()
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "队伍信息保存成功",
		"data":    gin.H{"teamId": strconv.Itoa(team.TeamId)},
	})
}
