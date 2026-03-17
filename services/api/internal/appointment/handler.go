package appointment

import (
	"context"

	localerrors "github.com/abdurrahimagca/appointment-task/internal/errors"
	"github.com/abdurrahimagca/appointment-task/internal/turnstile"
	"github.com/danielgtaylor/huma/v2"
)

type handler struct {
	service   *service
	turnstile turnstile.Service
}

type Handler interface {
	GetProviderByUsername(ctx context.Context, input *GetProviderInput) (*GetProviderOutput, error)
	GetAvailability(ctx context.Context, input *GetAvailabilityInput) (*GetAvailabilityOutput, error)
	CreateAppointment(ctx context.Context, input *CreateAppointmentInput) (*CreateAppointmentOutput, error)
}

func NewHandler(service *service, turnstileService turnstile.Service) *handler {
	return &handler{service: service, turnstile: turnstileService}
}

func (h *handler) GetProviderByUsername(ctx context.Context, input *GetProviderInput) (*GetProviderOutput, error) {
	provider, err := h.service.GetProviderByUsername(ctx, *input)
	if err != nil {
		return nil, localerrors.AsHumaError(err)
	}

	output := &GetProviderOutput{}
	output.Body.Provider = *provider
	return output, nil
}

func (h *handler) GetAvailability(ctx context.Context, input *GetAvailabilityInput) (*GetAvailabilityOutput, error) {
	availability, err := h.service.GetAvailability(ctx, *input)
	if err != nil {
		return nil, localerrors.AsHumaError(err)
	}

	output := &GetAvailabilityOutput{}
	output.Body = *availability
	return output, nil
}

func (h *handler) CreateAppointment(ctx context.Context, input *CreateAppointmentInput) (*CreateAppointmentOutput, error) {
	if h.turnstile != nil {
		ok, err := h.turnstile.Verify(ctx, input.Body.TurnstileToken)
		if err != nil {
			return nil, localerrors.AsHumaError(err)
		}
		if !ok {
			return nil, huma.Error400BadRequest("Turnstile verification failed")
		}
	}

	appointment, err := h.service.CreateAppointment(ctx, input.Body)
	if err != nil {
		return nil, localerrors.AsHumaError(err)
	}

	output := &CreateAppointmentOutput{}
	output.Body.Appointment = *appointment
	return output, nil
}
