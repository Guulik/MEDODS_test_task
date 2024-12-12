package app

import (
	"MEDODS-test/internal/api"
	"MEDODS-test/internal/configure"
	"MEDODS-test/internal/repo"
	"MEDODS-test/internal/service"
	"MEDODS-test/internal/util/jwtReader"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type App struct {
	api     *api.Api
	svc     *service.Service
	storage *repo.Storage
	echo    *echo.Echo
}

func New(ctx context.Context, log *slog.Logger, cfg *configure.Config) *App {
	app := &App{}

	app.echo = echo.New()

	dbpool := configure.NewPostgres(ctx, cfg.Postgres)
	defer dbpool.Close()
	jwtSecret := jwtReader.LoadJWTSecret()

	app.storage = repo.New(log, dbpool)

	app.svc = service.New(jwtSecret, cfg, log, app.storage, app.storage)

	app.api = api.New(log, *app.svc)

	app.echo.GET("/auth/generate?", app.api.GetTokens)
	app.echo.GET("/auth/refresh", app.api.RefreshTokens)

	return app
}

func (a *App) Run() error {

	err := a.echo.Start(":8080")
	if err != nil {
		return err
	}

	fmt.Println("server running")
	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (a *App) Stop(ctx context.Context) error {
	fmt.Println("stopping server..." + " op = app.Stop")

	if err := a.echo.Shutdown(ctx); err != nil {
		fmt.Println("failed to shutdown server")
		return err
	}
	return nil
}
