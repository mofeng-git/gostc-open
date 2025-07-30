package main

import (
	"gorm.io/gen"
	"server/model"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./repository/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.ApplyBasic(
		model.GostAuth{},
		model.GostClient{},
		model.GostClientAdmission{},
		model.GostClientConfig{},
		model.GostClientForward{},
		model.GostClientHost{},
		model.GostClientHostDomain{},
		model.GostClientTunnel{},
		model.GostClientProxy{},
		model.GostClientP2P{},
		model.GostNode{},
		model.GostNodeBind{},
		model.GostNodeConfig{},
		model.GostNodeDomain{},
		model.GostNodePort{},
		model.GostObs{},
		model.SystemConfig{},
		model.SystemNotice{},
		model.SystemUser{},
		model.SystemUserCheckin{},
		model.SystemUserEmail{},
		model.FrpClientCfg{},
	)
	g.Execute()
}
