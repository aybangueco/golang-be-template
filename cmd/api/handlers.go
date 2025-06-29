package main

import (
	"net/http"

	"github.com/aybangueco/golang-be-template/internal/response"
	"github.com/aybangueco/golang-be-template/internal/version"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"status":  "ok",
		"version": version.Get(),
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
