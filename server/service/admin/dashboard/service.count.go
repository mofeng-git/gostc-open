package service

import (
	"server/repository"
	cache2 "server/repository/cache"
	"time"
)

type CountResp struct {
	Client       int64 `json:"client"`
	ClientOnline int64 `json:"clientOnline"`

	Host    int64 `json:"host"`
	Forward int64 `json:"forward"`
	Tunnel  int64 `json:"tunnel"`
	Proxy   int64 `json:"proxy"`
	P2P     int64 `json:"p2p"`

	Node       int64 `json:"node"`
	NodeOnline int64 `json:"nodeOnline"`

	User        int64 `json:"user"`
	InputBytes  int64 `json:"inputBytes"`
	OutputBytes int64 `json:"outputBytes"`

	CheckInTotal int64 `json:"checkInTotal"`
}

func (service *service) Count() (result CountResp) {
	var dateOnly = time.Now().Format(time.DateOnly)
	db, _, _ := repository.Get("")
	var nodeCodes []string
	_ = db.GostNode.Pluck(db.GostNode.Code, &nodeCodes)
	result.Node = int64(len(nodeCodes))
	for _, nodeCode := range nodeCodes {
		obs := cache2.GetNodeObs(dateOnly, nodeCode)
		result.InputBytes += obs.InputBytes
		result.OutputBytes += obs.OutputBytes
		if cache2.GetNodeOnline(nodeCode) {
			result.NodeOnline++
		}
	}

	var clientCodes []string
	_ = db.GostClient.Pluck(db.GostClient.Code, &clientCodes)
	result.Client = int64(len(clientCodes))
	for _, clientCode := range clientCodes {
		if cache2.GetClientOnline(clientCode) {
			result.ClientOnline++
		}
	}

	result.Host, _ = db.GostClientHost.Count()
	result.Forward, _ = db.GostClientForward.Count()
	result.Tunnel, _ = db.GostClientTunnel.Count()
	result.Proxy, _ = db.GostClientProxy.Count()
	result.P2P, _ = db.GostClientP2P.Count()

	result.User, _ = db.SystemUser.Count()

	result.CheckInTotal, _ = db.SystemUserCheckin.Where(db.SystemUserCheckin.EventDate.Eq(time.Now().Format(time.DateOnly))).Count()
	return result
}
