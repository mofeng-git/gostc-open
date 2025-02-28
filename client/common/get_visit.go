package common

import (
	"encoding/json"
	"github.com/go-gost/x/config"
	"io"
	"net/http"
)

type VisitCfg struct {
	SvcList  []config.ServiceConfig
	Chain    config.ChainConfig
	Limiter  config.LimiterConfig
	RLimiter config.LimiterConfig
	CLimiter config.LimiterConfig
}

func GetVisitConfig(url string) (result VisitCfg, err error) {
	response, err := http.Get(url)
	if err != nil {
		return VisitCfg{}, err
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return VisitCfg{}, err
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return VisitCfg{}, err
	}
	return result, nil
}
