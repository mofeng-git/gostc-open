package main

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
	"gostc-sub/internal/common"
	service2 "gostc-sub/internal/service"
	service "gostc-sub/internal/service/visitor"
	"os"
	"strconv"
	"sync"
	"time"
)

var Cfg Config

type Config struct {
	Clients []ClientCfg `yaml:"clients"`
	Tunnels []TunnelCfg `yaml:"tunnels"`
	P2Ps    []P2PCfg    `yaml:"p2ps"`
}

type ClientCfg struct {
	Key     string `yaml:"key"`
	Remark  string `yaml:"remark"`
	Tls     bool   `yaml:"tls"`
	Address string `yaml:"address"`
}

type TunnelCfg struct {
	Key     string `yaml:"key"`
	Bind    string `yaml:"bind"`
	Port    int    `yaml:"port"`
	Remark  string `yaml:"remark"`
	Tls     bool   `yaml:"tls"`
	Address string `yaml:"address"`
}

type P2PCfg struct {
	Key     string `yaml:"key"`
	Bind    string `yaml:"bind"`
	Port    int    `yaml:"port"`
	Remark  string `yaml:"remark"`
	Tls     bool   `yaml:"tls"`
	Address string `yaml:"address"`
}

func loadConfig(file string) error {
	cfgFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("load config fail,%v", err)
	}
	if err := yaml.Unmarshal(cfgFile, &Cfg); err != nil {
		return fmt.Errorf("unmarshal fail,%v", err)
	}
	return nil
}

func configExample() {
	cfg := Config{
		Clients: []ClientCfg{
			{Key: uuid.NewString(), Remark: "客户端1", Tls: true, Address: "gost.sian.one"},
		},
		Tunnels: []TunnelCfg{
			{Key: uuid.NewString(), Bind: "127.0.0.1", Port: 8080, Remark: "私有隧道1", Tls: true, Address: "gost.sian.one"},
		},
		P2Ps: []P2PCfg{
			{Key: uuid.NewString(), Bind: "0.0.0.0", Port: 8081, Remark: "P2P隧道1", Tls: true, Address: "gost.sian.one"},
		},
	}
	marshal, _ := yaml.Marshal(cfg)
	fmt.Println(string(marshal))
}

func startForConfig() {
	wg := &sync.WaitGroup{}
	for _, client := range Cfg.Clients {
		svc := service2.NewClient(common.GenerateWsUrl(client.Tls, client.Address), common.GenerateHttpUrl(client.Tls, client.Address), client.Key)
		if err := svc.Start(); err != nil {
			fmt.Println(client.Key, client.Remark, "启动失败", err)
		}
		fmt.Println(client.Key, client.Remark, "客户端启动成功")
		wg.Add(1)
		go loop(wg, svc)
	}
	for _, tunnel := range Cfg.Tunnels {
		svc := service.NewTunnel(common.GenerateHttpUrl(tunnel.Tls, tunnel.Address), tunnel.Key, tunnel.Bind, tunnel.Port)
		if err := svc.Start(); err != nil {
			fmt.Println(tunnel.Key, tunnel.Remark, "启动失败", err)
		}
		fmt.Println(tunnel.Key, tunnel.Remark, "私有隧道启动成功", tunnel.Bind+":"+strconv.Itoa(tunnel.Port))
		wg.Add(1)
		go loop(wg, svc)
	}
	for _, p2p := range Cfg.P2Ps {
		svc := service.NewP2P(common.GenerateHttpUrl(p2p.Tls, p2p.Address), p2p.Key, p2p.Bind, p2p.Port)
		if err := svc.Start(); err != nil {
			fmt.Println(p2p.Key, p2p.Remark, "启动失败", err)
		}
		fmt.Println(p2p.Key, p2p.Remark, "P2P隧道启动成功", p2p.Bind+":"+strconv.Itoa(p2p.Port))
		wg.Add(1)
		go loop(wg, svc)
	}
	wg.Wait()
}

func loop(wg *sync.WaitGroup, svc service2.Service) {
	defer wg.Done()
	for svc.IsRunning() {
		time.Sleep(time.Second)
	}
}
