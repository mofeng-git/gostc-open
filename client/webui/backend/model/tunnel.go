package model

type Tunnel struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Bind      string `json:"bind"`
	Port      string `json:"port"`
	Address   string `json:"address"`
	Tls       int    `json:"tls"`
	AutoStart int    `json:"autoStart"`
}
