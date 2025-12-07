package controller

import (
	"my_project/model"

	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

// 辅助函数：根据中文名称获取运动ID
func (c *AdminController) getSportId(name string) (int, error) {
	var sport model.Sport
	err := c.DB.Where("sport_name = ?", name).First(&sport).Error
	if err != nil {
		// 如果不存在，自动创建（可选）
		sport = model.Sport{SportName: name}
		c.DB.Create(&sport)
		return sport.SportId, nil
	}
	return sport.SportId, nil
}
