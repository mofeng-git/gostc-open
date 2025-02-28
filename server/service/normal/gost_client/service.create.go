package service

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
)

type CreateReq struct {
	Name string `binding:"required" json:"name" label:"名称"`
}

func (service *service) Create(claims jwt.Claims, req CreateReq) error {
	db, _, log := repository.Get("")
	if err := db.Create(&model.GostClient{
		Name:     req.Name,
		UserCode: claims.Code,
		Key:      uuid.NewString(),
	}).Error; err != nil {
		log.Error("新增客户端失败", zap.Error(err))
		return errors.New("操作失败")
	}
	return nil
}
