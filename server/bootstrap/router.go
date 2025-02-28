package bootstrap

import (
	"github.com/gin-gonic/gin"
	"server/global"
	"server/pkg/logger"
	"server/pkg/middleware"
	"server/router"
	"server/router/admin"
	"server/router/auth"
	"server/router/normal"
	"server/router/public"
	"strings"
)

var engine *gin.Engine

func InitRouter() {
	if global.MODE == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	engine = gin.Default()
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

	router.InitStatic(engine)

	adminGroup := engine.Group("api/v1/admin")
	admin.InitGostClient(adminGroup)
	admin.InitGostClientLogger(adminGroup)
	admin.InitGostClientForward(adminGroup)
	admin.InitGostClientHost(adminGroup)
	admin.InitGostClientTunnel(adminGroup)
	admin.InitGostNodeConfig(adminGroup)
	admin.InitGostNode(adminGroup)
	admin.InitGostNodeBind(adminGroup)
	admin.InitGostNodeLogger(adminGroup)
	admin.InitGostNodeRule(adminGroup)
	admin.InitGostUserConfig(adminGroup)
	admin.InitSystemUser(adminGroup)
	admin.InitSystemNotice(adminGroup)
	admin.InitSystemConfig(adminGroup)

	publicGroup := engine.Group("api/v1/public")
	public.InitSystemConfig(publicGroup)
	public.InitGost(publicGroup)

	authGroup := engine.Group("api/v1/auth")
	auth.InitAuth(authGroup)

	normalGroup := engine.Group("api/v1/normal")
	normal.InitGost(normalGroup)
	normal.InitGostClient(normalGroup)
	normal.InitGostClientLogger(normalGroup)
	normal.InitGostClientForward(normalGroup)
	normal.InitGostClientHost(normalGroup)
	normal.InitGostClientTunnel(normalGroup)
	normal.InitGostNode(normalGroup)
	normal.InitGostObs(normalGroup)
	normal.InitSystemNotice(normalGroup)

	global.Logger.Info("init router finish")
}
