package service

import (
	"MEDODS-test/internal/configure"
	"MEDODS-test/internal/domain/model"
	sl "MEDODS-test/internal/lib/logger/slog"
	"MEDODS-test/internal/util/jwtReader"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=TokenProvider
type TokenProvider interface {
	Get(ctx context.Context, userID string) (*model.RefreshTokenDB, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=TokenModifier
type TokenModifier interface {
	Insert(ctx context.Context, userID string, tokenHash string, ipAddress string, expiresAt time.Time) error
	Delete(ctx context.Context, userID string) error
}

type TokenService struct {
	cfg *configure.Config
	log *slog.Logger

	emailNotifier EmailNotifier
	tokenModifier TokenModifier
	tokenProvider TokenProvider
}

func NewTokenService(
	cfg *configure.Config,
	log *slog.Logger,
	emailNotifier EmailNotifier,
	tokenModifier TokenModifier,
	tokenProvider TokenProvider,
) *TokenService {
	return &TokenService{
		cfg:           cfg,
		log:           log,
		emailNotifier: emailNotifier,
		tokenModifier: tokenModifier,
		tokenProvider: tokenProvider,
	}
}

func (s *TokenService) GenerateTokens(ctx context.Context, userID, ipAddress string) (*model.TokenPair, error) {
	log := s.log.With(
		slog.String("op", "Service.GenerateTokens"),
	)

	var (
		accessToken     string
		refreshTokenRaw string
		refreshHash     []byte
		existedToken    *model.RefreshTokenDB
		expiresAt       time.Time
		err             error
	)

	existedToken, err = s.tokenProvider.Get(ctx, userID)
	if err != nil {
		log.Info("user have no refresh tokens", sl.Err(err))
	}
	if existedToken != nil {
		log.Warn("user tried to generate multiple pairs")
		return nil, echo.NewHTTPError(http.StatusBadRequest, errors.New("user already have access token"))
	}

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

	refreshHash, err = bcrypt.GenerateFromPassword([]byte(refreshTokenRaw), bcrypt.DefaultCost)
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

func (s *TokenService) RefreshTokens(ctx context.Context, userID, refreshTokenRaw, ipAddress string) (*model.TokenPair, error) {
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
		return nil, echo.NewHTTPError(http.StatusBadRequest, errors.New("this token invalid"))
	}

	if token.IPAddress != ipAddress {
		mockEmail := gofakeit.Email()
		err = s.emailNotifier.SendWarning(ctx, mockEmail, ipAddress)
		if err != nil {
			log.Error("failed to send email warning", sl.Err(err))
		}
	}

	if err = s.tokenModifier.Delete(ctx, userID); err != nil {
		log.Error("failed to delete refresh token", sl.Err(err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to delete refresh token"))
	}

	return s.GenerateTokens(ctx, userID, ipAddress)
}

func (s *TokenService) generateAccessToken(ipAddress string) (string, error) {
	jwtSecret := jwtReader.LoadJWTSecret(s.cfg.Env)

	claims := jwt.MapClaims{
		"ip":  ipAddress,
		"exp": time.Now().Add(s.cfg.Auth.AccessTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtSecret)
}

func (s *TokenService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func (s *TokenService) expired(token *model.RefreshTokenDB) error {
	if token.ExpiresAt.Before(time.Now()) {
		return errors.New("token is expired")
	}
	return nil
}
