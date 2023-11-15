package main

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"time"
)

/*
*
  - Description:
    go语言发送邮件的入门案例（使用gomail包）
  - @Author Hollis
  - @Create 2023-11-04 16:01
*/
func main() {
	m := gomail.NewMessage()
	m.SetHeader("From", "hollis_pu@163.com")
	m.SetHeader("To", "3047595798@qq.com", "hollis_pu@163.com")
	m.SetAddressHeader("Cc", "854828902@qq.com", "Hollis")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", fmt.Sprintf("This is a test email, from go language program, current time is %s", time.Now().Format("2006-01-02 15:04:05")))
	m.Attach("./gomail/test.png")

	// 设置邮件服务器信息
	// 注意，163 邮箱的 SMTP 服务器地址是 smtp.163.com，端口为 465，并且需要启用 SSL 连接，因此我们设置了 TLSConfig。

	d := gomail.NewDialer("smtp.163.com", 465, "hollis_pu", "JNZXGGXBJBXEAQUF") // 这个需要要对应的邮件去申请SMTP授权密码

	// 启用 SSL 连接（163邮箱需要SSL）
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	fmt.Println("Mail sent successfully!")
}
