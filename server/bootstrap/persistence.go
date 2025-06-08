package bootstrap

import (
	"go.uber.org/zap"
	"server/global"
	"server/model"
	"server/pkg/orm"
	"server/pkg/orm/mysql"
	"server/pkg/orm/sqlite"
)

func InitPersistence() {
	var err error
	switch global.Config.DbType {
	case "mysql", "Mysql":
		global.DB, err = mysql.NewDB(
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
	case "sqlite", "Sqlite":
		global.DB, err = sqlite.NewDB(
			global.Config.Sqlite.File,
			global.Config.Sqlite.LogLevel,
			global.BASE_PATH+"/data/sql.log",
			global.MODE == "dev",
		)
	default:
		global.Logger.Fatal("init persistence fail", zap.Any("config", global.Config))
		Release()
	}
	if err != nil {
		global.Logger.Fatal("init persistence fail", zap.Any("config", global.Config), zap.Error(err))
		Release()
	}
	global.Logger.Info("init persistence finish")
	releaseFunc = append(releaseFunc, func() {
		global.DB.Close()
	})

	if err = autoMigrate(global.DB,
		&model.SystemUser{},
		&model.SystemUserEmail{},
		&model.SystemUserCheckin{},
		&model.SystemConfig{},
		&model.SystemNotice{},
		&model.GostAuth{},
		&model.GostObs{},
		&model.GostClient{},
		&model.GostClientHost{},
		&model.GostClientHostDomain{},
		&model.GostClientForward{},
		&model.GostClientTunnel{},
		&model.GostClientProxy{},
		&model.GostClientP2P{},
		&model.GostNode{},
		&model.GostNodeBind{},
		&model.GostNodeDomain{},
		&model.GostNodePort{},
		&model.GostNodeConfig{},
	); err != nil {
		global.Logger.Fatal("init table struct fail", zap.Error(err))
		Release()
	}
	global.Logger.Info("init table struct finish")
}

func autoMigrate(db orm.Interface, models ...any) error {
	if err := db.AutoMigrate(models...); err != nil {
		return err
	}
	for _, m := range models {
		db.GetDB().Model(m).Where("version IS NULL").Update("version", 1)
	}
	return nil
}
