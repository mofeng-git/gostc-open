package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"server/bootstrap"
	"server/global"
	"server/model"
	"server/pkg/utils"
)

var ResetPasswdCmd = cobra.Command{
	Use: "reset-pwd",
	Run: func(cmd *cobra.Command, args []string) {
		global.Init()
		bootstrap.InitLogger()
		bootstrap.InitConfig()
		bootstrap.InitJwt()
		bootstrap.InitPersistence()

		if len(args) == 0 {
			fmt.Println(`./server reset-pwd <account>`)
			return
		}

		var user model.SystemUser
		if global.DB.GetDB().Where("account = ?", args[0]).First(&user).RowsAffected == 0 {
			fmt.Println("account:", args[0], "non-existent")
			return
		}
		newPwd := utils.RandStr(8, utils.AllDict)
		user.Salt = utils.RandStr(8, utils.AllDict)
		user.Password = utils.MD5AndSalt(newPwd, user.Salt)
		if err := global.DB.GetDB().Save(&user).Error; err != nil {
			fmt.Println("reset password error:", err)
			return
		}
		fmt.Println("account:", args[0])
		fmt.Println("new password:", newPwd)
		fmt.Println("reset password success")
		defer bootstrap.Release()
	},
}
