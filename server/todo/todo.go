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
		gostNode()
		gostNodePort()

		// 需要先将obs回写到cache，在处理obs数据
		gostObsWriteBack()
		gostObs()

		// 修复一些线上数据错误
		//fix()
	}
}
