package todo

import (
	"server/bootstrap"
	"server/global"
	"server/repository/query"
)

func init() {
	bootstrap.TodoFunc = func() {
		query.SetDefault(global.DB.GetDB())
		systemUser()
		systemConfig()
		gostClient()
		gostClientLogger()
		gostNodeLogger()
		gostNodePort()
	}
}
