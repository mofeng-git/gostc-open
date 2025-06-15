package cache

import (
	"github.com/patrickmn/go-cache"
	"server/global"
)

const (
	admission_key = "admission:"
)

type AdmissionInfo struct {
	Code        string
	WhiteEnable int
	WhiteList   []string
	whiteMap    map[string]struct{}
	BlackEnable int
	Blacklist   []string
	blackMap    map[string]struct{}
}

func (admission *AdmissionInfo) ValidWhiteIp(ip string) bool {
	if admission.BlackEnable != 1 {
		return true
	}
	if _, ok := admission.whiteMap[ip]; ok {
		return true
	}
	return false
}

func (admission *AdmissionInfo) ValidBlackIp(ip string) bool {
	if admission.BlackEnable != 1 {
		return true
	}
	if _, ok := admission.blackMap[ip]; ok {
		return true
	}
	return false
}

func SetAdmissionInfo(info AdmissionInfo) {
	info.whiteMap = make(map[string]struct{})
	for _, item := range info.WhiteList {
		info.whiteMap[item] = struct{}{}
	}
	info.blackMap = make(map[string]struct{})
	for _, item := range info.Blacklist {
		info.blackMap[item] = struct{}{}
	}
	global.Cache.SetStruct(admission_key+info.Code, info, cache.NoExpiration)
}

func GetAdmissionInfo(code string) (result AdmissionInfo) {
	_ = global.Cache.GetStruct(admission_key+code, &result)
	return result
}

func DelAdmissionInfo(code string) {
	global.Cache.Del(admission_key + code)
}
