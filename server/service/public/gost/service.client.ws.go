package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/service/gost_engine"
)

func (service *service) ClientWs(c *gin.Context) {
	db, _, log := repository.Get("")
	key := c.GetHeader("key")
	if key == "" {
		return
	}

	client, _ := db.GostClient.Where(db.GostClient.Key.Eq(key)).First()
	if client == nil {
		client = &model.GostClient{}
	}

	value, ok := gost_engine.EngineRegistry.Get(client.Code)
	if ok {
		value.GetClient().Stop("客户端已在别处连接，连接IP：" + c.ClientIP())
	}

	engine, err := gost_engine.NewGwsClientEngine(client.Code, c.Writer, c.Request)
	if err != nil {
		log.Error("建立连接失败", zap.String("key", key), zap.Error(err))
		return
	}
	clientEngine := gost_engine.NewClientEngine(client.Code, engine)

	//engine, err := gost_engine.NewGwsClientEngine(client.Code, c.Writer, c.Request, gost_engine.NewClientEvent(client.Code, c.ClientIP(), log))
	//if err != nil {
	//	log.Error("建立连接失败", zap.String("key", key), zap.Error(err))
	//	return
	//}
	if client.Code == "" {
		engine.Stop("客户端不存在")
	} else {
		gost_engine.EngineRegistry.Set(clientEngine)
	}
}
