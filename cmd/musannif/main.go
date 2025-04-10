package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/masroof-maindak/musannif/internal/config"
	"github.com/masroof-maindak/musannif/internal/db"
	"github.com/masroof-maindak/musannif/internal/logger"
	"github.com/masroof-maindak/musannif/internal/middlewares"
	"github.com/masroof-maindak/musannif/internal/routes"
	"github.com/masroof-maindak/musannif/internal/utils"

	"github.com/MadAppGang/httplog"
)

func initialize() error {
	err := config.Initialize()
	if err != nil {
		return fmt.Errorf("error initializing config: %v\n", err)
	}

	err = logger.Initialize(logger.LoggerConfig{
		InfoLogPath:  path.Join(config.Cfg.App.LogDirectory, "info.log"),
		ErrorLogPath: path.Join(config.Cfg.App.LogDirectory, "error.log"),
	})
	if err != nil {
		log.Fatalf("error initializing logger: %v\n", err)
	}

	utils.SetJwtKeys(config.Cfg.Secrets.JWT_ACCESS_SECRET, config.Cfg.Secrets.JWT_REFRESH_SECRET)

	err = db.InitDb(config.Cfg.App.SqliteDirectory)
	if err != nil {
		log.Fatalf("error initializing db: %v\n", err)
	}

	return nil
}

func main() {
	err := initialize()
	ctx := context.Background()
	err = run(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	defer func() {
		if err := db.CleanupDb(); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to cleanup database")
		}
	}()

	srv := newServer(&config.Cfg)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Cfg.Server.Host, strconv.Itoa(config.Cfg.Server.Port)),
		Handler: srv,
	}

	go func() {
		logger.Log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to shutdown server")
		}
	}()

	wg.Wait()

	return nil
}

func newServer(config *config.AppConfig) http.Handler {
	mux := http.NewServeMux()

	routes.AddRoutes(mux)

	var handler http.Handler = middlewares.CORS(mux)

	if config.App.Environment == "debug" {
		handler = httplog.Logger(handler)
	}

	return handler
}
