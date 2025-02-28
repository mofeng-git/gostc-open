package middleware

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
)

func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims.Data == nil {
			bean.Response.AuthNoAllow(c)
			c.Abort()
			return
		}
		if claims.Data["admin"] != "1" {
			bean.Response.AuthNoAllow(c)
			c.Abort()
			return
		}
	}
}
