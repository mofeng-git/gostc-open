package main

import (
	"fmt"
	"gostc-sub/common"
	"gostc-sub/p2p/frpc"
	v1 "gostc-sub/p2p/pkg/config/v1"
	"gostc-sub/p2p/registry"
	"gostc-sub/pkg/signal"
	"strconv"
	"strings"
)

type P2P struct {
	VKey string
	Port string
}

func runP2P(vTunnels string, apiurl string) {
	var tunnelList []P2P
	tunnels := strings.Split(vTunnels, ",")
	for _, tunnel := range tunnels {
		kp := strings.Split(tunnel, ":")
		if len(kp) != 2 {
			continue
		}
		tunnelList = append(tunnelList, P2P{
			VKey: kp[0],
			Port: kp[1],
		})
	}

	for _, tunnel := range tunnelList {
		data, err := common.GetP2PConfig(apiurl + "/api/v1/public/p2p/visit?key=" + tunnel.VKey)
		if err != nil {
			fmt.Println("获取P2P隧道配置失败", tunnel.VKey, tunnel.Port, err)
			continue
		}
		data.XTCPCfg.BindPort, _ = strconv.Atoi(tunnel.Port)

		registry.Del(tunnel.VKey)
		svc := frpc.NewService(data.Common, nil, []v1.VisitorConfigurer{
			&data.XTCPCfg,
			&data.STCPCfg,
		})
		if err := svc.Start(); err != nil {
			fmt.Println("运行P2P隧道配置失败", tunnel.VKey, tunnel.Port, err)
			continue
		}
		_ = registry.Set(tunnel.VKey, svc)
		fmt.Println("P2P隧道启动成功", tunnel.VKey, "0.0.0.0:"+tunnel.Port)
	}
	<-signal.Free()
}
