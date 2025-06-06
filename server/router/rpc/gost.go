package rpc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	"server/service/gost_engine"
	"time"
)

func checkRegister(client *arpc.Client) {
	time.Sleep(time.Second * 10)
	_, ok := client.Get("registered")
	if !ok {
		_ = client.Notify("stop", "客户端连接注册超时", time.Second*5)
		time.Sleep(time.Second * 2)
		client.Stop()
	}
}

func InitGost(engine *gin.Engine, ln net.Listener) {
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
		code, ok := client.Get("code")
		if ok {
			cache.SetClientOnline(code.(string), false, cache2.NoExpiration)
			cache.SetNodeOnline(code.(string), false, cache2.NoExpiration)
		}
	})

	// 路由
	rpcServer.Handler.Handle("rpc/client/reg", func(c *arpc.Context) {
		var reply = make(map[string]string)
		if err := c.Bind(&reply); err != nil {
			return
		}
		key := reply["key"]
		version := reply["version"]
		if key == "" {
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

		value, ok := gost_engine.EngineRegistry.Get(gostClient.Code)
		if ok {
			value.GetClient().Stop("客户端已在别处连接")
			time.Sleep(time.Second)
		}

		e := gost_engine.NewARpcClientEngine(gostClient.Code, client)
		clientEngine := gost_engine.NewClientEngine(gostClient.Code, e)
		if gostClient.Code == "" {
			e.Stop("客户端不存在")
		} else {
			gost_engine.EngineRegistry.Set(clientEngine)
		}
		cache.SetClientLastTime(gostClient.Code)
		cache.SetClientOnline(gostClient.Code, true, cache2.NoExpiration)
		cache.SetClientVersion(gostClient.Code, version)
		client.Set("code", gostClient.Code)
		clientTunnelInit(gostClient.Code)
	})
	rpcServer.Handler.Handle("rpc/node/reg", func(c *arpc.Context) {
		var reply = make(map[string]string)
		if err := c.Bind(&reply); err != nil {
			return
		}
		key := reply["key"]
		version := reply["version"]
		if key == "" {
			return
		}

		client := c.Client
		client.Set("key", key)
		client.Set("version", version)
		client.Set("registered", 1)
		client.Set("session", uuid.NewString())

		db, _, _ := repository.Get("")
		gostNode, _ := db.GostNode.Where(db.GostNode.Key.Eq(key)).First()
		if gostNode == nil {
			gostNode = &model.GostNode{}
		}

		value, ok := gost_engine.EngineRegistry.Get(gostNode.Code)
		if ok {
			value.GetNode().Stop("节点已在别处连接")
			time.Sleep(time.Second)
		}

		e := gost_engine.NewARpcNodeEngine(gostNode.Code, "", client)
		clientEngine := gost_engine.NewNodeEngine(gostNode.Code, e)
		if gostNode.Code == "" {
			e.Stop("节点不存在")
		} else {
			gost_engine.EngineRegistry.Set(clientEngine)
		}
		cache.SetNodeOnline(gostNode.Code, true, cache2.NoExpiration)
		cache.SetNodeVersion(gostNode.Code, version)
		client.Set("code", gostNode.Code)
		nodeServerInit(gostNode.Code)
	})

	// 入口
	engine.GET("rpc/client/ws", func(c *gin.Context) {
		ln.(*websocket.Listener).Handler(c.Writer, c.Request)
	})
	engine.GET("rpc/node/ws", func(c *gin.Context) {
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

func clientTunnelInit(code string) {
	db, _, _ := repository.Get("")
	var hostCodes []string
	_ = db.GostClientHost.Where(db.GostClientHost.ClientCode.Eq(code)).Pluck(db.GostClientHost.Code, &hostCodes)
	for _, code := range hostCodes {
		gost_engine.ClientHostConfig(db, code)
	}
	var forwardCodes []string
	_ = db.GostClientForward.Where(db.GostClientForward.ClientCode.Eq(code)).Pluck(db.GostClientForward.Code, &forwardCodes)
	for _, code := range forwardCodes {
		gost_engine.ClientForwardConfig(db, code)
	}
	var tunnelCodes []string
	_ = db.GostClientTunnel.Where(db.GostClientTunnel.ClientCode.Eq(code)).Pluck(db.GostClientTunnel.Code, &tunnelCodes)
	for _, code := range tunnelCodes {
		gost_engine.ClientTunnelConfig(db, code)
	}
	var proxyCodes []string
	_ = db.GostClientProxy.Where(db.GostClientProxy.ClientCode.Eq(code)).Pluck(db.GostClientProxy.Code, &proxyCodes)
	for _, code := range proxyCodes {
		gost_engine.ClientProxyConfig(db, code)
	}
	var p2pCodes []string
	_ = db.GostClientP2P.Where(db.GostClientP2P.ClientCode.Eq(code)).Pluck(db.GostClientP2P.Code, &p2pCodes)
	for _, code := range p2pCodes {
		gost_engine.ClientP2PConfig(db, code)
	}
}

func nodeServerInit(code string) {
	db, _, _ := repository.Get("")
	gost_engine.NodeConfig(db, code)
}
