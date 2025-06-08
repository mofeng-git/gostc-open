package global

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"server/configs"
	"server/pkg/jwt"
	"server/pkg/memory"
	"server/pkg/orm"
)

var (
	Logger *zap.Logger
	Jwt    *jwt.Tool
	Cache  memory.Interface
	Config configs.Config
	DB     orm.Interface
	Cron   *cron.Cron
)
