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
	cache3 "server/repository/cache"
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
			cache3.SetClientOnline(code.(string), false, cache2.NoExpiration)
			cache3.SetNodeOnline(code.(string), false, cache2.NoExpiration)
		}
	})
	rpcServer.Handler.HandleMessageDone(func(c *arpc.Client, m *arpc.Message) {
		//fmt.Println(string(m.Data()))
	})

	// 路由
	rpcServer.Handler.Handle("rpc/client/ping", func(c *arpc.Context) {
		code, ok := c.Client.Get("code")
		if ok {
			cache3.SetClientOnline(code.(string), true, cache2.NoExpiration)
		}
		_ = c.Write("success")
	})
	rpcServer.Handler.Handle("rpc/node/ping", func(c *arpc.Context) {
		code, ok := c.Client.Get("code")
		if ok {
			cache3.SetNodeOnline(code.(string), true, cache2.NoExpiration)
		}
		_ = c.Write("success")
	})
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
		cache3.SetClientLastTime(gostClient.Code)
		cache3.SetClientOnline(gostClient.Code, true, cache2.NoExpiration)
		cache3.SetClientVersion(gostClient.Code, version)
		client.Set("code", gostClient.Code)
		engine.ClientAllConfigUpdateByClientCode(db, gostClient.Code)
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
		cache3.SetNodeOnline(gostNode.Code, true, cache2.NoExpiration)
		cache3.SetNodeVersion(gostNode.Code, version)
		cache3.SetNodeCustomDomain(gostNode.Code, domain == "1")
		cache3.SetNodeCache(gostNode.Code, domainCache == "1")
		client.Set("code", gostNode.Code)
		engine.NodeConfig(db, gostNode.Code)
	})

	rpcServer.Handler.Handle("rpc/metrics/input", func(c *arpc.Context) {
		var req TrafficOutput
		if err := c.Bind(&req); err != nil {
			_ = c.Write(err.Error())
			return
		}
		code := strings.Split(req.Name, "_")[0]
		tunnelInfo := cache3.GetTunnelInfo(code)
		go cache3.IncreaseObs(time.Now().Format(time.DateOnly), tunnelInfo.Code, tunnelInfo.ClientCode, tunnelInfo.NodeCode, tunnelInfo.UserCode, cache3.TunnelObs{
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
		tunnelInfo := cache3.GetTunnelInfo(code)
		go cache3.IncreaseObs(time.Now().Format(time.DateOnly), tunnelInfo.Code, tunnelInfo.ClientCode, tunnelInfo.NodeCode, tunnelInfo.UserCode, cache3.TunnelObs{
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
