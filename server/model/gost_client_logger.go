package model

type GostClientLogger struct {
	Id         int    `gorm:"primaryKey"`
	Level      string `gorm:"column:level;index"`
	ClientCode string `gorm:"column:client_code;index"`
	Content    string `gorm:"column:content"`
	CreatedAt  int64  `gorm:"column:created_at;index"`
}
