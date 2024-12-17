package service

import (
	"MEDODS-test/internal/configure"
	"MEDODS-test/internal/domain/model"
	"MEDODS-test/internal/service/mocks"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestService_generateRefreshToken(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
		{
			name: "2",
		},
		{
			name: "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TokenService{
				cfg:           nil,
				log:           nil,
				emailNotifier: nil,
				tokenModifier: nil,
				tokenProvider: nil,
			}
			refreshToken, err := s.generateRefreshToken()
			require.NoError(t, err)

			decodedToken, err := base64.StdEncoding.DecodeString(refreshToken)
			require.NoError(t, err)

			encoded := base64.StdEncoding.EncodeToString(decodedToken)
			require.Equal(t, refreshToken, encoded)
		})
	}
}

func TestTokenService_generateAccessToken(t *testing.T) {
	tests := []struct {
		name      string
		ipAddress string
	}{
		{
			name:      "random ip",
			ipAddress: "5.21.77.23",
		},
		{
			name:      "null ip",
			ipAddress: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TokenService{
				cfg:           &configure.Config{Env: "local", Auth: configure.Auth{AccessTTL: 5 * time.Minute}},
				log:           nil,
				emailNotifier: nil,
				tokenModifier: nil,
				tokenProvider: nil,
			}
			tokenString, err := s.generateAccessToken(tt.ipAddress)
			require.NoError(t, err)
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("default"), nil
			})
			require.NoError(t, err)
			require.True(t, token.Valid)

			claims, ok := token.Claims.(jwt.MapClaims)
			require.True(t, ok)
			require.Equal(t, tt.ipAddress, claims["ip"])
		})
	}
}

func TestTokenService_expired(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "1 millisec",
			duration: 1 * time.Millisecond,
		},
		{
			name:     "500 millisec",
			duration: 500 * time.Millisecond,
		},
		{
			name:     "900 millisec",
			duration: 900 * time.Millisecond,
		},
	}
	sleep := 100 * time.Millisecond
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TokenService{
				cfg:           nil,
				log:           nil,
				emailNotifier: nil,
				tokenModifier: nil,
				tokenProvider: nil,
			}
			refresh, err := s.generateRefreshToken()
			require.NoError(t, err)

			tokenDB := &model.RefreshTokenDB{
				UserGuid:  "hw",
				TokenHash: refresh,
				IPAddress: "666",
				ExpiresAt: time.Now().Add(tt.duration),
			}
			time.Sleep(sleep)

			err = s.expired(tokenDB)

			if tt.duration < sleep {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTokenService_GenerateTokens(t *testing.T) {
	type service struct {
		cfg *configure.Config
		log *slog.Logger
	}
	type args struct {
		userID    string
		ipAddress string
	}
	tests := []struct {
		name    string
		service service
		args    args
		want    *model.TokenPair
		wantErr bool
	}{
		{
			name: "a",
			service: service{
				cfg: &configure.Config{Env: "local",
					Auth: configure.Auth{AccessTTL: 5 * time.Minute, RefreshTTL: 10 * time.Minute}},
				log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
			},
			args: args{
				"u78", "62.112.77.99",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			tokenModifier := mocks.NewTokenModifier(t)

			tokenProvider := mocks.NewTokenProvider(t)
			s := &TokenService{
				cfg:           tt.service.cfg,
				log:           tt.service.log,
				emailNotifier: nil,
				tokenModifier: tokenModifier,
				tokenProvider: tokenProvider,
			}
			tokenModifier.On("Insert", ctx, tt.args.userID, mock.Anything, tt.args.ipAddress, mock.Anything).
				Return(nil).Once()

			tokenProvider.On("Get", ctx, tt.args.userID).
				Return(nil, nil).Once()
			_, err := s.GenerateTokens(ctx, tt.args.userID, tt.args.ipAddress)
			require.NoError(t, err)
		})
	}
}
