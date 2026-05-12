package mailer

import (
	"context"
	"log/slog"
	"net/smtp"
	"strings"

	domain "github.com/chanler/prosel/backend/internal/domain/subscribe"
	"github.com/chanler/prosel/backend/internal/infrastructure/config"
)

type SMTPMailer struct {
	cfg config.MailConfig
	log *slog.Logger
}

func NewSMTPMailer(cfg config.MailConfig, log *slog.Logger) *SMTPMailer {
	return &SMTPMailer{cfg: cfg, log: log}
}

func (m *SMTPMailer) Send(ctx context.Context, message domain.MailMessage) error {
	if !m.cfg.Enabled {
		if m.log != nil {
			m.log.Info("mail disabled", slog.String("to", message.To), slog.String("subject", message.Subject))
		}
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	from := m.cfg.From
	if strings.TrimSpace(from) == "" {
		from = m.cfg.Username
	}
	body := "From: " + from + "\r\n" + "To: " + message.To + "\r\n" + "Subject: " + message.Subject + "\r\n" + "Content-Type: text/plain; charset=UTF-8\r\n\r\n" + message.Body
	addr := m.cfg.Host + ":" + m.cfg.Port
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	if m.cfg.Username == "" {
		auth = nil
	}
	return smtp.SendMail(addr, auth, from, []string{message.To}, []byte(body))
}
