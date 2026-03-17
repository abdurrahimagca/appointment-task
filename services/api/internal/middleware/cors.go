package middleware

import (
	"net/http"
	"strings"

	"github.com/abdurrahimagca/appointment-task/internal/environment"
)

func CORS(next http.Handler, env *environment.Environment) http.Handler {
	allowedOrigins := toSet(env.AllowedOrigins)
	allowedMethods := strings.Join(env.AllowedMethods, ", ")
	allowedHeaders := strings.Join(env.AllowedHeaders, ", ")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && originAllowed(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func toSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		set[trimmed] = struct{}{}
	}
	return set
}

func originAllowed(origin string, allowedOrigins map[string]struct{}) bool {
	if _, ok := allowedOrigins["*"]; ok {
		return true
	}
	_, ok := allowedOrigins[origin]
	return ok
}
