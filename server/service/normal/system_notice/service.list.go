package service

import (
	"server/model"
	"server/repository"
	"time"
)

type ListItem struct {
	Code       string `json:"code"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	IndexValue int    `json:"indexValue"`
	Date       string `json:"date"`
}

func (service *service) List() (list []ListItem) {
	db, _, _ := repository.Get("")
	var notices []model.SystemNotice
	db.Where("hidden = 2").Order("index_value asc").Order("id desc").Find(&notices)
	for _, notice := range notices {
		list = append(list, ListItem{
			Code:       notice.Code,
			Title:      notice.Title,
			Content:    notice.Content,
			IndexValue: notice.IndexValue,
			Date:       notice.UpdatedAt.Format(time.DateTime),
		})
	}
	return list
}
