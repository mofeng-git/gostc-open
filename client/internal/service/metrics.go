package service

import (
	"github.com/SianHH/frp-package/server/metrics"
	"github.com/lesismal/arpc"
	"time"
)

type Metrics struct {
	client *arpc.Client
}

func InitMetrics(client *arpc.Client) error {
	metrics.Server = &Metrics{
		client: client,
	}
	return nil
}

func (m *Metrics) NewClient() {

}

func (m *Metrics) CloseClient() {

}

func (m *Metrics) NewProxy(name string, proxyType string) {
}

func (m *Metrics) CloseProxy(name string, proxyType string) {
}

func (m *Metrics) OpenConnection(name string, proxyType string) {
}

func (m *Metrics) CloseConnection(name string, proxyType string) {
}

type TrafficData struct {
	Name      string
	ProxyType string
	Total     int64
}

func (m *Metrics) AddTrafficIn(name string, proxyType string, trafficBytes int64) {
	_ = m.client.CallAsync("rpc/metrics/input", TrafficData{
		Name:      name,
		ProxyType: proxyType,
		Total:     trafficBytes,
	}, func(c *arpc.Context, err error) {
	}, time.Second*10)
}

func (m *Metrics) AddTrafficOut(name string, proxyType string, trafficBytes int64) {
	_ = m.client.CallAsync("rpc/metrics/output", TrafficData{
		Name:      name,
		ProxyType: proxyType,
		Total:     trafficBytes,
	}, func(c *arpc.Context, err error) {
	}, time.Second*10)
}
