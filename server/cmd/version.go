package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"server/global"
)

var VersionCmd = cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		global.Init()
		fmt.Println("VERSION:", global.VERSION)
	},
}
