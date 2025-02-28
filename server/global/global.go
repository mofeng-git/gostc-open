package global

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"server/configs"
	"server/pkg/jwt"
	"server/pkg/memory"
	"server/pkg/orm"
)

var Logger *zap.Logger

var Jwt *jwt.Tool

var Cache memory.Interface

var Config configs.Config

var DB orm.Interface

var Cron *cron.Cron
