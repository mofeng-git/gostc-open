package service

import (
	v1 "server/pkg/p2p_cfg/v1"
	"server/repository"
)

type VisitCfgResp struct {
	Common  v1.ClientCommonConfig
	XTCPCfg v1.XTCPVisitorConfig
	STCPCfg v1.STCPVisitorConfig
}

func (service *service) VisitCfg(key string) (result VisitCfgResp) {
	db, _, _ := repository.Get("")
	p2p, _ := db.GostClientP2P.Preload(db.GostClientP2P.Node).Where(db.GostClientP2P.VKey.Eq(key)).First()
	if p2p == nil {
		return result
	}
	result.Common, _ = p2p.GenerateCommonCfg()
	result.STCPCfg, result.XTCPCfg = p2p.GenerateVisitorCfgs()
	return result
}
