package service

import (
	"server/model"
	"server/pkg/bean"
	"server/repository"
	"time"
)

type PageReq struct {
	bean.PageParam
}

type Item struct {
	Code       string `json:"code"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Hidden     int    `json:"hidden"`
	IndexValue int    `json:"indexValue"`
	Date       string `json:"date"`
}

func (service *service) Page(req PageReq) (list []Item, total int64) {
	db, _, _ := repository.Get("")
	var notices []model.SystemNotice
	var where = db
	db.Where(where).Model(&notices).Count(&total)
	db.Where(where).Order("index_value asc").Order("id desc").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&notices)
	for _, notice := range notices {
		list = append(list, Item{
			Code:       notice.Code,
			Title:      notice.Title,
			Content:    notice.Content,
			Hidden:     notice.Hidden,
			IndexValue: notice.IndexValue,
			Date:       notice.UpdatedAt.Format(time.DateTime),
		})
	}
	return list, total
}
