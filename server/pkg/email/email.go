package email

import (
	"gopkg.in/gomail.v2"
)

type Config struct {
	Host string // 邮件地址
	Port int    // 端口
	User string // 用户
	Pwd  string // 密码
}

type Body struct {
	Mails    []string // 接收者
	NickName string   // 发件人昵称
	Title    string   // 标题
	Body     string   // 内容
}

func Send(config Config, data Body) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(config.User, data.NickName))
	m.SetHeader("To", data.Mails...)   // 发送给多个用户
	m.SetHeader("Subject", data.Title) // 设置邮件主题
	m.SetBody("text/html", data.Body)  // 设置邮件正文
	d := gomail.NewDialer(config.Host, config.Port, config.User, config.Pwd)
	err := d.DialAndSend(m) // 发送邮件
	return err
}
