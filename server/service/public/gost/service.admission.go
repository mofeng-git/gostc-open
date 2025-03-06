package service

type AdmissionReq struct {
	Addr string `json:"addr"`
}

type AdmissionResp struct {
	Ok bool `json:"ok"`
}

func (service *service) Admission(req AdmissionReq) AdmissionResp {
	return AdmissionResp{Ok: false}
}
