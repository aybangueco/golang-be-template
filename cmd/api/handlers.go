package main

import (
	"net/http"
	"time"

	"github.com/aybangueco/golang-be-template/internal/response"
	"github.com/aybangueco/golang-be-template/internal/version"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"status":  "ok",
		"version": version.Get(),
		"uptime":  time.Since(time.Now()).String(),
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
