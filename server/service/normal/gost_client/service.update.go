package service

import (
	"errors"
	"go.uber.org/zap"
	"server/pkg/jwt"
	"server/repository"
)

type UpdateReq struct {
	Code string `binding:"required" json:"code"`
	Name string `binding:"required" json:"name" label:"名称"`
}

func (service *service) Update(claims jwt.Claims, req UpdateReq) error {
	db, _, log := repository.Get("")
	client, _ := db.GostClient.Where(db.GostClient.Code.Eq(req.Code), db.GostClient.UserCode.Eq(claims.Code)).First()
	if client == nil {
		return errors.New("客户端错误")
	}

	client.Name = req.Name
	if err := db.GostClient.Save(client); err != nil {
		log.Error("修改客户端失败", zap.Error(err))
		return errors.New("操作失败")
	}
	return nil
}
