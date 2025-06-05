package event_node

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type DomainData struct {
	Domain     string
	Target     string
	Cert       string
	Key        string
	ForceHttps int
}

func (e *Event) WsDomain(baseUrl string, data DomainData) {
	ReqUpdateDomain(baseUrl, DomainReqData{
		Domain:     data.Domain,
		Target:     data.Target,
		Cert:       data.Cert,
		Key:        data.Key,
		ForceHttps: data.ForceHttps,
	})
}

type DomainReqData struct {
	Domain     string `json:"domain"`
	Target     string `json:"target"`
	Cert       string `json:"cert"`
	Key        string `json:"key"`
	ForceHttps int    `json:"forceHttps"`
}

func ReqUpdateDomain(baseUrl string, data DomainReqData) {
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
