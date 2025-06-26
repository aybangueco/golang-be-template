package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aybangueco/golang-be-template/internal/response"
	"github.com/aybangueco/golang-be-template/internal/token"
	"github.com/google/uuid"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			pv := recover()
			if pv != nil {
				app.serverError(w, r, fmt.Errorf("%v", pv))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := w.Header().Get("Authorization")

		if authorization == "" {
			app.contextSetAuthenticatedUser(r, nil)
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			app.invalidAuthorizationToken(w, r)
			return
		}

		authorizationToken := parts[1]

		tokenClaims, err := token.ValidateToken(app.config.tokenSecret, authorizationToken)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		parsedUserID, err := uuid.Parse(tokenClaims.UserID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		user, err := app.db.GetUserById(r.Context(), parsedUserID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		r = app.contextSetAuthenticatedUser(r, &user)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetAuthenticatedUser(r)

		if user != nil {
			next.ServeHTTP(w, r)
			return
		}

		err := response.JSON(w, http.StatusUnauthorized, "")
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	})
}
