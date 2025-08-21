package model

type FrpClientCfg struct {
	Base
	Name        string     `gorm:"column:name;index;comment:名称"`
	Enable      int        `gorm:"column:enable;size:1;default:1;comment:启用状态"`
	Address     string     `gorm:"column:address;comment:访问地址(仅提示作用)"`
	Platform    string     `gorm:"column:platform;comment:FRP平台(仅提示作用)"`
	Content     string     `gorm:"column:content;comment:配置内容"`
	ContentType string     `gorm:"column:content_type;comment:配置类型"`
	ClientCode  string     `gorm:"column:client_code;index;comment:客户端编号"`
	Client      GostClient `gorm:"foreignKey:ClientCode;references:Code"`
	UserCode    string     `gorm:"column:user_code;index;comment:用户编号"`
	User        SystemUser `gorm:"foreignKey:UserCode;references:Code"`
}
