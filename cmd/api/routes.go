package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.NotFound(app.notFound)

	mux.Use(app.recoverPanic)
	mux.Use(app.securityHeaders)

	mux.Route("/api", func(r chi.Router) {
		mux.Get("/health", app.health)
	})

	return mux
}
