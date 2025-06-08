package open

import (
	"github.com/gin-gonic/gin"
	"server/controller/open/auth"
	"server/router/middleware"
)

func InitAuth(group *gin.RouterGroup) {
	group.POST("checkin", middleware.AuthApi(), auth.Checkin)
}
