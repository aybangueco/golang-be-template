package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/aybangueco/golang-be-template/internal/response"
	"github.com/aybangueco/golang-be-template/internal/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			pv := recover()
			if pv != nil {
				app.serverErrorResponse(w, r, fmt.Errorf("%v", pv))
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

func (app *application) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid IP", http.StatusInternalServerError)
			return
		}

		limiter := app.limiter.GetLimiter(ip)
		if !limiter.Allow() {
			app.tooManyRequestResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if authorization == "" {
			app.contextSetAuthenticatedUser(r, nil)
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			app.authorizationTokenMissingResponse(w, r)
			return
		}

		authorizationToken := parts[1]

		tokenClaims, err := token.ValidateToken(app.config.tokenSecret, authorizationToken)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				app.invalidAuthorizationTokenResponse(w, r)
				return
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				app.authorizationTokenExpiredResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		parsedUserID, err := uuid.Parse(tokenClaims.UserID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		user, err := app.db.GetUserById(r.Context(), parsedUserID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
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

		err := response.JSON(w, http.StatusUnauthorized, envelope{"error": "You need to be authenticated to access this resource"})
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	})
}
