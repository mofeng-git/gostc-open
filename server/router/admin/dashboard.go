package admin

import (
	"github.com/gin-gonic/gin"
	"server/controller/admin/dashboard"
	"server/global"
	"server/router/middleware"
)

func InitDashboard(group *gin.RouterGroup) {
	g := group.Group("dashboard", middleware.Auth(global.Jwt), middleware.AuthAdmin())
	g.POST("userObs", dashboard.UserObs)
	g.POST("nodeObs", dashboard.NodeObs)
	g.POST("userObsDate", dashboard.UserObsDate)
	g.POST("nodeObsDate", dashboard.NodeObsDate)
	g.POST("count", dashboard.Count)
}
