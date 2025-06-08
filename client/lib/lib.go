package lib

import (
	"fmt"
	"gostc-sub/webui/backend/bootstrap"
	_ "gostc-sub/webui/backend/router"
	_ "gostc-sub/webui/backend/todo"
	"runtime"
)

func init() {
	var basePath = "/data/user/0/one.sian.gostc"
	bootstrap.InitLogger()
	bootstrap.InitFS(basePath)
	bootstrap.InitTodo()
	bootstrap.InitTask()
	bootstrap.InitRouter()
}

var errMsg string

func Start(address string) {
	if err := bootstrap.StartServer(address); err != nil {
		errMsg = err.Error()
		return
	}
	errMsg = ""
	fmt.Println("start server")
}
func Stop() {
	bootstrap.StopServer()
	fmt.Println("stop server")
}

func IsRunning() string {
	if bootstrap.IsRunningServer() {
		return "1"
	}
	return "2"
}

func GetPlatform() string {
	return runtime.GOOS
}

func GetErrMsg() string {
	return errMsg
}
