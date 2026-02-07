package api

import (
	"fmt"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := extractBearer(r)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, errBody("unauthorized", err.Error()))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func extractBearer(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	if len(auth) < 7 || strings.ToLower(auth[:7]) != "bearer " {
		return "", fmt.Errorf("Authorization must be 'Bearer <token>'")
	}

	token := strings.TrimSpace(auth[7:])
	if token == "" {
		return "", fmt.Errorf("empty bearer token")
	}

	return token, nil
}
