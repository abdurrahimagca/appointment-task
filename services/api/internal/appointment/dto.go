package appointment

type GetProviderInput struct {
	Username string `path:"username" validate:"required,min=3,max=255"`
}

type GetProviderOutput struct {
	Body struct {
		Provider Provider `json:"provider"`
	}
}

type GetAvailabilityInput struct {
	Username        string `path:"username" validate:"required,min=3,max=255"`
	DateOfDay       string `query:"date" format:"date" validate:"required"`
	Limit           int    `query:"limit" default:"10" minimum:"1" maximum:"100"`
	CursorStartTime string `query:"cursorStartTime" format:"date-time" required:"false"`
	CursorID        string `query:"cursorId" format:"uuid" required:"false"`
}

type AvailabilityResponse struct {
	Username   string            `json:"username"`
	Timezone   string            `json:"timezone"`
	DateOfDay  string            `json:"dateOfDay" format:"date"`
	Slots      []AppointmentSlot `json:"slots"`
	NextCursor *SlotCursor       `json:"nextCursor" required:"false"`
}
type GetAvailabilityOutput struct {
	Body AvailabilityResponse `json:"body"`
}

type CreateAppointmentInput struct {
	Body CreateAppointmentRequest
}

type CreateAppointmentRequest struct {
	SlotID         string                    `json:"slotId" format:"uuid"`
	TurnstileToken string                    `json:"turnstileToken" required:"false"`
	Client         *CreateAppointmentClient  `json:"client" required:"false"`
	Clients        []CreateAppointmentClient `json:"clients" required:"false"`
}

type CreateAppointmentClient struct {
	FullName string `json:"fullName"`
	Email    string `json:"email" format:"email"`
}

type CreateAppointmentOutput struct {
	Body struct {
		Appointment Appointment `json:"appointment"`
	}
}
