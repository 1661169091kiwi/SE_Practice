package model

type Team struct {
	TeamId    int    `gorm:"primaryKey;autoIncrement"`
	TeamName  string `gorm:"not null"`
	SportId   int    `gorm:"not null"`
	CreatedBy string `gorm:"default:'admin'"` // 简化处理，默认填admin
}

type Athlete struct {
	AthleteId    int    `gorm:"primaryKey;autoIncrement"`
	StudentId    string `gorm:"not null"`
	TeamId       int    `gorm:"not null"`
	JerseyNumber string
	SportType    string
	Position     string
}
