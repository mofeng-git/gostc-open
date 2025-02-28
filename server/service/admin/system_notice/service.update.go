package service

import (
	"errors"
	"server/model"
	"server/repository"
)

type UpdateReq struct {
	Code       string `binding:"required" json:"code" label:"编号"`
	Title      string `binding:"required" json:"title"`
	Content    string `binding:"required" json:"content"`
	Hidden     int    `binding:"required" json:"hidden"`
	IndexValue int    `json:"indexValue"`
}

func (service *service) Update(req UpdateReq) error {
	db, _, _ := repository.Get("")
	var notice model.SystemNotice
	if db.Where("code = ?", req.Code).First(&notice).RowsAffected == 0 {
		return errors.New("通知不存在")
	}
	notice.Title = req.Title
	notice.Content = req.Content
	notice.Hidden = req.Hidden
	notice.IndexValue = req.IndexValue
	if err := db.Save(&notice).Error; err != nil {
		return errors.New("操作失败")
	}
	return nil
}
