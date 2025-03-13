package service

import (
	"server/repository"
	"server/service/common/cache"
	"time"
)

type CountResp struct {
	Client       int64 `json:"client"`
	ClientOnline int64 `json:"clientOnline"`

	Host    int64 `json:"host"`
	Forward int64 `json:"forward"`
	Tunnel  int64 `json:"tunnel"`

	Node       int64 `json:"node"`
	NodeOnline int64 `json:"nodeOnline"`

	User        int64 `json:"user"`
	InputBytes  int64 `json:"inputBytes"`
	OutputBytes int64 `json:"outputBytes"`
}

func (service *service) Count() (result CountResp) {
	var dateOnly = time.Now().Format(time.DateOnly)
	db, _, _ := repository.Get("")
	var nodeCodes []string
	_ = db.GostNode.Pluck(db.GostNode.Code, &nodeCodes)
	result.Node = int64(len(nodeCodes))
	for _, nodeCode := range nodeCodes {
		obs := cache.GetNodeObs(dateOnly, nodeCode)
		result.InputBytes += obs.InputBytes
		result.OutputBytes += obs.OutputBytes
		if cache.GetNodeOnline(nodeCode) {
			result.NodeOnline++
		}
	}

	var clientCodes []string
	_ = db.GostClient.Pluck(db.GostClient.Code, &clientCodes)
	result.Client = int64(len(clientCodes))
	for _, clientCode := range clientCodes {
		if cache.GetClientOnline(clientCode) {
			result.ClientOnline++
		}
	}

	result.Host, _ = db.GostClientHost.Count()
	result.Forward, _ = db.GostClientForward.Count()
	result.Tunnel, _ = db.GostClientTunnel.Count()

	result.User, _ = db.SystemUser.Count()
	return result
}
