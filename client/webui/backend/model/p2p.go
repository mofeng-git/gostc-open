package model

type P2P struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Port      string `json:"port"`
	Address   string `json:"address"`
	Tls       int    `json:"tls"`
	AutoStart int    `json:"autoStart"`
}
