package service

import (
	"MEDODS-test/internal/domain/model"
	sl "MEDODS-test/internal/lib/logger/slog"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func (s *Service) RefreshTokens(ctx context.Context, userID, refreshTokenRaw, ipAddress string) (*model.TokenPair, error) {
	log := s.log.With(slog.String("op", "Service.RefreshTokens"))

	var (
		tokenData *model.RefreshTokenDB
		err       error
	)
	tokenData, err = s.tokenProvider.Get(ctx, userID)
	if err != nil {
		log.Error("failed to get refresh token", sl.Err(err))
		return nil, errors.New("invalid or expired refresh token")
	}

	// TODO: хз пока
	if err = bcrypt.CompareHashAndPassword([]byte(tokenData.TokenHash), []byte(refreshTokenRaw)); err != nil {
		log.Error("failed to check refresh token", sl.Err(err))
		return nil, errors.New("invalid or expired refresh token")
	}

	// TODO: отправка предупреждения на емейл
	/*
		if tokenData.IPAddress != ipAddress {
			mockEmail := "user@example.com"
			s.emailService.SendWarning(mockEmail, ipAddress)
		}*/

	if err = s.tokenModifier.Delete(ctx, userID); err != nil {
		log.Error("failed to delete refresh token", sl.Err(err))
		return nil, err
	}

	return s.GenerateTokens(ctx, userID, ipAddress)
}
