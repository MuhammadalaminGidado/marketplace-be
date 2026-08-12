package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"time"

	"go.uber.org/zap"
)

const (
	defaultAppName = "Local Marketplace"
	defaultTimeout = 30 * time.Second
)

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
	host        string
	port        string
	from        string
	auth        smtp.Auth
	appName     string
	timeout     time.Duration
	implicitTLS bool
	logger      *zap.Logger
}

func (m *SMTPMailer) AppName() string { return m.appName }

func (m *SMTPMailer) Send(ctx context.Context, to, subject, body string) error {
	addr := net.JoinHostPort(m.host, m.port)

	msg := fmt.Appendf(nil,
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, to, subject, body,
	)

	deadline := time.Now().Add(m.timeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}

	if err := sendMail(ctx, addr, m.host, m.from, []string{to}, msg, m.auth, m.implicitTLS, deadline); err != nil {
		return err
	}

	m.logger.Info("email sent", zap.String("to", to), zap.String("subject", subject))
	return nil
}

func sendMail(
	ctx context.Context,
	addr, host, from string,
	to []string,
	msg []byte,
	a smtp.Auth,
	implicitTLS bool,
	deadline time.Time,
) error {

	d := net.Dialer{Timeout: 15 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}

	if implicitTLS {
		tlsConn := tls.Client(conn, &tls.Config{ServerName: host})
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			conn.Close()
			return fmt.Errorf("smtp tls handshake: %w", err)
		}
		conn = tlsConn
	}

	if err := conn.SetDeadline(deadline); err != nil {
		conn.Close()
		return fmt.Errorf("smtp set deadline: %w", err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("smtp greeting: %w", err)
	}
	defer c.Close()

	if err := c.Hello("localhost"); err != nil {
		return fmt.Errorf("smtp hello: %w", err)
	}

	if !implicitTLS {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(&tls.Config{ServerName: host}); err != nil {
				return fmt.Errorf("smtp starttls: %w", err)
			}
		}
	}

	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err := c.Auth(a); err != nil {
				return fmt.Errorf("smtp auth: %w", err)
			}
		}
	}

	if err := c.Mail(from); err != nil {
		return fmt.Errorf("smtp mail: %w", err)
	}
	for _, recipient := range to {
		if err := c.Rcpt(recipient); err != nil {
			return fmt.Errorf("smtp rcpt: %w", err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp data write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp data close: %w", err)
	}

	_ = c.Quit()
	return nil
}

func SendTimeout() time.Duration {
	if raw := os.Getenv("SMTP_TIMEOUT"); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil && d > 0 {
			return d
		}
	}
	return defaultTimeout
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
		logger.Warn("SMTP_FROM not set, providers may reject the sender")
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

	implicitTLS := options.Port == "465"
	if os.Getenv("SMTP_IMPLICIT_TLS") == "true" {
		implicitTLS = true
	}

	var auth smtp.Auth
	if options.Username != "" {
		auth = smtp.PlainAuth("", options.Username, options.Password, options.Host)
	}

	logger.Info(
		"SMTP transport configured",
		zap.String("host", options.Host),
		zap.String("port", options.Port),
		zap.Bool("implicit_tls", implicitTLS),
		zap.String("from", options.From),
		zap.Duration("timeout", SendTimeout()),
	)

	return &SMTPMailer{
		host:        options.Host,
		port:        options.Port,
		from:        options.From,
		auth:        auth,
		appName:     options.AppName,
		timeout:     SendTimeout(),
		implicitTLS: implicitTLS,
		logger:      logger,
	}
}
