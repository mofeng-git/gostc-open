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
	RootCmd.AddCommand(&ResetPasswdCmd)
	RootCmd.AddCommand(&CloseTotp)
	RootCmd.AddCommand(&MigrateCmd)
	for _, cmd := range []*cobra.Command{
		&RootCmd,
		&VersionCmd,
		&ServiceCmd,
		&ResetPasswdCmd,
		&CloseTotp,
		&MigrateCmd,
	} {
		cmd.Flags().StringVarP(&global.BASE_PATH, "path", "p", "", "app run dir")
		cmd.Flags().StringVarP(&global.LOGGER_LEVEL, "log-level", "", "warn", "log level debug|info|warn|error|fatal")
		cmd.Flags().BoolVarP(&global.FLAG_DEV, "dev", "d", false, "app run dev")
		cmd.Flags().StringVarP(&global.LICENCE, "licence", "l", "", "app licence")
		cmd.Flags().StringVarP(&global.LICENCE_URL, "licence-url", "", "https://licence.sian.one", "app licence url")
		cmd.Flags().StringVarP(&global.LICENCE_FILE, "licence-file", "", "", "load licence from file")
		cmd.Flags().StringVarP(&global.LICENCE_PROXY, "licence-proxy", "", "", "load licence use http proxy, example: http://username:password@proxy.example.com:8080")
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
		bootstrap.InitTodo()
		bootstrap.InitTask()
		bootstrap.InitRouter()
		bootstrap.InitServer()

		<-signal.Free()
		bootstrap.Release()
	},
}
