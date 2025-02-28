package service

import "fmt"

type AdmissionReq struct {
	Addr string `json:"addr"`
}

type AdmissionResp struct {
	Ok bool `json:"ok"`
}

func (service *service) Admission(req AdmissionReq) AdmissionResp {
	fmt.Println(req)
	return AdmissionResp{Ok: false}
}
