package service

import (
	"MEDODS-test/internal/configure"
	"MEDODS-test/internal/domain/model"
	sl "MEDODS-test/internal/lib/logger/slog"
	"MEDODS-test/internal/util/email"
	"context"
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
	"log/slog"
)

type SMTPSender struct {
	from string
	host string
	port int
}

func NewSMTPSender(from, host string, port int) (*SMTPSender, error) {
	if !email.IsEmailValid(from) {
		return nil, errors.New("invalid from email")
	}

	return &SMTPSender{from: from, host: host, port: port}, nil
}

type EmailService struct {
	log *slog.Logger

	SMTPSender *SMTPSender
}

func NewEmailService(
	cfg *configure.Config,
	log *slog.Logger,
) (*EmailService, error) {
	sender, err := NewSMTPSender("your_auth_service@gg.com", cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		return &EmailService{}, err
	}
	return &EmailService{
		log:        log,
		SMTPSender: sender,
	}, nil
}

func (s *EmailService) SendWarning(ctx context.Context, address string, newIp string) error {
	input := model.SendEmailInput{
		To:      address,
		Subject: "We detected new login attempt",
		Body:    fmt.Sprintf("Someone with ip: %s refreshed your tokens. If it's not you, think about it...", newIp),
	}
	err := s.SMTPSender.send(input)
	if err != nil {
		s.log.Error("op: EmailService.SendWarning", sl.Err(err))
		return err
	}
	return nil
}

func (s *SMTPSender) send(input model.SendEmailInput) error {
	if err := email.ValidateEmailInput(input); err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", s.from)
	msg.SetHeader("To", input.To)
	msg.SetHeader("Subject", input.Subject)
	msg.SetBody("text/html", input.Body)

	dialer := gomail.NewDialer(s.host, s.port, "", "")
	if err := dialer.DialAndSend(msg); err != nil {
		return errors.Wrap(err, "failed to sent email via smtp")
	}

	return nil
}
