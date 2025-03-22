package router

import (
	"github.com/gin-gonic/gin"
	"server/bootstrap"
	"server/global"
	"server/pkg/logger"
	"server/pkg/middleware"
	"server/router/admin"
	"server/router/auth"
	"server/router/normal"
	"server/router/public"
	"strings"
)

func init() {
	bootstrap.Route = func(engine *gin.Engine) {
		httpLog := logger.NewLogger(logger.Option{
			To:         []string{global.BASE_PATH + "/data/httpLog/access.log"},
			Level:      "debug",
			MaxSize:    100,
			MaxAge:     30,
			MaxBackups: 10,
			Compress:   true,
		})
		engine.Use(middleware.Logger(httpLog, global.MODE == "dev", func(c *gin.Context) bool {
			return strings.Contains(c.Request.RequestURI, "api/v")
		}))

		InitStatic(engine)

		adminGroup := engine.Group("api/v1/admin")
		admin.InitGostClient(adminGroup)
		admin.InitGostClientLogger(adminGroup)
		admin.InitGostClientForward(adminGroup)
		admin.InitGostClientHost(adminGroup)
		admin.InitGostClientTunnel(adminGroup)
		admin.InitGostClientProxy(adminGroup)
		admin.InitGostClientP2P(adminGroup)
		admin.InitGostNodeConfig(adminGroup)
		admin.InitGostNode(adminGroup)
		admin.InitGostNodeBind(adminGroup)
		admin.InitGostNodeLogger(adminGroup)
		admin.InitGostNodeRule(adminGroup)
		admin.InitSystemUser(adminGroup)
		admin.InitSystemNotice(adminGroup)
		admin.InitSystemConfig(adminGroup)
		admin.InitDashboard(adminGroup)

		publicGroup := engine.Group("api/v1/public")
		public.InitSystemConfig(publicGroup)
		public.InitGost(publicGroup)
		public.InitP2P(publicGroup)

		authGroup := engine.Group("api/v1/auth")
		auth.InitAuth(authGroup)

		normalGroup := engine.Group("api/v1/normal")
		normal.InitGost(normalGroup)
		normal.InitGostClient(normalGroup)
		normal.InitGostClientLogger(normalGroup)
		normal.InitGostClientForward(normalGroup)
		normal.InitGostClientHost(normalGroup)
		normal.InitGostClientTunnel(normalGroup)
		normal.InitGostClientProxy(normalGroup)
		normal.InitGostClientP2P(normalGroup)
		normal.InitGostNode(normalGroup)
		normal.InitGostObs(normalGroup)
		normal.InitSystemNotice(normalGroup)
	}
}
