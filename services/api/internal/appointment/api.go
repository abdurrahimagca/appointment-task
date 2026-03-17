package appointment

import (
	"log/slog"

	"github.com/abdurrahimagca/appointment-task/internal/mailer"
	"github.com/abdurrahimagca/appointment-task/internal/turnstile"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(api huma.API, pool *pgxpool.Pool, mailerService mailer.Service, turnstileService turnstile.Service, logger *slog.Logger) Service {
	r := NewRepository(pool)
	s := NewService(r, mailerService, logger)
	h := NewHandler(s, turnstileService)

	huma.Register(api, huma.Operation{
		Method:                       "GET",
		Path:                         "/providers/{username}",
		Description:                  "Get a provider by username",
		Summary:                      "Get a provider by username",
		RejectUnknownQueryParameters: true,
		Errors:                       []int{422, 404, 500},
	}, h.GetProviderByUsername)

	huma.Register(api, huma.Operation{
		Method:                       "GET",
		Path:                         "/providers/{username}/availability",
		Description:                  "Get provider availability for a local date",
		Summary:                      "Get provider availability",
		RejectUnknownQueryParameters: true,
		Errors:                       []int{400, 404, 422, 500},
	}, h.GetAvailability)

	huma.Register(api, huma.Operation{
		Method:      "POST",
		Path:        "/appointments",
		Description: "Create an appointment for a free slot",
		Summary:     "Create appointment",
		Errors:      []int{400, 404, 409, 422, 500},
	}, h.CreateAppointment)

	return s
}
