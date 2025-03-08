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
	RootCmd.AddCommand(&VersionCmd)
	RootCmd.AddCommand(&ServiceCmd)
	for _, cmd := range []*cobra.Command{
		&RootCmd,
		&VersionCmd,
		&ServiceCmd,
	} {
		cmd.Flags().StringVarP(&global.BASE_PATH, "path", "p", "", "app run dir")
		cmd.Flags().StringVarP(&global.LOGGER_LEVEL, "log-level", "", "warn", "log level debug|info|warn|error|fatal")
		cmd.Flags().BoolVarP(&global.FLAG_DEV, "dev", "d", false, "app run dev")
	}
}

var RootCmd = cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		global.Init()
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
