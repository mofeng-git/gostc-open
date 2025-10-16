package model

type GostClientProxy struct {
	Base
	Name           string     `gorm:"column:name;index;comment:名称"`
	Protocol       string     `gorm:"column:protocol;index;comment:代理类型"`
	Port           string     `gorm:"column:port;comment:访问端口"`
	AuthUser       string     `gorm:"column:auth_user;comment:认证用户"`
	AuthPwd        string     `gorm:"column:auth_pwd;comment:认证密码"`
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
	PoolCount      int        `gorm:"column:pool_count;default:0;comment:复用连接数量"`
	GostClientConfig
}
