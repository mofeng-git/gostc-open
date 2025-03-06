package model

const (
	GOST_TUNNEL_TYPE_HOST    = 1 // 域名解析
	GOST_TUNNEL_TYPE_FORWARD = 2 // 端口转发
	GOST_TUNNEL_TYPE_TUNNEL  = 3 // 私有隧道
)

type GostAuth struct {
	Base
	TunnelType int    `gorm:"column:tunnel_type;index;default:1;comment:隧道类型"`
	TunnelCode string `gorm:"column:tunnel_code;size:100;uniqueIndex;default:'';comment:隧道编号"`
	User       string `gorm:"column:user;size:100;uniqueIndex:gost_auth_uidx;comment:连接用户"`
	Password   string `gorm:"column:password;size:100;uniqueIndex:gost_auth_uidx;comment:连接密码"`
}
