package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
)

type EmailReq struct {
	Enable      string `binding:"required" json:"enable"` // 是否启用邮件服务
	NickName    string `json:"nickName"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
	ResetPwdTpl string `json:"resetPwdTpl"` // 重置密码邮件模板
}

func (service *service) Email(req EmailReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		_, _ = tx.SystemConfig.Where(tx.SystemConfig.Kind.Eq(model.SYSTEM_CONFIG_KIND_EMAIL)).Delete()
		if err := tx.SystemConfig.Create(model.GenerateSystemConfigEmail(
			req.Enable,
			req.NickName,
			req.Host,
			req.Port,
			req.User,
			req.Pwd,
			req.ResetPwdTpl,
		)...); err != nil {
			log.Error("修改系统Email配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetSystemConfigEmail(model.SystemConfigEmail{
			Enable:      req.Enable,
			NickName:    req.NickName,
			Host:        req.Host,
			Port:        req.Port,
			User:        req.User,
			Pwd:         req.Pwd,
			ResetPwdTpl: req.ResetPwdTpl,
		})
		return nil
	})
}
