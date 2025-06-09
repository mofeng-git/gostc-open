package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/pkg/utils"
	"server/repository"
	"server/repository/query"
	"server/service/engine"
)

type CreateReq struct {
	Name       string `binding:"required" json:"name" label:"名称"`
	ClientCode string `binding:"required" json:"clientCode"`
	Type       string `binding:"required" json:"type"`
	Content    string `binding:"required" json:"content"`
}

func (service *service) Create(claims jwt.Claims, req CreateReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}
		client, _ := tx.GostClient.Where(
			tx.GostClient.UserCode.Eq(claims.Code),
			tx.GostClient.Code.Eq(req.ClientCode),
		).First()
		if client == nil {
			return errors.New("客户端错误")
		}

		var cfg = model.FrpClientCfg{
			Name:        req.Name,
			Enable:      1,
			Content:     req.Content,
			ContentType: req.Type,
			ClientCode:  req.ClientCode,
			UserCode:    claims.Code,
		}
		if err := tx.FrpClientCfg.Create(&cfg); err != nil {
			log.Error("新增用户配置隧道失败", zap.Error(err))
			return errors.New("操作失败")
		}
		var auth = model.GostAuth{
			TunnelType: model.GOST_TUNNEL_TYPE_HOST,
			TunnelCode: cfg.Code,
			User:       utils.RandStr(10, utils.AllDict),
			Password:   utils.RandStr(10, utils.AllDict),
		}
		if err := tx.GostAuth.Create(&auth); err != nil {
			log.Error("生成授权信息失败", zap.Error(err))
			return errors.New("操作失败")
		}
		return engine.ClientCfgConfig(tx, cfg.Code)
	})
}
