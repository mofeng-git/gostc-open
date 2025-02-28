package model

import "encoding/json"

type GostClientAdmission struct {
	WhiteEnable int    `gorm:"column:white_enable;size:1;default:2;comment:是否启用准入控制"`
	BlackEnable int    `gorm:"column:black_enable;size:1;default:2;comment:是否启用准入控制"`
	WhiteList   string `gorm:"column:white_list;comment:白名单"`
	BlackList   string `gorm:"column:black_list;comment:黑名单"`
}

func (forward *GostClientAdmission) GetWhiteList() (result []string) {
	_ = json.Unmarshal([]byte(forward.WhiteList), &result)
	return result
}

func (forward *GostClientAdmission) SetWhiteList(data []string) {
	marshal, _ := json.Marshal(data)
	forward.WhiteList = string(marshal)
}

func (forward *GostClientAdmission) GetBlackList() (result []string) {
	_ = json.Unmarshal([]byte(forward.BlackList), &result)
	return result
}

func (forward *GostClientAdmission) SetBlackList(data []string) {
	marshal, _ := json.Marshal(data)
	forward.BlackList = string(marshal)
}
