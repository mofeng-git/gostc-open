package public

import (
	"github.com/gin-gonic/gin"
	"server/controller/public/p2p"
)

func InitP2P(group *gin.RouterGroup) {
	g := group.Group("p2p")
	g.Any("new", p2p.New)
	g.Any("visit", p2p.Visit)
}
