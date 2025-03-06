package service

import (
	"errors"
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
	notice, _ := db.SystemNotice.Where(db.SystemNotice.Code.Eq(req.Code)).First()
	if notice == nil {
		return errors.New("通知不存在")
	}
	notice.Title = req.Title
	notice.Content = req.Content
	notice.Hidden = req.Hidden
	notice.IndexValue = req.IndexValue
	if err := db.SystemNotice.Save(notice); err != nil {
		return errors.New("操作失败")
	}
	return nil
}
