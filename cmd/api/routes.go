package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.MethodNotAllowed(app.methodNotAllowed)
	mux.NotFound(app.notFound)

	mux.Use(app.recoverPanic)
	mux.Use(app.securityHeaders)
	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{app.config.baseURL},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	mux.Use(app.authenticate)
	mux.Use(app.rateLimitMiddleware)

	mux.Get("/health", app.health)

	mux.Route("/auth", func(r chi.Router) {
		r.Get("/me", app.requireAuthenticated(app.handlerCurrentAuthenticated))

		r.Post("/register", app.handlerRegister)
		r.Post("/login", app.handlerLogin)
	})

	return mux
}
