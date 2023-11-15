package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"gopkg.in/ini.v1"
	"log"
	"strings"
	"time"
)

/*
*
  - Description:
    发送方邮箱，密码，接收方、抄送方邮箱信息都通过配置文件获取。
  - @Author Hollis
  - @Create 2023-11-04 17:22
*/
func main() {

	// 初始化邮件发送配置
	cfg, err := ini.Load("./gomail/email_config.ini")
	if err != nil {
		log.Fatalf("Failed to load email config: %v", err)
	}

	// 读取发送邮箱信息
	senderEmail := cfg.Section("sender").Key("email").String()
	password := cfg.Section("sender").Key("password").String()
	smtpServer := cfg.Section("sender").Key("smtp_server").String()
	smtpPort := cfg.Section("sender").Key("smtp_port").MustInt(587)

	// 连接到 SMTP 服务器
	d := gomail.NewDialer(smtpServer, smtpPort, senderEmail, password)
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}  // 是否跳过TLS证书验证

	// 读取接收方和抄送方邮箱信息
	toListStr := cfg.Section("recipient").Key("to_list").String()
	toLists := strings.Split(toListStr, ",")
	ccListStr := cfg.Section("recipient").Key("cc_list").String()
	ccLists := strings.Split(ccListStr, ",")

	// 设置发送方、接收方、抄送方邮箱
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", toLists...)
	if len(ccLists[0]) > 0 {
		m.SetHeader("Cc", ccLists...)
	}

	// 设置邮件内容
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", fmt.Sprintf("This is a test email, from go language program, current time is %s", time.Now().Format("2006-01-02 15:04:05")))
	m.Attach("./gomail/test.png")

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Failed to send email, err: %v", err)
	}
	log.Println("Email sent successfully!")
}
