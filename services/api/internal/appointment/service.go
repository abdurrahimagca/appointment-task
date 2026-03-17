package appointment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	localerrors "github.com/abdurrahimagca/appointment-task/internal/errors"
	"github.com/abdurrahimagca/appointment-task/internal/mailer"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
)

type service struct {
	logger     *slog.Logger
	repository Repository
	mailer     mailer.Service
}

type Service interface {
	GetProviderByUsername(ctx context.Context, input GetProviderInput) (*Provider, error)
	GetAvailability(ctx context.Context, input GetAvailabilityInput) (*AvailabilityResponse, error)
	CreateAppointment(ctx context.Context, input CreateAppointmentRequest) (*Appointment, error)
}

func NewService(repository Repository, mailerService mailer.Service, logger *slog.Logger) *service {
	return &service{logger: logger, repository: repository, mailer: mailerService}
}

func (s *service) GetProviderByUsername(ctx context.Context, input GetProviderInput) (*Provider, error) {
	return s.repository.GetProviderByUsername(ctx, input)
}

func (s *service) GetAvailability(ctx context.Context, input GetAvailabilityInput) (*AvailabilityResponse, error) {
	provider, err := s.repository.GetProviderByUsername(ctx, GetProviderInput{Username: input.Username})
	if err != nil {
		return nil, err
	}

	cursor, err := getCursor(input.CursorStartTime, input.CursorID)
	if err != nil {
		return nil, err
	}

	if input.Limit <= 0 {
		input.Limit = 10
	}
	input.Limit = input.Limit + 1

	slots, err := s.repository.GetSlotsByProviderUsernameAndDate(ctx, input, provider.Timezone, cursor)
	if err != nil {
		return nil, err
	}

	limit := input.Limit - 1
	response := &AvailabilityResponse{
		Username:  provider.Username,
		Timezone:  provider.Timezone,
		DateOfDay: input.DateOfDay,
	}

	if len(slots) > limit {
		next := slots[limit]
		response.NextCursor = &SlotCursor{StartTime: next.StartTimeUTC, ID: next.ID}
		slots = slots[:limit]
	}

	response.Slots = slots
	return response, nil
}

func (s *service) CreateAppointment(ctx context.Context, input CreateAppointmentRequest) (*Appointment, error) {
	clients := normalizeAppointmentClients(input)
	if len(clients) == 0 {
		return nil, huma.Error400BadRequest("at least one client is required")
	}

	s.logger.Info("creating appointment", "slot_id", input.SlotID, "client_count", len(clients), "primary_client_email", clients[0].Email)

	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return nil, localerrors.Wrapf(err, "begin transaction")
	}
	defer tx.Rollback(ctx)

	txRepo := s.repository.WithTx(tx)
	clientID, err := txRepo.UpsertClientByEmail(ctx, clients[0])
	if err != nil {
		return nil, localerrors.Wrapf(err, "upsert client")
	}

	appointmentID, err := txRepo.BookSlot(ctx, input.SlotID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, huma.Error409Conflict("Slot is no longer available")
		}
		return nil, localerrors.Wrapf(err, "book slot")
	}

	err = txRepo.AddParticipantToAppointment(ctx, appointmentID, clientID)
	if err != nil {
		return nil, localerrors.Wrapf(err, "add participant")
	}

	appointment, err := txRepo.GetAppointmentWithDetails(ctx, appointmentID)
	if err != nil {
		return nil, localerrors.Wrapf(err, "load appointment")
	}

	if len(clients) > 1 && !appointment.AcceptMultipleParticipants {
		return nil, huma.Error400BadRequest("selected slot does not allow multiple participants")
	}

	for _, client := range clients[1:] {
		clientID, err := txRepo.UpsertClientByEmail(ctx, client)
		if err != nil {
			return nil, localerrors.Wrapf(err, "upsert additional client")
		}

		err = txRepo.AddParticipantToAppointment(ctx, appointmentID, clientID)
		if err != nil {
			return nil, localerrors.Wrapf(err, "add additional participant")
		}
	}

	if len(clients) > 1 {
		appointment, err = txRepo.GetAppointmentWithDetails(ctx, appointmentID)
		if err != nil {
			return nil, localerrors.Wrapf(err, "reload appointment")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, localerrors.Wrapf(err, "commit transaction")
	}
	if s.mailer == nil {
		s.logger.Info("mailer is not configured, skipping appointment email")
	}
	if s.mailer != nil {
		for _, client := range clients {
			mailInput := mailer.MailInput{
				To:               client.Email,
				ReceiverName:     appointment.ProviderName,
				Date:             formatAppointmentWindow(appointment.StartTimeLocal, appointment.EndTimeLocal, appointment.ProviderTimezone),
				StartTimeUTC:     appointment.StartTimeUTC,
				EndTimeUTC:       appointment.EndTimeUTC,
				ProviderTimezone: appointment.ProviderTimezone,
			}
			err = s.mailer.SendAppointmentCreatedEmail(ctx, mailInput)
			if err != nil {
				s.logger.Error("failed to send appointment email", "error", err, "to", client.Email)
				return nil, localerrors.Wrapf(err, "send appointment email")
			}
		}
	}

	return appointment, nil
}

func normalizeAppointmentClients(input CreateAppointmentRequest) []CreateAppointmentClient {
	if len(input.Clients) > 0 {
		return input.Clients
	}
	if input.Client != nil {
		return []CreateAppointmentClient{*input.Client}
	}
	return nil
}

func formatAppointmentWindow(startTimeLocal string, endTimeLocal string, timezoneName string) string {
	start, startErr := time.Parse(time.RFC3339, startTimeLocal)
	end, endErr := time.Parse(time.RFC3339, endTimeLocal)
	if startErr != nil || endErr != nil {
		return fmt.Sprintf("%s - %s", startTimeLocal, endTimeLocal)
	}

	if start.Format(time.DateOnly) == end.Format(time.DateOnly) {
		return fmt.Sprintf("%s, %s - %s (%s)",
			start.Format("January 2, 2006"),
			start.Format("3:04 PM"),
			end.Format("3:04 PM"),
			timezoneName,
		)
	}

	return fmt.Sprintf("%s - %s (%s)",
		start.Format("January 2, 2006 3:04 PM"),
		end.Format("January 2, 2006 3:04 PM"),
		timezoneName,
	)
}

func getCursor(cursorStartTime string, cursorID string) (*SlotCursor, error) {
	if cursorStartTime == "" && cursorID == "" {
		return nil, nil
	}
	if cursorStartTime == "" || cursorID == "" {
		return nil, fmt.Errorf("cursorStartTime and cursorId must be provided together")
	}
	return &SlotCursor{StartTime: cursorStartTime, ID: cursorID}, nil
}
