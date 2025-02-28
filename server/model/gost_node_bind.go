package model

type GostNodeBind struct {
	Base
	NodeCode string     `gorm:"column:node_code;uniqueIndex:uidx_gost_node_bind;comment:节点编号"`
	Node     GostNode   `gorm:"foreignKey:NodeCode;references:Code"`
	UserCode string     `gorm:"column:user_code;uniqueIndex:uidx_gost_node_bind;comment:用户编号"`
	User     SystemUser `gorm:"foreignKey:UserCode;references:Code"`
}
