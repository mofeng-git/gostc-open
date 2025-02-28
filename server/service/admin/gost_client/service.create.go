package service

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
)

type CreateReq struct {
	Name     string `binding:"required" json:"name" label:"名称"`
	UserCode string `binding:"required" json:"userCode" label:"用户编号"`
}

func (service *service) Create(req CreateReq) error {
	db, _, log := repository.Get("")
	if err := db.Create(&model.GostClient{
		Name:     req.Name,
		UserCode: req.UserCode,
		Key:      uuid.NewString(),
	}).Error; err != nil {
		log.Error("新增用户客户端失败", zap.Error(err))
		return errors.New("操作失败")
	}
	return nil
}
