package model

type GostNodeDomain struct {
	Base
	Prefix   string `gorm:"column:prefix;uniqueIndex:gost_node_domain_uidx;comment:域名前缀"`
	NodeCode string `gorm:"column:node_code;uniqueIndex:gost_node_domain_uidx;comment:节点编号"`
}
