package public

import (
	"github.com/gin-gonic/gin"
	"server/controller/public/gost"
)

func InitGost(group *gin.RouterGroup) {
	g := group.Group("gost")
	g.Any("client/ws", gost.ClientWs)
	g.Any("node/ws", gost.NodeWs)
	g.Any("ingress", gost.Ingress)
	g.Any("auther", gost.Auther)
	g.Any("admission", gost.Admission)
	g.Any("limiter", gost.Limiter)
	g.Any("obs", gost.Obs)
	g.Any("visit", gost.Visit)
}
