package mailer

type MailInput struct {
	From             string
	To               string
	ReceiverName     string
	Date             string
	CalendarURL      string
	StartTimeUTC     string
	EndTimeUTC       string
	ProviderTimezone string
}
