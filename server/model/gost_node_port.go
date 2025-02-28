package model

type GostNodePort struct {
	Base
	Port     string `gorm:"column:port;uniqueIndex:gost_node_port_uidx;comment:节点端口"`
	NodeCode string `gorm:"column:node_code;uniqueIndex:gost_node_port_uidx;comment:节点编号"`
}
