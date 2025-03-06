package service

import (
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
	notices, _ := db.SystemNotice.Where(db.SystemNotice.Hidden.Eq(2)).Order(db.SystemNotice.IndexValue.Asc(), db.SystemNotice.Id.Desc()).Find()
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
