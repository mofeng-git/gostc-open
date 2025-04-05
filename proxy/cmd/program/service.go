package program

import (
	"github.com/kardianos/service"
	"proxy/bootstrap"
	"time"
)

var SvcCfg = &service.Config{
	Name:        "gostc-proxy",
	DisplayName: "GOSTC-PROXY",
	Description: "GOSTC的代理网关，用于扩展自定义域名功能",
	Option:      make(service.KeyValue),
}

var Program = &program{
	stopChan: make(chan struct{}),
}

type program struct {
	stopChan chan struct{}
}

func (p *program) run() {
	bootstrap.InitLogger()
	bootstrap.InitConfig()
	bootstrap.InitServer()
	bootstrap.InitApi()

	<-p.stopChan
	bootstrap.Release()
}

func (p *program) Run() error {
	p.run()
	return nil
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) Stop(s service.Service) error {
	p.stopChan <- struct{}{}
	time.Sleep(time.Second)
	return nil
}
