package engine

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HttpsDomainData struct {
	Domain string
	Target string
	Cert   string
	Key    string
}

func (e *Event) WsHttpsDomain(baseUrl string, data HttpsDomainData) {
	//_ = HttpsServer.AddDomain(data.Domain, data.Target, data.Cert, data.Key)
	ReqUpdateDomain(baseUrl, ReqData{
		Domain: data.Domain,
		Target: data.Target,
		Cert:   data.Cert,
		Key:    data.Key,
	})
}

type ReqData struct {
	Domain string `json:"domain"`
	Target string `json:"target"`
	Cert   string `json:"cert"`
	Key    string `json:"key"`
}

func ReqUpdateDomain(baseUrl string, data ReqData) {
	marshal, _ := json.Marshal(data)
	request, err := http.NewRequest(http.MethodPost, baseUrl+"/domain", bytes.NewReader(marshal))
	if err != nil {
		return
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
}
