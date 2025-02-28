package service

import (
	"server/service/common/cache"
)

type AutherReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

type AutherResp struct {
	Ok bool   `json:"ok"`
	Id string `json:"id"`
}

func (service *service) Auther(req AutherReq) AutherResp {
	code := cache.GetGostAuth(req.Username, req.Password)
	if code == "" {
		return AutherResp{Ok: false, Id: ""}
	}
	return AutherResp{Ok: true, Id: code}
}
