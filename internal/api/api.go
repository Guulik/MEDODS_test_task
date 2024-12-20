package api

import (
	"MEDODS-test/internal/domain/model"
	"MEDODS-test/internal/domain/request"
	sl "MEDODS-test/internal/lib/logger/slog"
	"MEDODS-test/internal/service"
	"MEDODS-test/internal/util/binder"
	"MEDODS-test/internal/util/ip"
	goContext "context"
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type Api struct {
	log *slog.Logger
	svc *service.Service
}

func New(
	log *slog.Logger,
	service *service.Service,
) *Api {
	return &Api{
		log: log,
		svc: service,
	}
}

func (a *Api) GetTokens(ctx echo.Context) error {
	log := a.log.With(
		slog.String("op", "Api.GetTokens"),
	)
	userIp := ip.GetIPAddress(ctx.Request())
	if userIp == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to extract IP"))
	}
	var (
		req       request.GetTokensRequest
		tokenPair *model.TokenPair
		err       error
	)

	context := goContext.Background()
	err = binder.BindReq(log, ctx, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	log.Info(sl.Req(req))

	tokenPair, err = a.svc.Tokens.GenerateTokens(context, req.UserGUID, userIp)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *tokenPair)
}

func (a *Api) RefreshTokens(ctx echo.Context) error {
	log := a.log.With(
		slog.String("op", "Api.RefreshTokens"),
	)
	userIp := ip.GetIPAddress(ctx.Request())
	if userIp == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to extract IP"))
	}
	var (
		req       request.RefreshTokensRequest
		tokenPair *model.TokenPair
		err       error
	)

	context := goContext.Background()
	err = binder.BindReq(log, ctx, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	log.Info(sl.Req(req))

	tokenPair, err = a.svc.Tokens.RefreshTokens(context, req.UserGUID, req.RefreshToken, userIp)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *tokenPair)
}
