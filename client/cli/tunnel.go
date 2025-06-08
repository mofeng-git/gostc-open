package main

import (
	"fmt"
	service "gostc-sub/internal/service/visitor"
	"strconv"
	"strings"
	"sync"
)

type Tunnel struct {
	VKey string
	Bind string
	Port string
}

func runTunnels(mode string, vTunnels string, apiurl string) {
	var wg = &sync.WaitGroup{}
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
			port, err := strconv.Atoi(tunnel.Port)
			if err != nil {
				fmt.Println("私有隧道启动失败", tunnel.VKey, tunnel.Bind+":"+tunnel.Port, err)
				continue
			}
			svc := service.NewTunnel(apiurl, tunnel.VKey, tunnel.Bind, port)
			if err := svc.Start(); err != nil {
				fmt.Println("私有隧道启动失败", tunnel.VKey, tunnel.Bind+":"+tunnel.Port, err)
				continue
			}
			fmt.Println("私有隧道启动成功", tunnel.VKey, tunnel.Bind+":"+tunnel.Port)
			go loop(wg, svc)
		case "p2p":
			port, err := strconv.Atoi(tunnel.Port)
			if err != nil {
				fmt.Println("P2P隧道启动失败", tunnel.VKey, tunnel.Bind+":"+tunnel.Port, err)
				continue
			}
			svc := service.NewP2P(apiurl, tunnel.VKey, tunnel.Bind, port)
			if err := svc.Start(); err != nil {
				fmt.Println("P2P隧道启动失败", tunnel.VKey, tunnel.Bind+":"+tunnel.Port, err)
				continue
			}
			fmt.Println("P2P隧道启动成功", tunnel.VKey, tunnel.Bind+":"+tunnel.Port)
			wg.Add(1)
			go loop(wg, svc)
		}
	}
	wg.Wait()
}
