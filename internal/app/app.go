package app

import (
	"MEDODS-test/internal/api"
	"MEDODS-test/internal/configure"
	sl "MEDODS-test/internal/lib/logger/slog"
	"MEDODS-test/internal/repo"
	"MEDODS-test/internal/service"
	"MEDODS-test/internal/util/jwtReader"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type App struct {
	api     *api.Api
	svc     *service.Service
	storage *repo.Storage
	echo    *echo.Echo
	dbpool  *pgxpool.Pool
}

func New(ctx context.Context, log *slog.Logger, cfg *configure.Config) *App {
	app := &App{}

	app.echo = echo.New()

	app.dbpool = configure.NewPostgres(ctx, cfg.Postgres)

	if err := cfg.Postgres.MigrationsUp(); err != nil && err.Error() != "no change" {
		log.Error("migration failed", sl.Err(err))
		panic(err)
	}

	jwtSecret := jwtReader.LoadJWTSecret()
	log.Info("secret:", string(jwtSecret))

	app.storage = repo.New(log, app.dbpool)

	app.svc = service.New(jwtSecret, cfg, log, app.storage, app.storage)

	app.api = api.New(log, *app.svc)

	app.echo.GET("/api/auth/generate", app.api.GetTokens)
	app.echo.GET("/api/auth/refresh", app.api.RefreshTokens)

	return app
}

func (a *App) Run() error {

	err := a.echo.Start(":8888")
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

	defer a.dbpool.Close()
	if err := a.echo.Shutdown(ctx); err != nil {
		fmt.Println("failed to shutdown server")
		return err
	}
	return nil
}
