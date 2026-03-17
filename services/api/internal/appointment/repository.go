package appointment

import (
	"context"
	"encoding/json"
	"time"

	db "github.com/abdurrahimagca/appointment-task/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	WithTx(tx pgx.Tx) *repository
	GetProviderByUsername(ctx context.Context, input GetProviderInput) (*Provider, error)
	GetSlotsByProviderUsernameAndDate(ctx context.Context, input GetAvailabilityInput, timezoneName string, cursor *SlotCursor) ([]AppointmentSlot, error)
	UpsertClientByEmail(ctx context.Context, input CreateAppointmentClient) (string, error)
	BookSlot(ctx context.Context, slotID string) (string, error)
	AddParticipantToAppointment(ctx context.Context, appointmentID string, clientID string) error
	GetAppointmentWithDetails(ctx context.Context, appointmentID string) (*Appointment, error)
}

type repository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *repository) WithTx(tx pgx.Tx) *repository {
	return &repository{
		pool:    r.pool,
		queries: r.queries.WithTx(tx),
	}
}

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *repository) GetProviderByUsername(ctx context.Context, input GetProviderInput) (*Provider, error) {
	dbProvider, err := r.queries.GetProviderByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	return toProvider(dbProvider), nil
}

func (r *repository) GetSlotsByProviderUsernameAndDate(ctx context.Context, input GetAvailabilityInput, timezoneName string, cursor *SlotCursor) ([]AppointmentSlot, error) {
	var dbParams db.GetSlotsByProviderUsernameAndDateParams
	dbParams.Username = input.Username
	dbParams.DateOfDay = input.DateOfDay
	dbParams.Limit = int32(input.Limit)

	if cursor != nil {
		var cursorStartTime pgtype.Timestamptz
		if err := cursorStartTime.Scan(cursor.StartTime); err != nil {
			return nil, err
		}
		var cursorID pgtype.UUID
		if err := cursorID.Scan(cursor.ID); err != nil {
			return nil, err
		}
		dbParams.CursorStartTime = cursorStartTime
		dbParams.CursorID = cursorID
	}

	rows, err := r.queries.GetSlotsByProviderUsernameAndDate(ctx, dbParams)
	if err != nil {
		return nil, err
	}

	slots := make([]AppointmentSlot, 0, len(rows))
	for _, row := range rows {
		slots = append(slots, toAppointmentSlot(row, timezoneName))
	}

	return slots, nil
}

func (r *repository) UpsertClientByEmail(ctx context.Context, input CreateAppointmentClient) (string, error) {
	var dbParams db.UpsertClientByEmailParams
	dbParams.FullName = input.FullName
	dbParams.Email = input.Email

	row, err := r.queries.UpsertClientByEmail(ctx, dbParams)
	if err != nil {
		return "", err
	}
	return row.ID.String(), nil
}

func (r *repository) BookSlot(ctx context.Context, slotID string) (string, error) {
	var parsedSlotID pgtype.UUID
	if err := parsedSlotID.Scan(slotID); err != nil {
		return "", err
	}

	row, err := r.queries.BookSlot(ctx, parsedSlotID)
	if err != nil {
		return "", err
	}

	return row.ID.String(), nil
}

func (r *repository) AddParticipantToAppointment(ctx context.Context, appointmentID string, clientID string) error {
	var dbParams db.AddParticipantToAppointmentParams
	if err := dbParams.AppointmentID.Scan(appointmentID); err != nil {
		return err
	}
	if err := dbParams.ClientID.Scan(clientID); err != nil {
		return err
	}

	return r.queries.AddParticipantToAppointment(ctx, dbParams)
}

func (r *repository) GetAppointmentWithDetails(ctx context.Context, appointmentID string) (*Appointment, error) {
	var parsedAppointmentID pgtype.UUID
	if err := parsedAppointmentID.Scan(appointmentID); err != nil {
		return nil, err
	}

	row, err := r.queries.GetAppointmentWithDetails(ctx, parsedAppointmentID)
	if err != nil {
		return nil, err
	}

	return toAppointment(row), nil
}

func toProvider(dbProvider db.GetProviderByUsernameRow) *Provider {
	provider := &Provider{
		ID:       dbProvider.ID.String(),
		FullName: dbProvider.FullName,
		Username: dbProvider.Username,
		Timezone: dbProvider.Timezone,
	}
	if dbProvider.Bio.Valid {
		provider.Bio = &dbProvider.Bio.String
	}
	return provider
}

func toAppointmentSlot(row db.GetSlotsByProviderUsernameAndDateRow, timezoneName string) AppointmentSlot {
	startUTC := row.StartTime.Time.UTC()
	endUTC := row.EndTime.Time.UTC()
	startLocal := toLocalTime(startUTC, timezoneName)
	endLocal := toLocalTime(endUTC, timezoneName)

	return AppointmentSlot{
		ID:                         row.ID.String(),
		DateOfDay:                  row.DateOfDay.Time.Format(time.DateOnly),
		Timezone:                   timezoneName,
		StartTimeUTC:               startUTC.Format(time.RFC3339),
		EndTimeUTC:                 endUTC.Format(time.RFC3339),
		StartTimeLocal:             startLocal.Format(time.RFC3339),
		EndTimeLocal:               endLocal.Format(time.RFC3339),
		AcceptMultipleParticipants: row.AvailableMultipleParticipants,
		IsAvailable:                true,
	}
}

func toAppointment(row db.GetAppointmentWithDetailsRow) *Appointment {
	startUTC := row.StartTime.Time.UTC()
	endUTC := row.EndTime.Time.UTC()
	startLocal := toLocalTime(startUTC, row.ProviderTimezone)
	endLocal := toLocalTime(endUTC, row.ProviderTimezone)

	participants := make([]AppointmentParticipant, 0)
	_ = json.Unmarshal(row.Participants, &participants)

	return &Appointment{
		ID:                         row.AppointmentID.String(),
		ProviderName:               row.ProviderName,
		ProviderTimezone:           row.ProviderTimezone,
		DateOfDay:                  row.DateOfDay.Time.Format(time.DateOnly),
		StartTimeUTC:               startUTC.Format(time.RFC3339),
		EndTimeUTC:                 endUTC.Format(time.RFC3339),
		StartTimeLocal:             startLocal.Format(time.RFC3339),
		EndTimeLocal:               endLocal.Format(time.RFC3339),
		AcceptMultipleParticipants: row.AvailableMultipleParticipants,
		Participants:               participants,
	}
}

func toLocalTime(value time.Time, timezoneName string) time.Time {
	location, err := time.LoadLocation(timezoneName)
	if err != nil {
		return value
	}
	return value.In(location)
}
