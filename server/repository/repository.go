package repository

import (
	"go.uber.org/zap"
	"server/global"
	"server/pkg/memory"
	"server/repository/query"
)

func Get(domain string) (*query.Query, memory.Interface, *zap.Logger) {
	return query.Q, global.Cache, global.Logger
}
