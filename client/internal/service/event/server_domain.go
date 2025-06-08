package event

import (
	"bytes"
	"encoding/json"
	"github.com/lesismal/arpc"
	"net/http"
)

func ServerDomainHandle(client *arpc.Client, proxyBaseUrl string) {
	client.Handler.Handle("server_domain_config", func(c *arpc.Context) {
		if proxyBaseUrl == "" {
			_ = c.Write("not support")
			return
		}
		var req ServerDomain
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		if err := domainRequestDo(proxyBaseUrl, domainRequestData{
			Domain:     req.Domain,
			Target:     req.Target,
			Cert:       req.Cert,
			Key:        req.Key,
			ForceHttps: req.ForceHttps,
		}); err != nil {
			_ = c.Write(err)
			return
		}
		_ = c.Write("success")
	})
}

type domainRequestData struct {
	Domain     string `json:"domain"`
	Target     string `json:"target"`
	Cert       string `json:"cert"`
	Key        string `json:"key"`
	ForceHttps int    `json:"forceHttps"`
}

func domainRequestDo(baseUrl string, data domainRequestData) error {
	marshal, _ := json.Marshal(data)
	request, err := http.NewRequest(http.MethodPost, baseUrl+"/domain", bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	_ = response.Body.Close()
	return nil
}
