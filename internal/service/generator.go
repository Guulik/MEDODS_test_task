package service

import (
	"MEDODS-test/internal/domain/model"
	sl "MEDODS-test/internal/lib/logger/slog"
	"MEDODS-test/internal/util/refresher"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"log/slog"
	"math/rand"
	"net/http"
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
		expiresAt       time.Time
		err             error
	)

	accessToken, err = s.generateAccessToken(ipAddress)
	if err != nil {
		log.Error("failed to generate access token", sl.Err(err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	refreshTokenRaw, err = s.generateRefreshToken()
	if err != nil {
		log.Error("failed to generate refresh token", sl.Err(err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	refreshHash, err = refresher.Encrypt(refreshTokenRaw)
	if err != nil {
		log.Error("failed to encrypt refresh token", sl.Err(err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	expiresAt = time.Now().Add(s.cfg.Auth.RefreshTTL)
	if err = s.tokenModifier.Insert(ctx, userID, string(refreshHash), ipAddress, expiresAt); err != nil {
		log.Error("failed to save refresh token", sl.Err(err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenRaw,
	}, nil
}

func (s *Service) generateAccessToken(ipAddress string) (string, error) {
	claims := jwt.MapClaims{
		"ip":  ipAddress,
		"exp": time.Now().Add(s.cfg.Auth.AccessTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) generateRefreshToken() (string, error) {
	b := make([]byte, 32)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
