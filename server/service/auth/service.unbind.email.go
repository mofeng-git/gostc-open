package service

import (
	"server/pkg/jwt"
	"server/repository"
)

func (service *service) UnBindEmail(claims jwt.Claims) {
	db, _, _ := repository.Get("")
	_, _ = db.SystemUserEmail.Where(db.SystemUserEmail.UserCode.Eq(claims.Code)).Delete()
}
