package normal

import (
	"github.com/gin-gonic/gin"
	"server/controller/normal/frp_client_cfg"
	"server/global"
	"server/router/middleware"
)

func InitFrpClientCfg(group *gin.RouterGroup) {
	g := group.Group("frp/client/cfg", middleware.Auth(global.Jwt))
	g.POST("create", frp_client_cfg.Create)
	g.POST("migrate", frp_client_cfg.Migrate)
	g.POST("update", frp_client_cfg.Update)
	g.POST("enable", frp_client_cfg.Enable)
	g.POST("delete", frp_client_cfg.Delete)
	g.POST("page", frp_client_cfg.Page)
}
