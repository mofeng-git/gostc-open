package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"server/bootstrap"
	"server/global"
	"server/model"
)

var CloseTotp = cobra.Command{
	Use: "close-totp",
	Run: func(cmd *cobra.Command, args []string) {
		global.Init()
		bootstrap.InitLogger()
		bootstrap.InitConfig()
		bootstrap.InitJwt()
		bootstrap.InitPersistence()

		if len(args) == 0 {
			fmt.Println(`./server close-totp <account>`)
			return
		}

		var user model.SystemUser
		if global.DB.GetDB().Where("account = ?", args[0]).First(&user).RowsAffected == 0 {
			fmt.Println("account:", args[0], "non-existent")
			return
		}
		if user.OtpKey == "" {
			fmt.Println("account:", args[0], "not enabled totp")
			return
		}

		user.OtpKey = ""
		if err := global.DB.GetDB().Save(&user).Error; err != nil {
			fmt.Println("close totp error:", err)
			return
		}
		fmt.Println("account:", args[0])
		fmt.Println("close totp success")
		defer bootstrap.Release()
	},
}
