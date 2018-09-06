package xmail

import (
	"crypto/tls"
	cf "digicon/token_service/conf"
	"gopkg.in/gomail.v2"
	"strconv"
)

func NewGomail() (*gomail.Dialer, string) {
	// 读取配置
	host := cf.Cfg.MustValue("mail", "host")
	user := cf.Cfg.MustValue("mail", "user")
	portStr := cf.Cfg.MustValue("mail", "port")
	port, _ := strconv.ParseInt(portStr, 10, 64)
	password := cf.Cfg.MustValue("mail", "password")

	// new SMTP Dialer
	d := gomail.NewDialer(host, int(port), user, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d, user
}
