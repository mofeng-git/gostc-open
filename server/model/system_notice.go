package model

type SystemNotice struct {
	Base
	Title      string `gorm:"column:title;comment:标题"`
	Content    string `gorm:"column:content;comment:内容"`
	Hidden     int    `gorm:"column:hidden;size:2;comment:是否隐藏"`
	IndexValue int    `gorm:"column:index_value;index;default:1000;comment:排序，升序"`
}
