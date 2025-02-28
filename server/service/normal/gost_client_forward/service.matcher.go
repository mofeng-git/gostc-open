package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/service/gost_engine"
)

type MatcherReq struct {
	Code       string        `binding:"required" json:"code"`
	Enable     int           `json:"enable"`
	Matchers   []MatcherItem `json:"matchers"`
	TcpMatcher ItemMatcher   `json:"tcpMatcher"`
	SSHMatcher ItemMatcher   `json:"sshMatcher"`
}

type MatcherItem struct {
	Host       string `json:"host"`
	TargetIp   string `json:"targetIp"`
	TargetPort string `json:"targetPort"`
}

func (service *service) Matcher(claims jwt.Claims, req MatcherReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *gorm.DB) error {
		var user model.SystemUser
		if tx.Where("code = ?", claims.Code).First(&user).RowsAffected == 0 {
			return errors.New("用户错误")
		}

		var forward model.GostClientForward
		if tx.Where("code = ? AND user_code = ?", req.Code, user.Code).First(&forward).RowsAffected == 0 {
			return errors.New("操作失败")
		}

		forward.MatcherEnable = req.Enable
		var matchers []model.ForwardMatcher
		for _, matcher := range req.Matchers {
			matchers = append(matchers, model.ForwardMatcher{
				Host:       matcher.Host,
				TargetIp:   matcher.TargetIp,
				TargetPort: matcher.TargetPort,
			})
		}
		forward.SetMatcher(matchers)
		forward.SetTcpMatcher(req.TcpMatcher.TargetIp, req.TcpMatcher.TargetPort)
		forward.SetSSHMatcher(req.SSHMatcher.TargetIp, req.SSHMatcher.TargetPort)

		if err := tx.Save(&forward).Error; err != nil {
			log.Error("修改端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		gost_engine.ClientForwardConfig(tx, forward.Code)
		return nil
	})
}
