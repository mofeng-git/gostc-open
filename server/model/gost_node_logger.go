package model

type GostNodeLogger struct {
	Id        int    `gorm:"primaryKey"`
	Level     string `gorm:"column:level;index"`
	NodeCode  string `gorm:"column:node_code;index"`
	Content   string `gorm:"column:content"`
	CreatedAt int64  `gorm:"column:created_at;index"`
}
