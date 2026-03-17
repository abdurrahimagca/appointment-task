package mailer

import (
	"bytes"
	"html/template"
	"net/url"
	"strings"
	"time"
)

func BuildAppointmentCreatedTemplate(input MailInput) (string, error) {
	input.CalendarURL = buildGoogleCalendarURL(input)

	tmpl, err := template.New("appointment-created").Parse(appointmentCreatedTemplate)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, input); err != nil {
		return "", err
	}

	return body.String(), nil
}

func buildGoogleCalendarURL(input MailInput) string {
	values := url.Values{}
	values.Set("action", "TEMPLATE")
	values.Set("text", "Appointment with "+input.ReceiverName)
	values.Set("location", "Online")

	details := []string{"This meeting is held online."}
	if input.ProviderTimezone != "" {
		details = append(details, "Provider timezone: "+input.ProviderTimezone)
	}
	values.Set("details", strings.Join(details, " "))

	startUTC, startErr := time.Parse(time.RFC3339, input.StartTimeUTC)
	endUTC, endErr := time.Parse(time.RFC3339, input.EndTimeUTC)
	if startErr == nil && endErr == nil {
		values.Set("dates", formatGoogleCalendarTime(startUTC)+"/"+formatGoogleCalendarTime(endUTC))
	}

	return "https://www.google.com/calendar/render?" + values.Encode()
}

func formatGoogleCalendarTime(value time.Time) string {
	return value.UTC().Format("20060102T150405Z")
}

const appointmentCreatedTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        .calendar-btn {
            background-color: #4285F4;
            color: #ffffff !important;
            padding: 14px 24px;
            text-decoration: none;
            font-family: 'Helvetica', 'Arial', sans-serif;
            font-size: 16px;
            font-weight: bold;
            border-radius: 4px;
            display: inline-block;
        }
        .tiny-text {
            font-size: 10px;
            color: #9aa0a6;
            margin-top: 15px;
            font-family: 'Helvetica', 'Arial', sans-serif;
        }
    </style>
</head>
<body style="margin: 0; padding: 40px; background-color: #ffffff; text-align: center;">

    <div style="max-width: 500px; margin: auto;">
        <h2 style="color: #202124; font-family: Arial, sans-serif;">Created a new appointment</h2>
        <h3 style="color: #202124; font-family: Arial, sans-serif;">You and {{ .ReceiverName }} created a new appointment in {{ .Date }}</h3>
        <p style="color: #5f6368; font-family: Arial, sans-serif; font-size: 16px;">
            Click below to add this event directly to your Google Calendar.
        </p>

        <a href="{{ .CalendarURL }}"
           class="calendar-btn">
            Add to Google Calendar
        </a>

    </div>

</body>
</html>`
