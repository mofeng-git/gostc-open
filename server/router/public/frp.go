package public

import (
	"github.com/gin-gonic/gin"
	"server/controller/public/frp"
)

func InitFrp(group *gin.RouterGroup) {
	g := group.Group("frp")
	g.Any("login", frp.Login)
	g.Any("newProxy", frp.NewProxy)
	g.Any("closeProxy", frp.CloseProxy)
	g.Any("ping", frp.Ping)
	g.Any("newWorkConn", frp.NewWorkConn)
	g.Any("newUserConn", frp.NewUserConn)
	g.Any("visitorTunnel", frp.VisitorTunnel)
	g.Any("visitorP2P", frp.VisitorP2P)
}
