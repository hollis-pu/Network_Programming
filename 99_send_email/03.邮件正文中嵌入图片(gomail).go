package main

import (
	"gopkg.in/gomail.v2"
	"gopkg.in/ini.v1"
	"log"
	"strings"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-04 22:32
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
	if ccListStr != "" {
		m.SetHeader("Cc", ccLists...)
	}

	// 设置邮件内容
	m.SetHeader("Subject", "邮件中嵌入图片")
	m.Embed("./gomail/test.png")
	m.SetBody("text/html", `
        <html>
            <body>
                <p>这是一封包含图片的邮件</p>
                <img src="cid:test.png" alt="Example Image" style="transform: scale(0.4);">
            </body>
        </html>
    `)
	m.Attach("./gomail/test.png")

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Failed to send email, err: %v", err)
	}
	log.Println("Email sent successfully!")
}
