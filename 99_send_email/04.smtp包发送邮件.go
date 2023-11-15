package main

import (
	"net/smtp"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-04 23:04
 */

func main() {
	// 定义 SMTP 服务器的地址和端口
	smtpServer := "smtp.163.com:465"

	// 发件人和收件人的电子邮件地址
	from := "ecodego@163.com"
	to := "hollis_pu@163.com"

	// 邮件服务器的身份验证信息
	auth := smtp.PlainAuth("", "ecodego@163.com", "YKSVKJKRLKWQBIXC", "smtp.163.com")

	// 电子邮件的主体内容
	subject := "Hello, World!"
	body := "This is the email body."

	// 构建邮件消息
	msg := "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	// 连接到 SMTP 服务器
	client, err := smtp.Dial(smtpServer)
	if err != nil {
		panic(err)
	}

	// 发起 TLS 连接（如果 SMTP 服务器需要）
	err = client.StartTLS(nil)
	if err != nil {
		panic(err)
	}

	// 登录到 SMTP 服务器
	err = client.Auth(auth)
	if err != nil {
		panic(err)
	}

	// 设置发件人
	err = client.Mail(from)
	if err != nil {
		panic(err)
	}

	// 设置收件人
	err = client.Rcpt(to)
	if err != nil {
		panic(err)
	}

	// 写入邮件内容
	wc, err := client.Data()
	if err != nil {
		panic(err)
	}
	defer wc.Close()

	_, err = wc.Write([]byte(msg))
	if err != nil {
		panic(err)
	}

	// 发送邮件并退出
	err = client.Quit()
	if err != nil {
		panic(err)
	}
}
