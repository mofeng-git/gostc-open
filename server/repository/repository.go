package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/global"
	"server/pkg/memory"
)

func Get(domain string) (*gorm.DB, memory.Interface, *zap.Logger) {
	return global.DB.GetDB(), global.Cache, global.Logger
}
