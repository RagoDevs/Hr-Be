package main

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	db "github.com/Hopertz/Hr-Be/internal/db/sqlc"
	"github.com/Hopertz/Hr-Be/internal/token"
	"github.com/labstack/echo/v4"
)

func (app *application) authenticate(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		authorizationHeader := c.Request().Header.Get("Authorization")

		if authorizationHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {

			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid auth token"})

		}

		tokenString := headerParts[1]
		tokenHash := sha256.Sum256([]byte(tokenString))

		params := db.GetUserByTokenParams{
			Hash:   tokenHash[:],
			Scope:  token.ScopeAuthentication,
			Expiry: time.Now(),
		}

		user, err := app.store.GetUserByToken(c.Request().Context(), params)

		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired auth token"})
			default:
				slog.Error("Error on getting user associated with token ", "Error", err.Error())
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			}
		}

		c.Set("user", user)

		return next(c)

	}

}

func (app *application) requireEnabledUser(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		user, ok := c.Get("user").(db.GetUserByTokenRow)

		if !ok {
			slog.Error("Error on requireEnabledUser middle ", "Error", "type assertion of user in ctx")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}

		if !user.IsEnabled {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "you are disabled user"})
		}

		return next(c)
	}
}


func (app *application) requireAdminOrHr(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		user, ok := c.Get("user").(db.GetUserByTokenRow)

		if !ok {
			slog.Error("Error on requireEnabledUser middle ", "Error", "type assertion of user in ctx")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}

		if user.RoleName != "hr" && user.RoleName != "admin" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "you must be an hr or admin"})
		}

		return next(c)
	}
}

