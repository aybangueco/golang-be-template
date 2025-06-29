package main

import (
	"context"
	"net/http"

	"github.com/aybangueco/golang-be-template/internal/database"
)

type contextKey string

const (
	userCtxKey contextKey = "user"
)

func (app *application) contextSetAuthenticatedUser(r *http.Request, user *database.User) *http.Request {
	ctx := context.WithValue(r.Context(), userCtxKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetAuthenticatedUser(r *http.Request) *database.User {
	user, ok := r.Context().Value(userCtxKey).(*database.User)
	if !ok {
		return nil
	}

	return user
}
