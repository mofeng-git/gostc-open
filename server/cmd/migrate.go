package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"server/bootstrap"
	"server/global"
	"server/model"
	"server/pkg/orm"
	"server/pkg/orm/mysql"
	"server/pkg/orm/sqlite"
)

var MigrateCmd = cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {
		global.Init()
		bootstrap.InitLogger()
		bootstrap.InitConfig()
		bootstrap.InitJwt()
		defer bootstrap.Release()

		var validTypeFunc = func(dbType string) bool {
			switch dbType {
			case "sqlite":
				return true
			case "mysql":
				return true
			default:
				return false
			}
		}

		if len(args) != 2 {
			fmt.Println(`./server migrate sourceDbType targetDbType`)
			fmt.Println(`dbType: sqlite|mysql`)
			return
		}

		var sourceDbType, targetDbType string
		for index, item := range args {
			switch index {
			case 0:
				if !validTypeFunc(item) {
					fmt.Println(`dbType: sqlite|mysql`)
					return
				}
				sourceDbType = item
			case 1:
				if !validTypeFunc(item) {
					fmt.Println(`dbType: sqlite|mysql`)
					return
				}
				targetDbType = item
			}
		}

		var dbMap = make(map[string]orm.Interface)
		mysqlDB, err := mysql.NewDB(
			mysql.Config{
				Username: global.Config.Mysql.User,
				Password: global.Config.Mysql.Pwd,
				Host:     global.Config.Mysql.Host,
				Port:     global.Config.Mysql.Port,
				Prefix:   global.Config.Mysql.Prefix,
				Extend:   global.Config.Mysql.Extend,
				DbName:   global.Config.Mysql.DB,
			},
			global.Config.Mysql.LogLevel,
			global.BASE_PATH+"/data/sql.log",
			global.MODE == "dev",
		)
		if err != nil {
			fmt.Println("connect mysql fail,error:", err)
			marshal, _ := json.Marshal(global.Config.Mysql)
			fmt.Println("config:", string(marshal))
			return
		}
		dbMap["mysql"] = mysqlDB
		sqliteDB, err := sqlite.NewDB(
			global.Config.Sqlite.File,
			global.Config.Sqlite.LogLevel,
			global.BASE_PATH+"/data/sql.log",
			global.MODE == "dev",
		)
		if err != nil {
			fmt.Println("connect sqlite fail,error:", err)
			marshal, _ := json.Marshal(global.Config.Sqlite)
			fmt.Println("config:", string(marshal))
			return
		}
		dbMap["sqlite"] = sqliteDB

		db := dbMap[targetDbType].GetDB()
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostAuth{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostClient{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostClientForward{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostClientHost{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostClientHostDomain{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostClientP2P{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostClientProxy{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostClientTunnel{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostNode{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostNodeBind{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostNodeConfig{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostNodeDomain{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostNodePort{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.GostObs{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.SystemConfig{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.SystemNotice{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.SystemUser{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.SystemUserEmail{}); err != nil {
				return err
			}
			if err := Migrate(dbMap[sourceDbType].GetDB(), tx, model.SystemUserCheckin{}); err != nil {
				return err
			}
			return nil
		}); err != nil {
			fmt.Println("migrate fail,error:", err)
			return
		}
		fmt.Println(fmt.Sprintf("migrate %s to %s success", sourceDbType, targetDbType))
	},
}

func Migrate[Model model.GostAuth |
	model.GostClient |
	model.GostClientForward |
	model.GostClientHost |
	model.GostClientHostDomain |
	model.GostClientP2P |
	model.GostClientProxy |
	model.GostClientTunnel |
	model.GostNode |
	model.GostNodeBind |
	model.GostNodeConfig |
	model.GostNodeDomain |
	model.GostNodePort |
	model.GostObs |
	model.SystemConfig |
	model.SystemNotice |
	model.SystemUser |
	model.SystemUserEmail |
	model.SystemUserCheckin](source *gorm.DB, target *gorm.DB, model Model) error {
	if err := target.AutoMigrate(&model); err != nil {
		return err
	}
	var oldList []Model
	source.Find(&oldList)
	return target.CreateInBatches(&oldList, 1000).Error
}
