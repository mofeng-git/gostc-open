package program

import (
	"github.com/kardianos/service"
	"server/bootstrap"
	"time"
)

var SvcCfg = &service.Config{
	Name:        "gostc-admin",
	DisplayName: "GOSTC-ADMIN",
	Description: "基于GOST开发的内网穿透 后台管理",
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
	bootstrap.InitJwt()
	bootstrap.InitPersistence()
	bootstrap.InitMemory()
	bootstrap.InitTodo()
	bootstrap.InitTask()
	bootstrap.InitRouter()
	bootstrap.InitServer()
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
