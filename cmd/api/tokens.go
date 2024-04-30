package main

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	db "github.com/Hopertz/Hr-Be/internal/db/sqlc"
	"github.com/Hopertz/Hr-Be/internal/token"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type TokenResponseFmt struct {
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}

type LoginResponseFmt struct {
	UserID     uuid.UUID        `json:"user_id"`
	EmployeeID uuid.UUID           `json:"employee_id"`
	Email      string           `json:"email"`
	Role       string           `json:"role"`
	Name  string           `json:"name"`
	Avatar     string           `json:"avatar"`
	Token      TokenResponseFmt `json:"token"`
}

func (app *application) login(c echo.Context) error {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	user, err := app.store.GetUserByEmail(c.Request().Context(), input.Email)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid email or password"})

		default:
			slog.Error("Error on login Handler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(input.Password))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid password"})

		default:
			slog.Error("Error on password matches function ", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	token, tokenText, err := token.NewToken(user.UserID, 24*time.Hour, token.ScopeAuthentication)
	if err != nil {
		slog.Error("Error on password matches function ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})

	}

	err = app.store.CreateToken(c.Request().Context(), db.CreateTokenParams(*token))

	if err != nil {
		slog.Error("Error on password matches function ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}


	tokenResponse := TokenResponseFmt{
		Token:  tokenText,
		Expiry: token.Expiry.Unix(),
	}

	login := LoginResponseFmt{
		UserID:        user.UserID,
		EmployeeID: user.EmployeeID,
		Email:     user.Email,
		Role:      user.RoleName,
		Name: user.Name,
		Avatar:    user.Avatar,
		Token:     tokenResponse,
	}

	return c.JSON(200, login)

}
