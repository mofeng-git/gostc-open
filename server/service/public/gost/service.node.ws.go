package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/service/gost_engine"
)

func (service *service) NodeWs(c *gin.Context) {
	db, _, log := repository.Get("")
	key := c.GetHeader("key")
	if key == "" {
		return
	}
	var node model.GostNode
	db.Where("key = ?", key).First(&node)

	value, ok := gost_engine.EngineRegistry.Get(node.Code)
	if ok {
		value.Close("节点已在别处连接，连接IP：" + c.ClientIP())
	}
	engine, err := gost_engine.NewEngine(node.Code, c.Writer, c.Request, gost_engine.NewNodeEvent(node.Code, log))
	if err != nil {
		log.Error("建立连接失败", zap.String("key", key), zap.Error(err))
		return
	}
	if node.Code == "" {
		engine.Close("节点不存在")
	} else {
		gost_engine.EngineRegistry.Set(engine)
	}
	go engine.ReadLoop()
}
