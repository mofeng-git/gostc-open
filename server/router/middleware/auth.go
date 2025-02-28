package middleware

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/pkg/jwt"
)

func Auth(tool *jwt.Tool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetToken(c)
		if token == "" {
			bean.Response.AuthNoLogin(c)
			c.Abort()
			return
		}
		claims, err := tool.ValidToken(token)
		if err != nil {
			bean.Response.AuthInvalid(c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
	}
}

func GetToken(c *gin.Context) string {
	return c.GetHeader("token")
}

func GetClaims(c *gin.Context) jwt.Claims {
	claims, exists := c.Get("claims")
	if !exists {
		panic("get claims fail")
	}
	return claims.(jwt.Claims)
}
