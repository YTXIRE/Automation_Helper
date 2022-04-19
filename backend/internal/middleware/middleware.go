package middleware

import (
	"backend/internal/apperror"
	"net/http"
)

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(apperror.NewAppError(nil, "user unauthorized", "", "UNAUTHORIZED").Marshal())
			return
		}
		h(w, r)
	}
}
