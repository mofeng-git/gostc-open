package model

type Client struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Tls       int    `json:"tls"`
	AutoStart int    `json:"autoStart"`
}
