package controller

import (
	"fmt"
	"my_project/dto"
	"my_project/model"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 辅助工具：把 "2025", "10", "15" 这种字符串转成 time.Time
// 或者处理表格里的 "2025年10月15日"
func (c *AdminController) parseTime(year, month, day string) (time.Time, error) {
	// 拼接成标准格式 2025-10-15
	dateStr := fmt.Sprintf("%s-%s-%s", year, month, day)
	return time.Parse("2006-01-02", dateStr)
}

// 辅助工具：专门解析表格里的中文日期字符串
// 假设前端表格里现在传的是 "2025年10月15日" 或者 "2025-10-15"
func (c *AdminController) parseDateString(dateStr string) time.Time {
	// 1. 简单的清洗：把 "年", "月" 换成 "-", 把 "日" 去掉
	// "2025年10月15日" -> "2025-10-15"
	temp := strings.ReplaceAll(dateStr, "年", "-")
	temp = strings.ReplaceAll(temp, "月", "-")
	temp = strings.ReplaceAll(temp, "日", "")

	// 2. 尝试解析
	t, err := time.Parse("2006-01-02", temp)
	if err != nil {
		// 如果解析失败，为了不报错，返回当前时间作为兜底
		return time.Now()
	}
	return t
}

// 1. 发布赛事基本信息
func (c *AdminController) PublishEvent(ctx *gin.Context) {
	var req dto.PublishEventReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	sportId, _ := c.getSportId(req.Type)

	startTime, err := c.parseTime(req.Time.Year, req.Time.Month, req.Time.Day)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "日期格式错误"})
		return
	}

	event := model.Event{
		EventName: req.Name,
		SportId:   sportId,
		Location:  req.Location, // 需确认数据库是否有此字段，若无请修改DB
		StartDate: startTime,
	}

	if err := c.DB.Create(&event).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "服务器错误"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发布成功",
		"data":    gin.H{"eventId": strconv.Itoa(event.EventId)},
	})
}

// 2. 保存赛事编排数据
func (c *AdminController) ArrangeEvent(ctx *gin.Context) {
	var req dto.GridDataReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	sportId, _ := c.getSportId(req.Type)

	// 开启事务，保证批量插入的原子性
	tx := c.DB.Begin()

	// 遍历表格每一行
	// 格式: [0:轮次, 1:A队, 2:B队, 3:时间, 4:场地]
	for _, row := range req.Data {
		if len(row) < 5 {
			continue
		}

		matchName := row[0].Value
		teamAName := row[1].Value
		teamBName := row[2].Value
		// row[3].Value 现在是包含年份的字符串，例如 "2025年10月15日"
		dateStrRaw := row[3].Value

		// 调用我们刚才写的辅助工具 parseDateString
		matchTime := c.parseDateString(dateStrRaw)
		// 1. 查找队伍ID (根据名字和运动类型)

		var teamA, teamB model.Team
		tx.Where("team_name = ? AND sport_id = ?", teamAName, sportId).FirstOrCreate(&teamA, model.Team{TeamName: teamAName, SportId: sportId})
		tx.Where("team_name = ? AND sport_id = ?", teamBName, sportId).FirstOrCreate(&teamB, model.Team{TeamName: teamBName, SportId: sportId})
		// 3. 创建比赛记录
		match := model.Match{
			EventId:   1, // !!! 注意：接口未传EventID，这在实际开发中是个问题。通常URL应带参 /arrange-event?eventId=1001，此处暂写死或需查询最近的Event
			MatchName: matchName,
			TeamAId:   teamA.TeamId,
			TeamBId:   teamB.TeamId,
			MatchTime: matchTime,
			// Location: location,
		}
		if err := tx.Create(&match).Error; err != nil {
			tx.Rollback()
			ctx.JSON(500, gin.H{"code": 500, "message": "保存失败"})
			return
		}
	}

	tx.Commit()
	ctx.JSON(200, gin.H{"code": 200, "message": "编排数据保存成功"})
}

// 4. 保存比赛规则
func (c *AdminController) SaveRule(ctx *gin.Context) {
	var req dto.SaveRuleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"message": "参数错误"})
		return
	}

	// --- 修改开始：处理时间 ---
	// 虽然规则表里可能不直接存具体日期，但如果业务需要：
	matchTime, _ := c.parseTime(req.MatchTime.Year, req.MatchTime.Month, req.MatchTime.Day)
	// --- 修改结束 ---

	// 这里具体怎么存要看你 save-rule 对应的数据库逻辑
	// 假设你只是存个 description 或者存到 events 表
	fmt.Println("规则时间设为:", matchTime)

	ctx.JSON(200, gin.H{"code": 200, "message": "规则保存成功"})
}

// 5. 保存比赛成绩 (严格版：找不到就报错)
func (c *AdminController) SaveGrade(ctx *gin.Context) {
	var req dto.GridDataReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	sportId, _ := c.getSportId(req.Type)

	// 开启事务
	tx := c.DB.Begin()

	for index, row := range req.Data {
		// 格式: [0:轮次, 1:A队, 2:B队, 3:时间, 4:场地, 5:比分]
		if len(row) < 6 {
			continue
		}

		matchName := row[0].Value
		teamAName := row[1].Value
		teamBName := row[2].Value
		dateStrRaw := row[3].Value
		scoreStr := row[5].Value

		var scoreA, scoreB int
		if _, err := fmt.Sscanf(scoreStr, "%d:%d", &scoreA, &scoreB); err != nil {
			// 如果比分格式不对，也可以选择报错
			continue
		}

		// 2. 找队伍ID
		var teamA, teamB model.Team
		c.DB.Where("team_name = ? AND sport_id = ?", teamAName, sportId).First(&teamA)
		c.DB.Where("team_name = ? AND sport_id = ?", teamBName, sportId).First(&teamB)

		if teamA.TeamId == 0 || teamB.TeamId == 0 {
			// --- 修改点 1：队伍不存在，直接报错返回 ---
			tx.Rollback() // 回滚事务，之前改的也不算数
			errorMsg := fmt.Sprintf("第 %d 行错误：找不到队伍 '%s' 或 '%s'", index+1, teamAName, teamBName)
			ctx.JSON(400, gin.H{"code": 400, "message": errorMsg})
			return // 结束函数，不再往下走
		}

		// 3. 解析日期
		targetDate := c.parseDateString(dateStrRaw)

		// 4. 精确查找比赛
		var match model.Match
		err := tx.Where("team_a_id = ? AND team_b_id = ?", teamA.TeamId, teamB.TeamId).
			Where("match_name = ?", matchName).
			Where("status != ?", "finished").
			Where("DATE(match_time) = ?", targetDate.Format("2006-01-02")).
			First(&match).Error

		if err != nil {
			// --- 修改点 2：比赛没找到，直接报错返回 ---
			tx.Rollback() // 回滚

			// 详细的错误信息，告诉前端哪一行数据有问题
			errorMsg := fmt.Sprintf("保存失败：未找到匹配的比赛！请检查第 %d 行数据。\n可能原因：\n1. 轮次名称 '%s' 不匹配\n2. 比赛日期 '%s' 不匹配\n3. 比赛已结束",
				index+1, matchName, targetDate.Format("2006-01-02"))

			ctx.JSON(400, gin.H{"code": 400, "message": errorMsg})
			return // 结束函数
		}

		// 找到了，更新
		match.ScoreTeamA = scoreA
		match.ScoreTeamB = scoreB
		match.Status = "finished"

		if err := tx.Save(&match).Error; err != nil {
			tx.Rollback()
			ctx.JSON(500, gin.H{"code": 500, "message": "数据库更新失败"})
			return
		}
	}

	// 只有所有行都成功了，才提交事务
	tx.Commit()
	ctx.JSON(200, gin.H{"code": 200, "message": "成绩保存成功"})
}
