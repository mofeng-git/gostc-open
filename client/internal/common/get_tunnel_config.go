package common

import (
	"encoding/json"
	"github.com/go-gost/x/config"
	v1 "gostc-sub/pkg/p2p/pkg/config/v1"
	"io"
	"net/http"
	"time"
)

type VisitCfg struct {
	SvcList  []config.ServiceConfig
	Chain    config.ChainConfig
	Limiter  config.LimiterConfig
	RLimiter config.LimiterConfig
	CLimiter config.LimiterConfig
}

func exec(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func GetVisitTunnelConfig(url string) (result VisitCfg, err error) {
	var bytes []byte
	retry := 0
	for {
		retry++
		bytes, err = exec(url)
		if err != nil {
			if retry >= 3 {
				return
			}
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return VisitCfg{}, err
	}
	return result, nil
}

type P2PCfg struct {
	Common         v1.ClientCommonConfig
	XTCPCfg        v1.XTCPVisitorConfig
	STCPCfg        v1.STCPVisitorConfig
	DisableForward int
}

func GetP2PTunnelConfig(url string) (result P2PCfg, err error) {
	var bytes []byte
	for retry := 0; retry <= 10; retry++ {
		bytes, err = exec(url)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return P2PCfg{}, err
	}
	return result, nil
}
