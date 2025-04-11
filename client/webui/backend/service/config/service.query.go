package service

import "gostc-sub/internal/common"

type QueryResp struct {
	Version string `json:"version"`
}

func (*service) Query() (result QueryResp) {
	return QueryResp{Version: common.VERSION}
}
