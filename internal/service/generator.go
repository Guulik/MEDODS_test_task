package service

import (
	"MEDODS-test/internal/domain/model"
	sl "MEDODS-test/internal/lib/logger/slog"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

func (s *Service) GenerateTokens(ctx context.Context, userID, ipAddress string) (*model.TokenPair, error) {
	log := s.log.With(
		slog.String("op", "Service.GenerateTokens"),
	)

	var (
		accessToken     string
		refreshTokenRaw string
		refreshHash     []byte
		err             error
	)

	accessToken, err = s.generateAccessToken(ipAddress)
	if err != nil {
		log.Error("failed to generate access token", sl.Err(err))
		return nil, err
	}

	refreshTokenRaw, err = s.generateRefreshToken(ipAddress)
	if err != nil {
		log.Error("failed to generate refresh token", sl.Err(err))
		return nil, err
	}

	refreshHash, err = bcrypt.GenerateFromPassword([]byte(refreshTokenRaw), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to encrypt refresh token", sl.Err(err))
		return nil, err
	}

	if err = s.tokenModifier.Insert(ctx, userID, string(refreshHash), ipAddress); err != nil {
		log.Error("failed to save refresh token", sl.Err(err))
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenRaw,
	}, nil
}

func (s *Service) generateAccessToken(ipAddress string) (string, error) {
	claims := jwt.MapClaims{
		"ip":  ipAddress,
		"exp": time.Now().Add(s.cfg.AccessTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) generateRefreshToken(ipAddress string) (string, error) {
	claims := jwt.MapClaims{
		"ip":  ipAddress,
		"exp": time.Now().Add(s.cfg.RefreshTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(s.jwtSecret)
}
