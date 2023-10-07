package middleware

import (
	"car_dealership/internal/auth"
	"car_dealership/internal/employee"
	"context"
	"net/http"
)

func AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the "Content-Type" header to "application/json"
		w.Header().Set("Content-Type", "application/json")

		// Check if the "Authorization" header is empty
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the token using your auth.ValidateToken function
		e, err := auth.ValidateToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		empl := employee.Find(e)
		if empl == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "employee", empl)
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
