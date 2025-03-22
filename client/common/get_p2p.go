package common

import (
	"encoding/json"
	v1 "gostc-sub/p2p/pkg/config/v1"
	"io"
	"net/http"
)

type P2PCfg struct {
	Common  v1.ClientCommonConfig
	XTCPCfg v1.XTCPVisitorConfig
	STCPCfg v1.STCPVisitorConfig
}

func GetP2PConfig(url string) (result P2PCfg, err error) {
	response, err := http.Get(url)
	if err != nil {
		return P2PCfg{}, err
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return P2PCfg{}, err
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return P2PCfg{}, err
	}
	return result, nil
}
