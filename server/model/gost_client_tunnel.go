package model

type GostClientTunnel struct {
	Base
	Name           string     `gorm:"column:name;index;comment:名称"`
	TargetIp       string     `gorm:"column:target_ip;index;comment:内网IP"`
	TargetPort     string     `gorm:"column:target_port;index;comment:内网端口"`
	VKey           string     `gorm:"column:v_key;comment:访问密钥"`
	NoDelay        int        `gorm:"column:no_delay;size:1;comment:无等待延迟"`
	NodeCode       string     `gorm:"column:node_code;index;comment:节点编号"`
	Node           GostNode   `gorm:"foreignKey:NodeCode;references:Code"`
	ClientCode     string     `gorm:"column:client_code;index;comment:客户端编号"`
	Client         GostClient `gorm:"foreignKey:ClientCode;references:Code"`
	UserCode       string     `gorm:"column:user_code;index;comment:用户编号"`
	User           SystemUser `gorm:"foreignKey:UserCode;references:Code"`
	Enable         int        `gorm:"column:enable;size:1;default:1;comment:启用状态"`
	Status         int        `gorm:"column:status;size:1;default:1;comment:状态"`
	UseEncryption  int        `gorm:"column:use_encryption;size:1;default:1;comment:加密"`
	UseCompression int        `gorm:"column:use_compression;size:1;default:1;comment:压缩"`
	GostClientConfig
}
