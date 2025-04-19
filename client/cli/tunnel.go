package main

import (
	"fmt"
	service2 "gostc-sub/internal/service"
	"strings"
	"sync"
	"time"
)

type Tunnel struct {
	VKey string
	Bind string
	Port string
}

var tunnelMap = &sync.Map{}

func runTunnels(mode string, vTunnels string, apiurl string) {
	var tunnelList []Tunnel
	tunnels := strings.Split(vTunnels, ",")
	for _, tunnel := range tunnels {
		kp := strings.Split(tunnel, ":")
		switch len(kp) {
		case 2:
			tunnelList = append(tunnelList, Tunnel{
				VKey: kp[0],
				Bind: "0.0.0.0",
				Port: kp[1],
			})
		case 3:
			tunnelList = append(tunnelList, Tunnel{
				VKey: kp[0],
				Bind: kp[1],
				Port: kp[2],
			})
		default:
			continue
		}
	}

	for _, tunnel := range tunnelList {
		switch mode {
		case "visit":
			svc := service2.NewTunnel(apiurl, tunnel.VKey, tunnel.Bind, tunnel.Port)
			if err := svc.Start(); err != nil {
				fmt.Println("私有隧道启动失败", tunnel.VKey, tunnel.Bind+":"+tunnel.Port, err)
				continue
			}
			fmt.Println("私有隧道启动成功", tunnel.VKey, tunnel.Bind+":"+tunnel.Port)
			tunnelMap.Store(tunnel.VKey, svc)
		case "p2p":
			svc := service2.NewP2P(apiurl, tunnel.VKey, tunnel.Bind, tunnel.Port)
			if err := svc.Start(); err != nil {
				fmt.Println("P2P隧道启动失败", tunnel.VKey, tunnel.Bind+":"+tunnel.Port, err)
				continue
			}
			fmt.Println("P2P隧道启动成功", tunnel.VKey, tunnel.Bind+":"+tunnel.Port)
			tunnelMap.Store(tunnel.VKey, svc)
		}
	}
	for {
		var isRunningCount = 0
		tunnelMap.Range(func(key, value any) bool {
			if value.(service2.Service).IsRunning() {
				isRunningCount++
			}
			return true
		})
		if isRunningCount == 0 {
			return
		}
		time.Sleep(time.Second * 3)
	}
}
