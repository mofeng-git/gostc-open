package bootstrap

import "server/global"

var TodoFunc func()

func InitTodo() {
	if TodoFunc != nil {
		TodoFunc()
	}
	global.Logger.Info("init todo finish")
}
