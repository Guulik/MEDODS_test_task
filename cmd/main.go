package main

import (
	"MEDODS-test/internal/app"
	"MEDODS-test/internal/configure"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	// можно было бы под разные среды: чтобы в Local был Debug, а в проде Info, но решил не нагромождать в рамках тестового
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	cfg := configure.MustConfig()

	a := app.New(ctx, log, cfg)

	go func() {
		a.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	if err := a.Stop(ctx); err != nil {
		fmt.Println(fmt.Errorf("failed to gracefully stop app err=%s", err))
	}

	fmt.Println("Gracefully stopped")
}
