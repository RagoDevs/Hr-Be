package main

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	db "github.com/Hopertz/Hr-Be/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) createAnnouncementHandler(c echo.Context) error {

	var input struct {
		Description      string    `json:"description" validate:"required"`
		CreatedBy        uuid.UUID `json:"createdby_id" validate:"required"`
		AnnouncementDate time.Time `json:"date" validate:"required"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	args := db.CreateAnnouncementParams{
		Description:      input.Description,
		CreatedBy:        input.CreatedBy,
		AnnouncementDate: input.AnnouncementDate,
	}

	err := app.store.CreateAnnouncement(c.Request().Context(), args)

	if err != nil {
		slog.Error("Error creating announcement ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"success": "announcement created successful"})

}

func (app *application) getAllAnnouncementsHandler(c echo.Context) error {

	announcements, err := app.store.GetAnnouncements(c.Request().Context())

	if err != nil {
		slog.Error("Error getting announcements ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, announcements)

}

func (app *application) updateAnnouncementByIdHandler(c echo.Context) error {

	announcement_id, err := uuid.Parse(c.Param("announcement_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	var input struct {
		Description      *string    `json:"description"`
		AnnouncementDate *time.Time `json:"date"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	announcement, err := app.store.GetAnnouncementById(c.Request().Context(), announcement_id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "not found"})

		default:
			slog.Error("Error on get announcement by id on updateAnnouncementHandler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	if input.AnnouncementDate != nil {
		announcement.AnnouncementDate = *input.AnnouncementDate
	}

	if input.Description != nil {
		announcement.Description = *input.Description
	}

	args := db.UpdateAnnouncementParams{
		Description:      announcement.Description,
		AnnouncementDate: announcement.AnnouncementDate,
		ID:               announcement.ID,
	}

	err = app.store.UpdateAnnouncement(c.Request().Context(), args)

	if err != nil {
		slog.Error("Error updating announcement ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "announcement updated successful"})

}

func (app *application) deleteAnnouncementHandler(c echo.Context) error {

	announcement_id, err := uuid.Parse(c.Param("announcement_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	err = app.store.DeleteAnnouncement(c.Request().Context(), announcement_id)

	if err != nil {
		slog.Error("Error deleting announcement ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "announcement deleted successful"})

}
