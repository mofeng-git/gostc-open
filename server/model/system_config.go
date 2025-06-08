package model

import (
	"strconv"
	"strings"
	"time"
)

const (
	SYSTEM_CONFIG_KIND_BASE  = "SystemConfigBase"  // 基础配置
	SYSTEM_CONFIG_KIND_GOST  = "SystemConfigGost"  // Gost配置
	SYSTEM_CONFIG_KIND_EMAIL = "SystemConfigEmail" // 邮件配置
)

type SystemConfig struct {
	Kind  string `gorm:"column:kind;size:100;uniqueIndex:system_config_uidx"`
	Key   string `gorm:"column:key;size:100;uniqueIndex:system_config_uidx"`
	Value string `gorm:"column:value"`
}

type SystemConfigBase struct {
	Title        string `json:"title"`
	Favicon      string `json:"favicon"`
	BaseUrl      string `json:"baseUrl"`
	ApiKey       string `json:"apiKey"`
	Register     string `json:"register"`
	CheckIn      string `json:"checkIn"`
	CheckInStart int    `json:"checkInStart"`
	CheckInEnd   int    `json:"checkInEnd"`
}

func GetSystemConfigBase(list []*SystemConfig) (result SystemConfigBase) {
	for _, item := range list {
		switch item.Key {
		case "title":
			result.Title = item.Value
		case "favicon":
			result.Favicon = item.Value
		case "baseUrl":
			result.BaseUrl = item.Value
		case "apiKey":
			result.ApiKey = item.Value
		case "register":
			result.Register = item.Value
		case "checkIn":
			result.CheckIn = item.Value
		case "checkInStart":
			result.CheckInStart, _ = strconv.Atoi(item.Value)
		case "checkInEnd":
			result.CheckInEnd, _ = strconv.Atoi(item.Value)
		}
	}
	return result
}

func GenerateSystemConfigBase(title, favicon, baseUrl, apiKey, register, checkIn string, checkInStart, checkInEnd int) []*SystemConfig {
	return []*SystemConfig{
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "title", Value: title},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "favicon", Value: favicon},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "baseUrl", Value: baseUrl},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "apiKey", Value: apiKey},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "register", Value: register},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "checkIn", Value: checkIn},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "checkInStart", Value: strconv.Itoa(checkInStart)},
		{Kind: SYSTEM_CONFIG_KIND_BASE, Key: "checkInEnd", Value: strconv.Itoa(checkInEnd)},
	}
}

type SystemConfigGost struct {
	Version string `json:"version"`
	Logger  string `json:"logger"`

	// 功能菜单
	FuncWeb     string `json:"funcWeb"`
	FuncForward string `json:"funcForward"`
	FuncTunnel  string `json:"funcTunnel"`
	FuncP2P     string `json:"funcP2P"`
	FuncProxy   string `json:"funcProxy"`
	FuncNode    string `json:"funcNode"`
}

func GenerateSystemConfigGost(version string, logger string, funcWeb, funcForward, funcTunnel, funcP2P, funcProxy, funcNode string) []*SystemConfig {
	return []*SystemConfig{
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "version", Value: version},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "logger", Value: logger},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcWeb", Value: funcWeb},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcForward", Value: funcForward},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcTunnel", Value: funcTunnel},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcP2P", Value: funcP2P},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcProxy", Value: funcProxy},
		{Kind: SYSTEM_CONFIG_KIND_GOST, Key: "funcNode", Value: funcNode},
	}
}

func GetSystemConfigGost(list []*SystemConfig) (result SystemConfigGost) {
	for _, item := range list {
		switch item.Key {
		case "version":
			result.Version = item.Value
		case "logger":
			result.Logger = item.Value
		case "funcWeb":
			result.FuncWeb = item.Value
		case "funcForward":
			result.FuncForward = item.Value
		case "funcTunnel":
			result.FuncTunnel = item.Value
		case "funcP2P":
			result.FuncP2P = item.Value
		case "funcProxy":
			result.FuncProxy = item.Value
		case "funcNode":
			result.FuncNode = item.Value
		}
	}
	return result
}

type SystemConfigEmail struct {
	Enable   string `json:"enable"` // 是否启用邮件服务
	NickName string `json:"nickName"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Pwd      string `json:"pwd"`

	ResetPwdTpl string `json:"resetPwdTpl"` // 重置密码邮件模板
}

func (cfg SystemConfigEmail) GenerateUserBindEmailTpl(code string, now time.Time) string {
	tpl := "您正在进行绑定邮箱操作，请勿将验证码告诉他人，验证码：{{CODE}}(5分钟有效)，发件时间：{{DATETIME}}"
	result := strings.ReplaceAll(tpl, "{{CODE}}", code)
	result = strings.ReplaceAll(result, "{{DATETIME}}", now.Format(time.DateTime))
	return result
}

func (cfg SystemConfigEmail) GenerateResetPwdTpl(code string, now time.Time) string {
	result := strings.ReplaceAll(cfg.ResetPwdTpl, "{{CODE}}", code)
	result = strings.ReplaceAll(result, "{{DATETIME}}", now.Format(time.DateTime))
	return result
}

func (cfg SystemConfigEmail) GenerateResetPwdResultTpl(account, password string, now time.Time) string {
	tpl := "重置密码成功，您的账号：{{ACCOUNT}}，您的新密码：{{PASSWORD}}，发件时间：{{DATETIME}}"
	result := strings.ReplaceAll(tpl, "{{PASSWORD}}", password)
	result = strings.ReplaceAll(result, "{{ACCOUNT}}", account)
	result = strings.ReplaceAll(result, "{{DATETIME}}", now.Format(time.DateTime))
	return result
}

func GenerateSystemConfigEmail(enable, nickName, host, port, user, pwd, resetPwdTpl string) []*SystemConfig {
	return []*SystemConfig{
		{Kind: SYSTEM_CONFIG_KIND_EMAIL, Key: "enable", Value: enable},
		{Kind: SYSTEM_CONFIG_KIND_EMAIL, Key: "nickName", Value: nickName},
		{Kind: SYSTEM_CONFIG_KIND_EMAIL, Key: "host", Value: host},
		{Kind: SYSTEM_CONFIG_KIND_EMAIL, Key: "port", Value: port},
		{Kind: SYSTEM_CONFIG_KIND_EMAIL, Key: "user", Value: user},
		{Kind: SYSTEM_CONFIG_KIND_EMAIL, Key: "pwd", Value: pwd},
		{Kind: SYSTEM_CONFIG_KIND_EMAIL, Key: "resetPwdTpl", Value: resetPwdTpl},
	}
}

func GetSystemConfigEmail(list []*SystemConfig) (result SystemConfigEmail) {
	for _, item := range list {
		switch item.Key {
		case "enable":
			result.Enable = item.Value
		case "nickName":
			result.NickName = item.Value
		case "host":
			result.Host = item.Value
		case "port":
			result.Port = item.Value
		case "user":
			result.User = item.Value
		case "pwd":
			result.Pwd = item.Value
		case "resetPwdTpl":
			result.ResetPwdTpl = item.Value
		}
	}
	return result
}
