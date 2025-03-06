package model

type GostNodeDomain struct {
	Base
	Prefix   string `gorm:"column:prefix;size:100;uniqueIndex:gost_node_domain_uidx;comment:域名前缀"`
	NodeCode string `gorm:"column:node_code;size:100;uniqueIndex:gost_node_domain_uidx;comment:节点编号"`
}
