package main

import (
	"context"
	"errors"
	"fmt"
	"jsunnykhan/go-clean-template/internal/adapter/http/handler"
	"jsunnykhan/go-clean-template/internal/adapter/http/router"
	repository "jsunnykhan/go-clean-template/internal/adapter/repository"
	"jsunnykhan/go-clean-template/internal/core/config"
	"jsunnykhan/go-clean-template/internal/core/database"
	userusecase "jsunnykhan/go-clean-template/internal/usecase/user"
	"jsunnykhan/go-clean-template/pkg/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()

	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	// 2. Initialize logger
	log := logger.New(cfg.App.Mode)
	log.Info("starting server", slog.String("env", cfg.App.Mode))

	// 3. Initialize database connection

	userRepo := repository.NewUserRepository(nil)

	db, err := database.ConnectPostgres(cfg.Database,
		userRepo.AutoMigrateModel(),
	)

	if err != nil {
		log.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("database connection established")

	userRepo = repository.NewUserRepository(db)

	// use-case  will go here

	userUC := userusecase.New(userRepo)

	// handler initialization will go here

	userHandler := handler.NewUserHandler(userUC, cfg.JWT)

	// router and server setup will go here

	// ── 7. Router ─────────────────────────────────────────────────────────────
	engine := router.NewRouter(userHandler, cfg.JWT, log.Logger)

	// ── 8. HTTP Server with graceful shutdown ─────────────────────────────────
	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine to allow for graceful shutdown
	go func() {
		log.Info("HTTP server listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Block until SIGINT or SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", slog.Any("error", err))
	}

	log.Info("server stopped")

}
