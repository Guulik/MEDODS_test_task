package service

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

func TestEmailService_SendWarning(t *testing.T) {
	tests := []struct {
		name    string
		address string
		ip      string
	}{
		{
			name:    "local",
			address: gofakeit.Email(),
			ip:      "999.561.888.460",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender, err := NewSMTPSender("your_auth_service@gg.com", "0.0.0.0", 1025)
			require.NoError(t, err)
			log := slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)

			s := &EmailService{
				log:        log,
				SMTPSender: sender,
			}
			err = s.SendWarning(context.Background(), tt.address, tt.ip)
			require.NoError(t, err)
		})
	}
}
