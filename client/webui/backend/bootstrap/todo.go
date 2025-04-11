package bootstrap

import "gostc-sub/webui/backend/global"

var TodoFunc func()

func InitTodo() {
	if TodoFunc != nil {
		TodoFunc()
	}
	global.Logger.Info("init todo finish")
}
