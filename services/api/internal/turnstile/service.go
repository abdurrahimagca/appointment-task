package turnstile

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/abdurrahimagca/appointment-task/internal/environment"
	localerrors "github.com/abdurrahimagca/appointment-task/internal/errors"
	"github.com/danielgtaylor/huma/v2"
)

type Service interface {
	Verify(ctx context.Context, token string) (bool, error)
}

type service struct {
	logger   *slog.Logger
	secret   string
	env_name string
}

func NewService(env *environment.Environment, logger *slog.Logger) Service {
	return &service{logger: logger, secret: env.TurnstileSecret, env_name: env.EnvName}
}

func (s *service) Verify(ctx context.Context, token string) (bool, error) {
	if s.secret == "" {
		if s.env_name == "production" {
			s.logger.Error("turnstile secret is not configured in production")
			return false, huma.Error500InternalServerError("appointment creation is unavailable, please check server logs for more details")
		}
		s.logger.Info("turnstile secret is not set, skipping verification", "env", s.env_name)
		return true, nil
	}

	if s.env_name != "production" {
		s.logger.Info("turnstile verification skipped in non-production environment", "env", s.env_name)
		return true, nil
	}

	if token == "" {
		return false, huma.Error400BadRequest("turnstileToken is required")
	}

	url := "https://challenges.cloudflare.com/turnstile/v0/siteverify?secret=" + s.secret + "&response=" + token
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, localerrors.Wrapf(err, "create turnstile request")
	}

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return false, localerrors.Wrapf(err, "execute turnstile request")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, localerrors.Wrapf(err, "read turnstile response")
	}

	var responseBody struct {
		Success bool `json:"success"`
	}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return false, localerrors.Wrapf(err, "unmarshal turnstile response")
	}
	if !responseBody.Success {
		return false, huma.Error400BadRequest("Turnstile verification failed")
	}
	return true, nil
}
