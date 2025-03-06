package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"server/bootstrap"
	"server/global"
	"server/pkg/signal"
	_ "server/router"
	_ "server/task"
	_ "server/todo"
)

// [SOURCE] https://patorjk.com/software/taag/#p=display&h=0&v=0&f=ANSI%20Shadow&t=GOSTC
func init() {
	fmt.Println(`
 ██████╗  ██████╗ ███████╗████████╗ ██████╗
██╔════╝ ██╔═══██╗██╔════╝╚══██╔══╝██╔════╝
██║  ███╗██║   ██║███████╗   ██║   ██║     
██║   ██║██║   ██║╚════██║   ██║   ██║     
╚██████╔╝╚██████╔╝███████║   ██║   ╚██████╗
 ╚═════╝  ╚═════╝ ╚══════╝   ╚═╝    ╚═════╝
`)
}

func init() {
	global.CmdPerInit(&RootCmd, &VersionCmd)
	RootCmd.AddCommand(&VersionCmd)
}

var RootCmd = cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		global.CmdInit()
		bootstrap.InitLogger()
		bootstrap.InitConfig()
		bootstrap.InitJwt()
		bootstrap.InitPersistence()
		bootstrap.InitMemory()
		bootstrap.InitTask()
		bootstrap.InitRouter()
		bootstrap.InitTodo()
		bootstrap.InitServer()

		<-signal.Free()
		bootstrap.Release()
	},
}
