package public

import (
	"github.com/gin-gonic/gin"
	"server/controller/public/gost"
)

func InitGost(group *gin.RouterGroup) {
	g := group.Group("gost")
	g.Any("auther", gost.Auther)
}
