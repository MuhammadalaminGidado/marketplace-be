package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"

	"go.uber.org/zap"
)

const defaultAppName = "Local Marketplace"

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
	AppName() string
}

type Options struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	AppName  string
}

type ConsoleMailer struct {
	logger  *zap.Logger
	appName string
}

func NewConsoleMailer(logger *zap.Logger, appName string) *ConsoleMailer {
	if appName == "" {
		appName = defaultAppName
	}
	return &ConsoleMailer{logger: logger, appName: appName}
}

func (m *ConsoleMailer) AppName() string { return m.appName }

func (m *ConsoleMailer) Send(_ context.Context, to, subject, body string) error {
	m.logger.Info(
		"email (console transport)",
		zap.String("to", to),
		zap.String("subject", subject),
		zap.String("body", body),
	)
	return nil
}

type SMTPMailer struct {
	host    string
	port    string
	from    string
	auth    smtp.Auth
	appName string
	logger  *zap.Logger
}

func (m *SMTPMailer) AppName() string { return m.appName }

func (m *SMTPMailer) Send(ctx context.Context, to, subject, body string) error {
	addr := net.JoinHostPort(m.host, m.port)

	msg := fmt.Appendf(nil,
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, to, subject, body,
	)

	done := make(chan error, 1)
	go func() {
		done <- sendMail(ctx, m.host, addr, m.auth, m.from, []string{to}, msg)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func sendMail(_ context.Context, host, addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Hello("localhost"); err != nil {
		return err
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(&tls.Config{ServerName: host}); err != nil {
			return err
		}
	}

	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(msg); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}

	return c.Quit()
}

func NewFromEnv(logger *zap.Logger) Mailer {
	options := Options{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		AppName:  os.Getenv("APP_NAME"),
	}

	if options.From == "" {
		options.From = "no-reply@localhost"
	}
	if options.AppName == "" {
		options.AppName = defaultAppName
	}

	if options.Host == "" {
		logger.Warn("SMTP_HOST not set, using console email transport")
		return NewConsoleMailer(logger, options.AppName)
	}

	if options.Port == "" {
		options.Port = "587"
	}

	var auth smtp.Auth
	if options.Username != "" {
		auth = smtp.PlainAuth("", options.Username, options.Password, options.Host)
	}

	return &SMTPMailer{
		host:    options.Host,
		port:    options.Port,
		from:    options.From,
		auth:    auth,
		appName: options.AppName,
		logger:  logger,
	}
}
