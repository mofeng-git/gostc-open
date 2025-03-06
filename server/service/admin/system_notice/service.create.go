package service

import (
	"errors"
	"server/model"
	"server/repository"
)

type CreateReq struct {
	Title      string `binding:"required" json:"title"`
	Content    string `binding:"required" json:"content"`
	Hidden     int    `binding:"required" json:"hidden"`
	IndexValue int    `json:"indexValue"`
}

func (service *service) Create(req CreateReq) error {
	db, _, _ := repository.Get("")
	if err := db.SystemNotice.Create(&model.SystemNotice{
		Title:      req.Title,
		Content:    req.Content,
		Hidden:     req.Hidden,
		IndexValue: req.IndexValue,
	}); err != nil {
		return errors.New("操作失败")
	}
	return nil
}
