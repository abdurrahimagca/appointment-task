package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/abdurrahimagca/appointment-task/internal/appointment"
	"github.com/abdurrahimagca/appointment-task/internal/environment"
	"github.com/abdurrahimagca/appointment-task/internal/mailer"
	"github.com/abdurrahimagca/appointment-task/internal/middleware"
	"github.com/abdurrahimagca/appointment-task/internal/turnstile"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	env := environment.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: env.SlogLevel}))
	slog.SetDefault(logger)
	requestLogger := logger.With("component", "http")
	appointmentLogger := logger.With("component", "appointment")
	mailerLogger := logger.With("component", "mailer")
	turnstileLogger := logger.With("component", "turnstile")
	pool, err := pgxpool.New(context.Background(), env.DBURL)
	if err != nil {
		logger.Error("failed to create pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	for i := range 30 {
		if err := pool.Ping(context.Background()); err == nil {
			break
		} else if i == 29 {
			logger.Error("database not reachable after 30 attempts", "error", err)
			os.Exit(1)
		} else {
			logger.Warn("waiting for database", "attempt", i+1, "error", err)
			time.Sleep(time.Second)
		}
	}

	router := http.NewServeMux()
	config := huma.DefaultConfig("Appointment API", env.APIVersion)
	config.DocsRenderer = huma.DocsRendererScalar
	api := humago.New(router, config)
	mailerService := mailer.New(env, mailerLogger)
	turnstileService := turnstile.NewService(env, turnstileLogger)
	appointment.Register(api, pool, mailerService, turnstileService, appointmentLogger)

	addr := fmt.Sprintf(":%s", env.Port)
	logger.Info("server starting", "address", addr)
	err = http.ListenAndServe(addr, middleware.RequestLog(middleware.CORS(router, env), requestLogger))
	if err != nil {
		logger.Error("server stopped", "error", err)
	}
	os.Exit(0)
}
