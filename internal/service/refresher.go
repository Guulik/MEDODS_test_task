package service

import (
	"MEDODS-test/internal/domain/model"
	sl "MEDODS-test/internal/lib/logger/slog"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"time"
)

func (s *Service) RefreshTokens(ctx context.Context, userID, refreshTokenRaw, ipAddress string) (*model.TokenPair, error) {
	log := s.log.With(slog.String("op", "Service.RefreshTokens"))

	var (
		token *model.RefreshTokenDB
		err   error
	)
	token, err = s.tokenProvider.Get(ctx, userID)
	if err != nil {
		log.Info("user have no refresh tokens", sl.Err(err))
		return nil, echo.NewHTTPError(http.StatusBadRequest, errors.New("user have no refresh tokens"))
	}

	if err = s.expired(token); err != nil {
		log.Info(fmt.Sprintf("token for user %s is expired: ", token.UserGuid))

		if err = s.tokenModifier.Delete(ctx, userID); err != nil {
			log.With("action:", "deleting expired").Error("failed to delete refresh token", sl.Err(err))
			return nil, echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to delete refresh token"))
		}
		return nil, echo.NewHTTPError(http.StatusBadRequest, errors.New("expired refresh token"))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(token.TokenHash), []byte(refreshTokenRaw)); err != nil {
		log.Info("failed to check refresh token", sl.Err(err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to decrypt refresh token"))
	}

	// TODO: отправка предупреждения на емейл
	/*
		if tokenData.IPAddress != ipAddress {
			mockEmail := "user@example.com"
			s.emailService.SendWarning(mockEmail, ipAddress)
		}*/

	if err = s.tokenModifier.Delete(ctx, userID); err != nil {
		log.Error("failed to delete refresh token", sl.Err(err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to delete refresh token"))
	}

	return s.GenerateTokens(ctx, userID, ipAddress)
}

func (s *Service) expired(token *model.RefreshTokenDB) error {
	if token.ExpiresAt.Unix() < time.Now().Unix() {
		return errors.New("token is expired")
	}
	return nil
}
