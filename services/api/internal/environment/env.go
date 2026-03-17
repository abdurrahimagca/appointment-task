package environment

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Environment struct {
	DBURL          string
	APIVersion     string
	Port           string
	LogLevel       string
	SlogLevel      slog.Level
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	TurnstileSecret string
	EnvName         string
	ResendAPIKey    string
	ResendFromEmail string
	ResendFromName  string
}

func New() *Environment {
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")
	port := getEnvOrPanic("PORT")
	return &Environment{
		DBURL:          getEnvOrPanic("DB_URL"),
		APIVersion:     getEnvOrPanic("API_VERSION"),
		Port:           port,
		LogLevel:       logLevel,
		SlogLevel:      getSlogLevel(logLevel),
		AllowedOrigins: strings.Split(getEnvOrDefault("ALLOWED_ORIGINS", fmt.Sprintf("http://localhost:%s", port)), ","),
		AllowedMethods: strings.Split(getEnvOrDefault("ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
		AllowedHeaders: strings.Split(getEnvOrDefault("ALLOWED_HEADERS", "Content-Type,Authorization"), ","),
		TurnstileSecret: getEnvOrDefault("TURNSTILE_SECRET", ""),
		EnvName:         getEnvOrDefault("ENV_NAME", "development"),
		ResendAPIKey:    getEnvOrDefault("RESEND_API_KEY", ""),
		ResendFromEmail: getEnvOrDefault("RESEND_FROM_EMAIL", ""),
		ResendFromName:  getEnvOrDefault("RESEND_FROM_NAME", "Appointment"),
	}
}

func getSlogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}
	return slog.LevelInfo
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		err := fmt.Errorf("environment variable %s is not set", key)
		panic(err)
	}
	return value
}

func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
