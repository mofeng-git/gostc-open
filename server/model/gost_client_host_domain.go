package model

type GostClientHostDomain struct {
	Base
	Domain string `gorm:"column:domain;size:50;uniqueIndex;comment:自定义域名"`
}
