package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

const (
	defaultAppName = "Local Marketplace"
	defaultTimeout = 30 * time.Second
	resendBaseURL  = "https://api.resend.com"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
	AppName() string
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

type ResendMailer struct {
	apiKey string
	from   string
	client *http.Client
	appName string
	logger *zap.Logger
}

func (m *ResendMailer) AppName() string { return m.appName }

func (m *ResendMailer) Send(ctx context.Context, to, subject, body string) error {
	payload := map[string]interface{}{
		"from":    m.from,
		"to":      []string{to},
		"subject": subject,
		"html":    body,
	}

	js, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("resend marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendBaseURL+"/emails", bytes.NewReader(js))
	if err != nil {
		return fmt.Errorf("resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("resend dial: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		var res struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(bodyBytes, &res)
		if res.Error != "" {
			return fmt.Errorf("resend: %s", res.Error)
		}
		return fmt.Errorf("resend: status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var res struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("resend decode: %w", err)
	}

	m.logger.Info("email sent", zap.String("to", to), zap.String("subject", subject), zap.String("id", res.ID))
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
	apiKey := os.Getenv("SMTP_PASSWORD")
	if apiKey == "" {
		apiKey = os.Getenv("RESEND_API_KEY")
	}

	from := os.Getenv("SMTP_FROM")
	appName := os.Getenv("APP_NAME")

	if from == "" {
		logger.Warn("SMTP_FROM not set, providers may reject the sender")
		from = "no-reply@localhost"
	}
	if appName == "" {
		appName = defaultAppName
	}

	if apiKey == "" {
		logger.Warn("RESEND_API_KEY not set, using console email transport")
		return NewConsoleMailer(logger, appName)
	}

	logger.Info(
		"Resend HTTP transport configured",
		zap.String("from", from),
		zap.String("timeout", defaultTimeout.String()),
	)

	return &ResendMailer{
		apiKey:  apiKey,
		from:    from,
		client:  &http.Client{Timeout: defaultTimeout},
		appName: appName,
		logger:  logger,
	}
}
