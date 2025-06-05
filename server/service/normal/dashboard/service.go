package service

type service struct {
}

var Service *service

type TunnelObsItem struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	InputBytes  int64  `json:"inputBytes"`
	OutputBytes int64  `json:"outputBytes"`
}

type tunnelObsSortable []TunnelObsItem

func (u tunnelObsSortable) Len() int {
	return len(u)
}

func (u tunnelObsSortable) Less(i, j int) bool {
	return u[i].InputBytes+u[i].OutputBytes > u[j].InputBytes+u[j].OutputBytes
}

func (u tunnelObsSortable) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

type ClientObsItem struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Online      int    `json:"online"`
	InputBytes  int64  `json:"inputBytes"`
	OutputBytes int64  `json:"outputBytes"`
}

type clientObsSortable []ClientObsItem

func (u clientObsSortable) Len() int {
	return len(u)
}

func (u clientObsSortable) Less(i, j int) bool {
	return u[i].InputBytes+u[i].OutputBytes > u[j].InputBytes+u[j].OutputBytes
}

func (u clientObsSortable) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
