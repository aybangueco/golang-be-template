package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/aybangueco/golang-be-template/internal/response"
	"github.com/aybangueco/golang-be-template/internal/validator"
)

func (app *application) reportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	requestAttrs := slog.Group("request", "method", method, "url", url)
	app.logger.Error(message, requestAttrs, "trace", trace)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.reportServerError(r, err)

	message := "The server encountered a problem and could not process your request"
	err = response.JSON(w, http.StatusInternalServerError, envelope{"error": message})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The requested method %s is not supported for specified resource", r.Method)

	err := response.JSON(w, http.StatusMethodNotAllowed, envelope{"error": message})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"

	err := response.JSON(w, http.StatusNotFound, envelope{"error": message})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	err = response.JSON(w, http.StatusBadRequest, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) tooManyRequest(w http.ResponseWriter, r *http.Request) {
	message := "Too many request"

	err := response.JSON(w, http.StatusTooManyRequests, envelope{"error": message})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) inputValidationFailed(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	err := response.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) invalidAuthorizationToken(w http.ResponseWriter, r *http.Request) {
	message := "Invalid authorization token"

	err := response.JSON(w, http.StatusUnauthorized, envelope{"error": message})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) invalidUserCredentials(w http.ResponseWriter, r *http.Request) {
	message := "Invalid email or password"

	err := response.JSON(w, http.StatusUnauthorized, envelope{"error": message})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
