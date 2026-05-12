package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/eun1e/sublinkE-plugins"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// EmailPlugin 实现了 API 事件通知插件
type EmailPlugin struct {
	smtpServer   string
	smtpPort     int
	smtpUser     string
	smtpPassword string
	toEmail      string
}

// 导出获取插件实例的函数
func GetPlugin() plugins.Plugin {
	return &EmailPlugin{}
}

func (p *EmailPlugin) Name() string {
	return "EmailNotificationPlugin"
}

func (p *EmailPlugin) Version() string {
	return "1.0.0"
}

func (p *EmailPlugin) Description() string {
	return "在尝试登录后发送电子邮件通知"
}

func (p *EmailPlugin) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"smtpServer":   "example.smtp.com",
		"smtpPort":     465,
		"smtpUser":     "sender@example.com",
		"smtpPassword": "password",
		"toEmail":      "recevier@example.com",
	}
}

// SetConfig 设置插件配置
func (p *EmailPlugin) SetConfig(config map[string]interface{}) {
	// 安全地获取配置，提供默认值以防止缺失或类型错误
	if server, ok := config["smtpServer"].(string); ok {
		p.smtpServer = server
	} else {
		p.smtpServer = "example.smtp.com"
		fmt.Println("警告: SMTP服务器配置错误，使用默认值")
	}

	if port, ok := config["smtpPort"].(float64); ok { // JSON解析后数字类型通常是float64
		p.smtpPort = int(port)
	} else if port, ok := config["smtpPort"].(int); ok {
		p.smtpPort = port
	} else {
		p.smtpPort = 465
		fmt.Println("警告: SMTP端口配置错误，使用默认值")
	}

	if user, ok := config["smtpUser"].(string); ok {
		p.smtpUser = user
	} else {
		p.smtpUser = "sender@example.com"
		fmt.Println("警告: SMTP用户配置错误，使用默认值")
	}

	if password, ok := config["smtpPassword"].(string); ok {
		p.smtpPassword = password
	} else {
		p.smtpPassword = "password"
		fmt.Println("警告: SMTP密码配置错误，使用默认值")
	}

	if email, ok := config["toEmail"].(string); ok {
		p.toEmail = email
	} else {
		p.toEmail = "recevier@example.com"
		fmt.Println("警告: 收件人邮箱配置错误，使用默认值")
	}
}

func (p *EmailPlugin) Init() error {
	// 验证关键配置
	if p.smtpServer == "" {
		return fmt.Errorf("SMTP服务器地址不能为空")
	}
	if p.smtpPort <= 0 {
		return fmt.Errorf("SMTP端口无效")
	}
	if p.smtpUser == "" || p.smtpPassword == "" {
		return fmt.Errorf("SMTP用户名或密码不能为空")
	}
	if p.toEmail == "" {
		return fmt.Errorf("收件人邮箱不能为空")
	}

	fmt.Printf("邮件通知插件初始化成功: SMTP=%s:%d, 用户=%s, 收件人=%s\n",
		p.smtpServer, p.smtpPort, p.smtpUser, p.toEmail)
	return nil
}

func (p *EmailPlugin) Close() error {
	return nil
}

func (p *EmailPlugin) OnAPIEvent(ctx *gin.Context, event plugins.EventType, path string, statusCode int, requestBody interface{}, responseBody interface{}) error {

	// 根据事件类型和路径决定是否发送邮件
	if event == plugins.EventAPISuccess && path == "/api/v1/auth/login" {
		// 用户登录成功
		username := p.parseUsernameFromRequest(requestBody)

		subject := "用户登录通知"
		body := fmt.Sprintf(`🎉 用户登录成功通知

		👤 用户信息:
		━━━━━━━━━━━━━━━━━━━━━━━━
		🏷️  用户名: %s
		🕐  登录时间: %s
		🌐  客户端IP: %s
		🔑  访问令牌: %s
		━━━━━━━━━━━━━━━━━━━━━━━━
		
		系统已为用户生成新的访问令牌。`,
			username,
			time.Now().Format("2006-01-02 15:04:05"),
			ctx.ClientIP(),
			p.parseTokenFromResponse(responseBody))

		return p.sendEmail(subject, body)
	} else if event == plugins.EventAPIError && path == "/api/v1/auth/login" {
		username := p.parseUsernameFromRequest(requestBody)

		subject := "用户登录失败通知"
		body := fmt.Sprintf(`⚠️ 用户登录失败通知

		🚨 安全警告:
		━━━━━━━━━━━━━━━━━━━━━━━━
		👤  尝试用户名: %s
		🕐  尝试时间: %s
		🌐  客户端IP: %s
		📋  失败原因: 登录验证失败
		━━━━━━━━━━━━━━━━━━━━━━━━
		
		请检查是否为恶意登录尝试。`,
			username,
			time.Now().Format("2006-01-02 15:04:05"),
			ctx.ClientIP())

		return p.sendEmail(subject, body)
	}

	return nil
}

func (p *EmailPlugin) InterestedAPIs() []string {
	return []string{
		"/api/v1/auth/login",
		"/api/v1/nodes/add",
	}
}

func (p *EmailPlugin) InterestedEvents() []plugins.EventType {
	return []plugins.EventType{
		plugins.EventAPISuccess,
		plugins.EventAPIError,
	}
}

func (p *EmailPlugin) sendEmail(subject, body string) error {
	addr := fmt.Sprintf("%s:%d", p.smtpServer, p.smtpPort)

	// 建立 TLS 连接（适用于端口 465）
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // 如果证书有问题可以跳过验证（生产建议改为 false）
		ServerName:         p.smtpServer,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %w", err)
	}

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, p.smtpServer)
	if err != nil {
		return fmt.Errorf("SMTP 客户端创建失败: %w", err)
	}

	// 登录认证
	auth := smtp.PlainAuth("", p.smtpUser, p.smtpPassword, p.smtpServer)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP 登录失败: %w", err)
	}

	// 设置发件人和收件人
	if err := client.Mail(p.smtpUser); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}
	if err := client.Rcpt(p.toEmail); err != nil {
		return fmt.Errorf("设置收件人失败: %w", err)
	}

	// 写入邮件内容
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("写入邮件失败: %w", err)
	}

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: blog@eunie.online\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n"+
		"%s", p.toEmail, subject, body))

	if _, err = writer.Write(message); err != nil {
		return fmt.Errorf("写入内容失败: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("关闭邮件体失败: %w", err)
	}

	return client.Quit()
}

// 解析请求体中的用户名（multipart/form-data格式）
func (p *EmailPlugin) parseUsernameFromRequest(requestBody interface{}) string {
	if requestBody == nil {
		return "未知用户"
	}

	// 将interface{}转换为字符串
	requestStr, ok := requestBody.(string)
	if !ok {
		return "未知用户"
	}

	// 使用正则表达式解析multipart/form-data中的username字段
	re := regexp.MustCompile(`name="username".*?\r\n\r\n(.*?)\r\n`)
	matches := re.FindStringSubmatch(requestStr)

	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return "未知用户"
}

// 解析响应体中的访问令牌（JSON格式）
func (p *EmailPlugin) parseTokenFromResponse(responseBody interface{}) string {
	if responseBody == nil {
		return "未知令牌"
	}

	// 将interface{}转换为字符串
	responseStr, ok := responseBody.(string)
	if !ok {
		return "未知令牌"
	}

	// 解析JSON响应
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return "解析失败"
	}

	// 提取访问令牌
	if data, exists := response["data"]; exists {
		if dataMap, ok := data.(map[string]interface{}); ok {
			if token, exists := dataMap["accessToken"]; exists {
				if tokenStr, ok := token.(string); ok {
					return tokenStr
				}
			}
		}
	}

	return "未找到令牌"
}
