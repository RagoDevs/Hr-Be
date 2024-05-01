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

func (app *application) createLeaveHandler(c echo.Context) error {

	var input struct {
		EmployeeID   uuid.UUID `json:"employee_id"`
		ApprovedByID uuid.UUID `json:"approved_by_id"`
		Description  string    `json:"description"`
		StartDate    time.Time `json:"start_date"`
		EndDate      time.Time `json:"end_date"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	duration := input.EndDate.Sub(input.StartDate)

	leave_count := duration.Hours() / 24

	args := db.CreateLeaveParams{
		EmployeeID:   input.EmployeeID,
		ApprovedByID: input.ApprovedByID,
		Description:  input.Description,
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		LeaveCount:   int16(leave_count),
	}

	err := app.store.CreateLeave(c.Request().Context(), args)

	if err != nil {
		slog.Error("Error creating leave ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"success": "leave created successful"})

}

func (app *application) getAllLeavesHandler(c echo.Context) error {

	leaves, err := app.store.GetAllLeaves(c.Request().Context())

	if err != nil {
		slog.Error("Error deleting leave ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, leaves)

}

func (app *application) getLeaveByIdHandler(c echo.Context) error {

	id := c.Param("id")

	leave_id, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	leave, err := app.store.GetLeaveById(c.Request().Context(), leave_id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "not found"})

		default:
			slog.Error("Error get leave by id on getLeaveByIdHandler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}
	return c.JSON(http.StatusOK, leave)

}

func (app *application) updateLeaveByIdHandler(c echo.Context) error {

	id := c.Param("id")

	leave_id, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	var input struct {
		ApprovedByID *uuid.UUID `json:"approved_by_id"`
		Approved     *bool      `json:"approved"`
		Description  *string    `json:"description"`
		StartDate    *time.Time `json:"start_date"`
		EndDate      *time.Time `json:"end_date"`
		Seen         *bool      `json:"seen"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	leave, err := app.store.GetLeaveById(c.Request().Context(), leave_id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "not found"})

		default:
			slog.Error("Error on get leave by id on updateHandler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	if input.ApprovedByID != nil {
		leave.ApprovedByID = *input.ApprovedByID
	}

	if input.Approved != nil {
		leave.Approved = *input.Approved
	}

	if input.Description != nil {
		leave.Description = *input.Description
	}

	if input.StartDate != nil {
		leave.StartDate = *input.StartDate
	}

	if input.EndDate != nil {
		leave.EndDate = *input.EndDate
	}

	if input.Seen != nil {
		leave.Seen = *input.Seen
	}

	leave.LeaveCount = int16(leave.EndDate.Sub(leave.StartDate).Hours() / 24)

	args := db.UpdateLeaveParams{
		ApprovedByID: leave.ApprovedByID,
		Approved:     leave.Approved,
		Description:  leave.Description,
		StartDate:    leave.StartDate,
		EndDate:      leave.EndDate,
		LeaveCount:   leave.LeaveCount,
		Seen:         leave.Seen,
		ID:           leave_id,
	}

	err = app.store.UpdateLeave(c.Request().Context(), args)

	if err != nil {
		slog.Error("Error updating leave ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "leave updated successful"})

}

func (app *application) deleteLeaveHandler(c echo.Context) error {

	id := c.Param("id")

	leave_id, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	err = app.store.DeleteLeave(c.Request().Context(), leave_id)

	if err != nil {
		slog.Error("Error deleting leave ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "leave deleted successful"})

}
