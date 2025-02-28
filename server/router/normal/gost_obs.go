package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/gost_obs"
	"server/global"
	"server/router/middleware"
)

func InitGostObs(group *gin.RouterGroup) {
	g := group.Group("gost/obs", middleware.Auth(global.Jwt))
	g.POST("tunnel/month", gost_obs.TunnelMonth)
	g.POST("client/month", gost_obs.ClientMonth)
	g.POST("node/month", gost_obs.NodeMonth)
	g.POST("user/month", gost_obs.UserMonth)
}
