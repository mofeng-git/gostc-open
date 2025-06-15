package service

import (
	"server/pkg/jwt"
	"server/repository"
	"server/repository/cache"
)

type InfoResp struct {
	Client      int64 `json:"client"`
	Host        int64 `json:"host"`
	Forward     int64 `json:"forward"`
	Tunnel      int64 `json:"tunnel"`
	Proxy       int64 `json:"proxy"`
	P2P         int64 `json:"p2p"`
	InputBytes  int64 `json:"inputBytes"`
	OutputBytes int64 `json:"outputBytes"`
}

func (service *service) Info(claims jwt.Claims) (result InfoResp) {
	db, _, _ := repository.Get("")
	result.Client, _ = db.GostClient.Where(db.GostClient.UserCode.Eq(claims.Code)).Count()
	result.Host, _ = db.GostClientHost.Where(db.GostClientHost.UserCode.Eq(claims.Code)).Count()
	result.Forward, _ = db.GostClientForward.Where(db.GostClientForward.UserCode.Eq(claims.Code)).Count()
	result.Tunnel, _ = db.GostClientTunnel.Where(db.GostClientTunnel.UserCode.Eq(claims.Code)).Count()
	result.Proxy, _ = db.GostClientProxy.Where(db.GostClientProxy.UserCode.Eq(claims.Code)).Count()
	result.P2P, _ = db.GostClientP2P.Where(db.GostClientP2P.UserCode.Eq(claims.Code)).Count()
	obsInfo := cache.GetUserObsDateRange(cache.MONTH_DATEONLY_LIST, claims.Code)
	result.InputBytes = obsInfo.InputBytes
	result.OutputBytes = obsInfo.OutputBytes
	return result
}
