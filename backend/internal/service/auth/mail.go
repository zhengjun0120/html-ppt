package auth

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// sendCodeMail 通过 SMTPS（TLS 直连，QQ 邮箱 465 端口）发送验证码邮件。
// 手写 SMTP 客户端而不是引第三方库：逻辑就几十行，且能看清邮件到底怎么发出去的。
func (s *Service) sendCodeMail(to, code string) error {
	if s.smtp.Host == "" || s.smtp.Username == "" || s.smtp.Password == "" || s.smtp.From == "" {
		return fmt.Errorf("smtp 未配置完整（config.yaml 的 smtp 段：username/password 需要填 QQ 邮箱和授权码）")
	}

	addr := fmt.Sprintf("%s:%d", s.smtp.Host, s.smtp.Port)
	d := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := tls.DialWithDialer(d, "tcp", addr, &tls.Config{ServerName: s.smtp.Host})
	if err != nil {
		return fmt.Errorf("连接 %s 失败: %w", addr, err)
	}
	defer conn.Close()

	cl, err := smtp.NewClient(conn, s.smtp.Host)
	if err != nil {
		return fmt.Errorf("SMTP 握手失败: %w", err)
	}
	defer cl.Close()

	if ok, _ := cl.Extension("AUTH"); ok {
		// PlainAuth 只肯在 TLS 连接上发凭据（且校验 host 名），这里恰好都满足
		if err := cl.Auth(smtp.PlainAuth("", s.smtp.Username, s.smtp.Password, s.smtp.Host)); err != nil {
			return fmt.Errorf("SMTP 认证失败（检查授权码）: %w", err)
		}
	}
	if err := cl.Mail(s.smtp.From); err != nil {
		return fmt.Errorf("SMTP MAIL FROM: %w", err)
	}
	if err := cl.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT TO: %w", err)
	}
	w, err := cl.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA: %w", err)
	}

	subject := mimeWord("html-ppt 注册验证码")
	body := fmt.Sprintf(
		"你的注册验证码是：%s\r\n\r\n验证码 10 分钟内有效。若非本人操作，请忽略本邮件。\r\n", code)

	msg := strings.Join([]string{
		"From: " + s.smtp.From,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: base64",
		"", // 头与正文之间的空行
		wrap76(base64.StdEncoding.EncodeToString([]byte(body))),
	}, "\r\n")
	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("写邮件正文: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("发送邮件: %w", err)
	}
	return cl.Quit()
}

// mimeWord 按 RFC 2047 编码 UTF-8 主题（=?UTF-8?B?...?=），否则中文标题在部分客户端乱码。
func mimeWord(s string) string {
	return fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(s)))
}

// wrap76 把 base64 按 76 字符折行——SMTP 规范要求单行不超过 998 字符，折行最稳。
func wrap76(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i += 76 {
		end := i + 76
		if end > len(s) {
			end = len(s)
		}
		b.WriteString(s[i:end])
		if end < len(s) {
			b.WriteString("\r\n")
		}
	}
	return b.String()
}
