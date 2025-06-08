package rpc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lesismal/arpc"
	arpcLog "github.com/lesismal/arpc/log"
	cache2 "github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"net"
	"server/global"
	"server/model"
	"server/pkg/rpc_protocol/websocket"
	"server/repository"
	"server/service/common/cache"
	"server/service/engine"
	"strings"
	"time"
)

func checkRegister(client *arpc.Client) {
	time.Sleep(time.Second * 120)
	_, ok := client.Get("registered")
	if !ok {
		_ = client.Notify("stop", "客户端连接注册超时", time.Second*5)
		time.Sleep(time.Second * 2)
		client.Stop()
	}
}

func InitFrp(ginEngine *gin.Engine, ln net.Listener) {
	if global.MODE == "dev" {
		arpcLog.SetLevel(arpcLog.LevelDebug)
	} else {
		arpcLog.SetLevel(arpcLog.LevelNone)
	}
	rpcServer := arpc.NewServer()
	rpcServer.Handler.SetReadTimeout(time.Second * 50)
	rpcServer.Handler.HandleConnected(func(client *arpc.Client) {
		go checkRegister(client)
	})
	rpcServer.Handler.HandleDisconnected(func(client *arpc.Client) {
		fmt.Println("HandleDisconnected")
		code, ok := client.Get("code")
		if ok {
			cache.SetClientOnline(code.(string), false, cache2.NoExpiration)
			cache.SetNodeOnline(code.(string), false, cache2.NoExpiration)
		}
	})
	rpcServer.Handler.HandleMessageDone(func(c *arpc.Client, m *arpc.Message) {
		//fmt.Println(string(m.Data()))
	})

	// 路由
	rpcServer.Handler.Handle("rpc/client/reg", func(c *arpc.Context) {
		var reply = make(map[string]string)
		if err := c.Bind(&reply); err != nil {
			_ = c.Write(err.Error())
			return
		}
		key := reply["key"]
		version := reply["version"]
		if key == "" {
			_ = c.Write("no key")
			return
		}

		client := c.Client
		client.Set("key", key)
		client.Set("version", version)
		client.Set("registered", 1)

		db, _, _ := repository.Get("")
		gostClient, _ := db.GostClient.Where(db.GostClient.Key.Eq(key)).First()
		if gostClient == nil {
			gostClient = &model.GostClient{}
		}

		value, ok := engine.EngineRegistry.Get(gostClient.Code)
		if ok {
			value.GetClient().Stop("客户端已在别处连接")
			time.Sleep(time.Second * 2)
		}

		e := engine.NewARpcClientEngine(gostClient.Code, client)
		clientEngine := engine.NewClientEngine(gostClient.Code, e)
		if gostClient.Code == "" {
			_ = c.Write("客户端不存在")
			return
		} else {
			engine.EngineRegistry.Set(clientEngine)
		}
		_ = c.Write("success")
		cache.SetClientLastTime(gostClient.Code)
		cache.SetClientOnline(gostClient.Code, true, cache2.NoExpiration)
		cache.SetClientVersion(gostClient.Code, version)
		client.Set("code", gostClient.Code)
		clientTunnelInit(gostClient.Code)
	})
	rpcServer.Handler.Handle("rpc/node/reg", func(c *arpc.Context) {
		var reply = make(map[string]string)
		if err := c.Bind(&reply); err != nil {
			_ = c.Write(err.Error())
			return
		}
		key := reply["key"]
		version := reply["version"]
		domain := reply["domain"]
		domainCache := reply["domain_cache"]
		if key == "" {
			_ = c.Write("no key")
			return
		}

		client := c.Client
		client.Set("key", key)
		client.Set("version", version)
		client.Set("registered", 1)

		db, _, _ := repository.Get("")

		gostNode, _ := db.GostNode.Where(db.GostNode.Key.Eq(key)).First()
		if gostNode == nil {
			_ = c.Write("节点不存在")
			return
		}

		value, ok := engine.EngineRegistry.Get(gostNode.Code)
		if ok {
			value.GetNode().Stop("节点已在别处连接")
			time.Sleep(time.Second)
		}

		e := engine.NewARpcNodeEngine(gostNode.Code, "", client)
		clientEngine := engine.NewNodeEngine(gostNode.Code, e)
		if gostNode.Code == "" {
			_ = c.Write("节点不存在")
			return
		} else {
			engine.EngineRegistry.Set(clientEngine)
		}
		_ = c.Write("success")
		cache.SetNodeOnline(gostNode.Code, true, cache2.NoExpiration)
		cache.SetNodeVersion(gostNode.Code, version)
		cache.SetNodeCustomDomain(gostNode.Code, domain == "1")
		cache.SetNodeCache(gostNode.Code, domainCache == "1")
		client.Set("code", gostNode.Code)
		nodeServerInit(gostNode.Code)
	})

	rpcServer.Handler.Handle("rpc/metrics/input", func(c *arpc.Context) {
		var req TrafficOutput
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		code := strings.Split(req.Name, "_")[0]
		tunnelInfo := cache.GetTunnelInfo(code)
		go cache.IncreaseObs(time.Now().Format(time.DateOnly), tunnelInfo.Code, tunnelInfo.ClientCode, tunnelInfo.NodeCode, tunnelInfo.UserCode, cache.TunnelObs{
			InputBytes:  req.Total,
			OutputBytes: 0,
		})
		_ = c.Write("success")
	})
	rpcServer.Handler.Handle("rpc/metrics/output", func(c *arpc.Context) {
		var req TrafficOutput
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		code := strings.Split(req.Name, "_")[0]
		tunnelInfo := cache.GetTunnelInfo(code)
		go cache.IncreaseObs(time.Now().Format(time.DateOnly), tunnelInfo.Code, tunnelInfo.ClientCode, tunnelInfo.NodeCode, tunnelInfo.UserCode, cache.TunnelObs{
			InputBytes:  0,
			OutputBytes: req.Total,
		})
		_ = c.Write("success")
	})

	// 入口
	ginEngine.GET("rpc/ws", func(c *gin.Context) {
		ln.(*websocket.Listener).Handler(c.Writer, c.Request)
	})

	go func() {
		if err := rpcServer.Serve(ln); err != nil {
			global.Logger.Fatal("rpc serve fail", zap.Error(err))
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	global.Logger.Info("rcp server on address: " + global.Config.Address)
}

type TrafficOutput struct {
	Name      string
	ProxyType string
	Total     int64
}

func clientTunnelInit(code string) {
	db, _, _ := repository.Get("")
	var hostCodes []string
	_ = db.GostClientHost.Where(db.GostClientHost.ClientCode.Eq(code)).Pluck(db.GostClientHost.Code, &hostCodes)
	for _, code := range hostCodes {
		engine.ClientHostConfig(db, code)
	}
	var forwardCodes []string
	_ = db.GostClientForward.Where(db.GostClientForward.ClientCode.Eq(code)).Pluck(db.GostClientForward.Code, &forwardCodes)
	for _, code := range forwardCodes {
		engine.ClientForwardConfig(db, code)
	}
	var tunnelCodes []string
	_ = db.GostClientTunnel.Where(db.GostClientTunnel.ClientCode.Eq(code)).Pluck(db.GostClientTunnel.Code, &tunnelCodes)
	for _, code := range tunnelCodes {
		engine.ClientTunnelConfig(db, code)
	}
	var proxyCodes []string
	_ = db.GostClientProxy.Where(db.GostClientProxy.ClientCode.Eq(code)).Pluck(db.GostClientProxy.Code, &proxyCodes)
	for _, code := range proxyCodes {
		engine.ClientProxyConfig(db, code)
	}
	var p2pCodes []string
	_ = db.GostClientP2P.Where(db.GostClientP2P.ClientCode.Eq(code)).Pluck(db.GostClientP2P.Code, &p2pCodes)
	for _, code := range p2pCodes {
		engine.ClientP2PConfig(db, code)
	}
}

func nodeServerInit(code string) {
	db, _, _ := repository.Get("")
	engine.NodeConfig(db, code)
}
