package main

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"gopkg.in/ini.v1"
	"html/template"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-05 11:44
 */

func main() {
	dialer, senderEmail, toLists, ccLists := initInfo() // 初始化dialer，并从配置文件中获取senderEmail，toLists，ccLists

	verificationCode := generateRandomCode()            // 生成随机验证码
	emailContent := renderHtmlContent(verificationCode) // 渲染html内容

	mail := addEmailContent(senderEmail, toLists, ccLists, emailContent) // 添加邮件内容

	if err := dialer.DialAndSend(mail); err != nil { // 发送邮件
		log.Fatalf("Failed to send email, err: %v", err)
	}
	log.Println("Email sent successfully!")
}

func initInfo() (dialer *gomail.Dialer, senderEmail string, toLists []string, ccLists []string) {
	// 初始化邮件发送配置
	cfg, err := ini.Load("./send_email/email_config.ini")
	if err != nil {
		log.Fatalf("Failed to load email config: %v", err)
	}

	// 读取发送邮箱信息
	senderEmail = cfg.Section("sender").Key("email").String()
	password := cfg.Section("sender").Key("password").String()
	smtpServer := cfg.Section("sender").Key("smtp_server").String()
	smtpPort := cfg.Section("sender").Key("smtp_port").MustInt(587)

	// 连接到 SMTP 服务器
	dialer = gomail.NewDialer(smtpServer, smtpPort, senderEmail, password)

	// 读取接收方和抄送方邮箱信息
	toListStr := cfg.Section("recipient").Key("to_list").String()
	toLists = strings.Split(toListStr, ",")
	ccListStr := cfg.Section("recipient").Key("cc_list").String()
	ccLists = strings.Split(ccListStr, ",")
	return dialer, senderEmail, toLists, ccLists
}

func generateRandomCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 指定随机种子
	minVal := 100000                                     // 最小的六位数
	maxVal := 999999                                     // 最大的六位数
	code := minVal + r.Intn(maxVal-minVal+1)
	return fmt.Sprintf("%06d", code) // 将随机数格式化为六位数的字符串
}

func renderHtmlContent(verificationCode string) string {
	// 读取 HTML 模板文件
	templateBytes, err := os.ReadFile("./send_email/注册模板.html")
	if err != nil {
		panic(err)
	}

	// 创建 HTML 模板
	tmpl, err := template.New("emailTemplate").Parse(string(templateBytes))
	if err != nil {
		panic(err)
	}

	// 创建缓冲区来存储模板渲染后的内容
	var renderedEmailContent string
	buffer := &bytes.Buffer{}

	// 将验证码插入到模板
	data := struct {
		VerificationCode string
		CurrentTime      string
	}{
		VerificationCode: verificationCode,
		CurrentTime:      time.Now().Format("2006-01-02 15:04:05"),
	}

	// 渲染模板
	err = tmpl.Execute(buffer, data)
	if err != nil {
		panic(err)
	}

	// 获取渲染后的 HTML 内容
	renderedEmailContent = buffer.String()
	return renderedEmailContent
}

func addEmailContent(senderEmail string, toLists []string, ccLists []string, emailContent string) *gomail.Message {
	// 添加邮件内容
	mail := gomail.NewMessage()
	mail.SetHeader("From", senderEmail)
	mail.SetHeader("To", toLists...)
	if len(ccLists[0]) > 0 {
		mail.SetHeader("Cc", ccLists...)
	}

	mail.SetHeader("Subject", "【Ecodego】邮箱验证码")
	mail.SetBody("text/html", emailContent)
	return mail
}
