package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/masroof-maindak/musannif/internal/config"
	"github.com/masroof-maindak/musannif/internal/db"
	"github.com/masroof-maindak/musannif/internal/middlewares"
	"github.com/masroof-maindak/musannif/internal/routes"
	"github.com/masroof-maindak/musannif/internal/utils"

	"github.com/MadAppGang/httplog"
)

func init() {
	err := config.Initialize()
	if err != nil {
		log.Fatalf("Error initializing config: %v\n", err)
	}

	// TODO: init logger!

	utils.SetJwtKeys(config.Cfg.Secrets.JWT_ACCESS_SECRET, config.Cfg.Secrets.JWT_REFRESH_SECRET)

	err = db.InitDb(config.Cfg.App.SqliteDirectory)
	if err != nil {
		log.Fatalf("Error initializing db: %v\n", err)
	}
}

func main() {
	ctx := context.Background()
	err := run(ctx)
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
		logger.Log.Printf("Listening on %s\n", httpServer.Addr)
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

	if config.App.Environment == "dev" {
		handler = httplog.Logger(handler)
	}

	return handler
}

