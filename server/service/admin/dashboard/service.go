package service

type service struct {
}

var Service *service

type UserObsItem struct {
	Account     string `json:"account"`
	InputBytes  int64  `json:"inputBytes"`
	OutputBytes int64  `json:"outputBytes"`
}

type userObsSortable []UserObsItem

func (u userObsSortable) Len() int {
	return len(u)
}

func (u userObsSortable) Less(i, j int) bool {
	return u[i].InputBytes+u[i].OutputBytes > u[j].InputBytes+u[j].OutputBytes
}

func (u userObsSortable) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

type NodeObsItem struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Online      int    `json:"online"`
	InputBytes  int64  `json:"inputBytes"`
	OutputBytes int64  `json:"outputBytes"`
}

type nodeObsSortable []NodeObsItem

func (u nodeObsSortable) Len() int {
	return len(u)
}

func (u nodeObsSortable) Less(i, j int) bool {
	return u[i].InputBytes+u[i].OutputBytes > u[j].InputBytes+u[j].OutputBytes
}

func (u nodeObsSortable) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

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
