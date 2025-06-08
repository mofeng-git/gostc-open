package auth

import (
	"github.com/gin-gonic/gin"
	"server/pkg/bean"
	"server/router/middleware"
	"server/service/auth"
)

var svr = service.Service

func Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	result, err := svr.Login(c.ClientIP(), req)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.OkData(c, result)
}

func Register(c *gin.Context) {
	var req service.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	err := svr.Register(c.ClientIP(), req)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func LoginOtp(c *gin.Context) {
	var req service.LoginOtpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	result, err := svr.LoginOtp(c.ClientIP(), req)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.OkData(c, result)
}

func Captcha(c *gin.Context) {
	result, err := svr.Captcha(c.ClientIP())
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.OkData(c, result)
}

func Password(c *gin.Context) {
	var req service.PasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.Password(middleware.GetClaims(c), req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func Renew(c *gin.Context) {
	claims := middleware.GetClaims(c)
	token := middleware.GetToken(c)
	result, err := svr.Renew(token, claims)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.OkData(c, result)
}

func UserInfo(c *gin.Context) {
	claims := middleware.GetClaims(c)
	userInfo, err := svr.UserInfo(claims)
	if err != nil {
		bean.Response.AuthInvalid(c)
		return
	}
	bean.Response.OkData(c, userInfo)
}

func Checkin(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if err := svr.Checkin(claims); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func GenOtp(c *gin.Context) {
	otp, err := svr.GenOtp(middleware.GetClaims(c))
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.OkData(c, otp)
}

func OpenOtp(c *gin.Context) {
	var req service.OpenOtpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.OpenOtp(middleware.GetClaims(c), req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}
func CloseOtp(c *gin.Context) {
	if err := svr.CloseOtp(middleware.GetClaims(c)); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func BindEmail(c *gin.Context) {
	var req service.BindEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	claims := middleware.GetClaims(c)
	if err := svr.BindEmail(claims, req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}

func GenBindEmailCode(c *gin.Context) {
	var req service.GenBindEmailCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	claims := middleware.GetClaims(c)
	result, err := svr.GenBindEmailCode(claims, req)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.OkData(c, result)
}

func UnBindEmail(c *gin.Context) {
	claims := middleware.GetClaims(c)
	svr.UnBindEmail(claims)
	bean.Response.Ok(c)
}

func GenResetPwdEmailCode(c *gin.Context) {
	var req service.GenResetPwdEmailCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	result, err := svr.GenResetPwdEmailCode(req)
	if err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.OkData(c, result)
}

func ResetPwd(c *gin.Context) {
	var req service.ResetPwdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		bean.Response.Param(c, err)
		return
	}
	if err := svr.ResetPwd(req); err != nil {
		bean.Response.Fail(c, err.Error())
		return
	}
	bean.Response.Ok(c)
}
