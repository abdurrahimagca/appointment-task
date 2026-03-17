package mailer

import (
	"context"
	"log/slog"

	"github.com/abdurrahimagca/appointment-task/internal/environment"
	localerrors "github.com/abdurrahimagca/appointment-task/internal/errors"
	"github.com/resend/resend-go/v3"
)

type Service interface {
	SendAppointmentCreatedEmail(ctx context.Context, input MailInput) error
}

type service struct {
	logger      *slog.Logger
	client      *resend.Client
	fromEmail   string
	fromName    string
	isAvailable bool
}

func New(env *environment.Environment, logger *slog.Logger) Service {
	apiKey := env.ResendAPIKey
	fromEmail := env.ResendFromEmail
	fromName := env.ResendFromName
	isAvailable := apiKey != "" && fromEmail != ""
	if !isAvailable {
		logger.Info("resend is not configured, mail delivery disabled", "has_api_key", apiKey != "", "from_email_configured", fromEmail != "")
	} else {
		logger.Info("resend configured", "from_email", fromEmail, "from_name", fromName)
	}

	var client *resend.Client
	if isAvailable {
		client = resend.NewClient(apiKey)
	}

	return &service{
		logger:      logger,
		client:      client,
		fromEmail:   fromEmail,
		fromName:    fromName,
		isAvailable: isAvailable,
	}
}

func (s *service) SendAppointmentCreatedEmail(ctx context.Context, input MailInput) error {
	if err := ctx.Err(); err != nil {
		return localerrors.Wrapf(err, "appointment email context")
	}

	if !s.isAvailable {
		s.logger.Info("resend is not available, skipping mail delivery", "to", input.To)
		return nil
	}

	body, err := BuildAppointmentCreatedTemplate(input)
	if err != nil {
		s.logger.Error("failed to build appointment created email template", "error", err, "to", input.To)
		return localerrors.Wrapf(err, "build appointment created template")
	}
	from := input.From
	if from == "" {
		from = s.fromEmail
		if s.fromName != "" {
			from = s.fromName + " <" + s.fromEmail + ">"
		}
	}

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{input.To},
		Html:    body,
		Subject: "Appointment created",
	}

	s.logger.Info("sending appointment created email via resend", "to", input.To, "from", from)
	sent, err := s.client.Emails.Send(params)
	if err != nil {
		s.logger.Error("failed to send appointment created email via resend", "error", err, "to", input.To, "from", from)
		return localerrors.Wrapf(err, "send appointment created email")
	}

	s.logger.Info("appointment created email sent", "to", input.To, "email_id", sent.Id)

	return nil
}
