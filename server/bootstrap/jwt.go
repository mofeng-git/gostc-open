package bootstrap

import (
	"server/global"
	"server/pkg/jwt"
)

func InitJwt() {
	global.Jwt = jwt.NewTool(global.Config.AuthKey)
}
