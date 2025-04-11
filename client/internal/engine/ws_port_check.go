package engine

import (
	"gostc-sub/pkg/utils"
	"strconv"
)

type PortCheckResp struct {
	Code string `json:"code"` // 节点编号
	Use  bool   `json:"use"`  // 是否被占用
	Port string `json:"port"` // 端口
}

func (e *Event) WsPortCheck(data map[string]string) PortCheckResp {
	port, _ := strconv.Atoi(data["port"])
	var result = PortCheckResp{
		Code: data["code"],
		Use:  true,
		Port: data["port"],
	}
	if port != 0 {
		result = PortCheckResp{
			Code: data["code"],
			Use:  utils.IsUse(port),
			Port: data["port"],
		}
	}
	return result
}
