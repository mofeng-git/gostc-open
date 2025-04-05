package global

import (
	"go.uber.org/zap"
	"proxy/configs"
)

var Config *configs.Config
var Logger *zap.Logger
