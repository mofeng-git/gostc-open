package model

const (
	GOST_OBS_ORIGIN_KIND_NODE   = 1 // 节点
	GOST_OBS_ORIGIN_KIND_CLIENT = 2 // 客户端
	GOST_OBS_ORIGIN_KIND_TUNNEL = 3 // 隧道
)

type GostObs struct {
	Base
	OriginKind  int    `gorm:"column:origin_kind;index"`
	OriginCode  string `gorm:"column:origin_code;size:100;uniqueIndex:gost_obs_uidx"`
	EventDate   string `gorm:"column:event_date;size:100;uniqueIndex:gost_obs_uidx"`
	InputBytes  int64  `gorm:"column:input_bytes;index"`
	OutputBytes int64  `gorm:"column:output_bytes;index"`
}
