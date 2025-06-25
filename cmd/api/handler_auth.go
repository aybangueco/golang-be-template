package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/aybangueco/golang-be-template/internal/database"
	"github.com/aybangueco/golang-be-template/internal/hashing"
	"github.com/aybangueco/golang-be-template/internal/request"
	"github.com/aybangueco/golang-be-template/internal/response"
	"github.com/aybangueco/golang-be-template/internal/token"
	"github.com/aybangueco/golang-be-template/internal/validator"
)

func (app *application) handlerCurrentAuthenticated(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetAuthenticatedUser(r)

	err := response.JSON(w, http.StatusOK, map[string]any{
		"user": user,
	})
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) handlerRegister(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Validator validator.Validator
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.FirstName != "", "firstName", "First name is required")
	input.Validator.CheckField(input.LastName != "", "lastName", "Last name is required")
	input.Validator.CheckField(input.Email != "", "email", "Email is required")
	input.Validator.CheckField(input.Password != "", "password", "Password is required")

	if input.Validator.HasErrors() {
		app.validationFailed(w, r, input.Validator)
		return
	}

	_, err = app.db.GetUserByEmail(r.Context(), input.Email)
	existingEmail := errors.Is(err, sql.ErrNoRows)
	if err != nil && !existingEmail {
		app.serverError(w, r, err)
		return
	}

	hashedPassword, err := hashing.HashPassword(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(!existingEmail, "email", "Email is already taken")

	if input.Validator.HasErrors() {
		app.validationFailed(w, r, input.Validator)
		return
	}

	createdUser, err := app.db.CreateUser(r.Context(), database.CreateUserParams{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	token, err := token.GenerateToken(app.config.tokenSecret, createdUser.ID.String())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusCreated, map[string]any{
		"token": token,
	})
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Validator validator.Validator
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.Email != "", "email", "Email is required")
	input.Validator.CheckField(input.Password != "", "password", "Password is required")
	if input.Validator.HasErrors() {
		app.validationFailed(w, r, input.Validator)
		return
	}

	user, err := app.db.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.customError(w, r, "Invalid email or password", http.StatusBadRequest)
			return
		}
		app.serverError(w, r, err)
		return
	}

	if !hashing.IsPasswordValid(user.Password, input.Password) {
		app.customError(w, r, "Invalid email or password", http.StatusBadRequest)
		return
	}

	token, err := token.GenerateToken(app.config.tokenSecret, user.ID.String())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, map[string]any{
		"token": token,
	})
	if err != nil {
		app.serverError(w, r, err)
	}
}
