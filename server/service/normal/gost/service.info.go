package service

import (
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/service/common/cache"
)

type InfoResp struct {
	Client      int64 `json:"client"`
	Host        int64 `json:"host"`
	Forward     int64 `json:"forward"`
	Tunnel      int64 `json:"tunnel"`
	InputBytes  int64 `json:"inputBytes"`
	OutputBytes int64 `json:"outputBytes"`
}

func (service *service) Info(claims jwt.Claims) (result InfoResp) {
	db, _, _ := repository.Get("")
	db.Model(&model.GostClient{}).Where("user_code = ?", claims.Code).Count(&result.Client)
	db.Model(&model.GostClientHost{}).Where("user_code = ?", claims.Code).Count(&result.Host)
	db.Model(&model.GostClientForward{}).Where("user_code = ?", claims.Code).Count(&result.Forward)
	db.Model(&model.GostClientTunnel{}).Where("user_code = ?", claims.Code).Count(&result.Tunnel)
	obsInfo := cache.GetUserObsDateRange(cache.MONTH_DATEONLY_LIST, claims.Code)
	result.InputBytes = obsInfo.InputBytes
	result.OutputBytes = obsInfo.OutputBytes
	return result
}
