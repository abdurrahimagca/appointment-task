package appointment

type Provider struct {
	ID       string  `json:"id" format:"uuid"`
	FullName string  `json:"fullName"`
	Username string  `json:"username"`
	Timezone string  `json:"timezone"`
	Bio      *string `json:"bio" required:"false" doc:"Provider bio in markdown format"`
}

type SlotCursor struct {
	StartTime string `json:"startTime" format:"date-time"`
	ID        string `json:"id" format:"uuid"`
}

type AppointmentSlot struct {
	ID                         string `json:"id" format:"uuid"`
	DateOfDay                  string `json:"dateOfDay" format:"date"`
	Timezone                   string `json:"timezone"`
	StartTimeUTC               string `json:"startTimeUtc" format:"date-time"`
	EndTimeUTC                 string `json:"endTimeUtc" format:"date-time"`
	StartTimeLocal             string `json:"startTimeLocal" format:"date-time"`
	EndTimeLocal               string `json:"endTimeLocal" format:"date-time"`
	AcceptMultipleParticipants bool   `json:"acceptMultipleParticipants"`
	IsAvailable                bool   `json:"isAvailable"`
}


type AppointmentParticipant struct {
	ClientID string `json:"clientId" format:"uuid"`
	FullName string `json:"fullName"`
	Email    string `json:"email" format:"email"`
}

type Appointment struct {
	ID                         string                   `json:"id" format:"uuid"`
	ProviderName               string                   `json:"providerName"`
	ProviderTimezone           string                   `json:"providerTimezone"`
	DateOfDay                  string                   `json:"dateOfDay" format:"date"`
	StartTimeUTC               string                   `json:"startTimeUtc" format:"date-time"`
	EndTimeUTC                 string                   `json:"endTimeUtc" format:"date-time"`
	StartTimeLocal             string                   `json:"startTimeLocal" format:"date-time"`
	EndTimeLocal               string                   `json:"endTimeLocal" format:"date-time"`
	AcceptMultipleParticipants bool                     `json:"acceptMultipleParticipants"`
	Participants               []AppointmentParticipant `json:"participants"`
}
