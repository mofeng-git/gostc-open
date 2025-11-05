package service

import (
    "errors"
    "go.uber.org/zap"
    "server/model"
    "server/repository"
    "server/repository/cache"
    "server/repository/query"
)

type HomeReq struct {
    HomeEnable string `binding:"required" json:"homeEnable"` // "1"=启用；其它值禁用并跳转 /login
    HomeTpl    string `json:"homeTpl"`
}

func (service *service) Home(req HomeReq) error {
    db, _, log := repository.Get("")
    return db.Transaction(func(tx *query.Query) error {
        _, _ = tx.SystemConfig.Where(tx.SystemConfig.Kind.Eq(model.SYSTEM_CONFIG_KIND_HOME)).Delete()
        if err := tx.SystemConfig.Create(model.GenerateSystemConfigHome(
            req.HomeEnable,
            req.HomeTpl,
        )...); err != nil {
            log.Error("修改系统首页配置失败", zap.Error(err))
            return errors.New("操作失败")
        }
        cache.SetSystemConfigHome(model.SystemConfigHome{
            HomeEnable: req.HomeEnable,
            HomeTpl:    req.HomeTpl,
        })
        return nil
    })
}

