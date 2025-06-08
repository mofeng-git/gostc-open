package middleware

import (
	"github.com/gin-gonic/gin"
	"server/model"
	"server/pkg/bean"
	"server/service/common/cache"
)

func AuthApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := GetApiKey(c)
		if key == "" {
			bean.Response.Fail(c, "API KEY不存在")
			c.Abort()
			return
		}

		var baseConfig model.SystemConfigBase
		cache.GetSystemConfigBase(&baseConfig)
		if baseConfig.ApiKey == "" {
			bean.Response.Fail(c, "未配置API KEY")
			c.Abort()
			return
		}

		if baseConfig.ApiKey != key {
			bean.Response.Fail(c, "API KEY错误")
			c.Abort()
			return
		}
	}
}

func GetApiKey(c *gin.Context) string {
	return c.GetHeader("key")
}
