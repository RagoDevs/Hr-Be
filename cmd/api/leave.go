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

type LeaveDistribution struct {
	Sick     int64 `json:"sick"`
	Personal int64 `json:"personal"`
	Paid     int64 `json:"paid"`
	Annual   int64 `json:"annual"`
}

func (app *application) createLeaveHandler(c echo.Context) error {

	var input struct {
		EmployeeID   uuid.UUID `json:"employee_id" validate:"required"`
		ApprovedByID uuid.UUID `json:"approved_by_id" validate:"required"`
		LeaveType    string    `json:"leave_type" validate:"oneof=paid sick annual personal"`
		Description  string    `json:"description" validate:"required"`
		StartDate    time.Time `json:"start_date" validate:"required"`
		EndDate      time.Time `json:"end_date" validate:"required"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	duration := input.EndDate.Sub(input.StartDate)

	leave_count := duration.Hours() / 24

	args := db.CreateLeaveParams{
		EmployeeID:   input.EmployeeID,
		ApprovedByID: input.ApprovedByID,
		LeaveType:    input.LeaveType,
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

func (app *application) getAllLeavesSeenUnseensHandler(c echo.Context) error {

	leaves, err := app.store.GetAllLeavesUnSeenTopBottomSeen(c.Request().Context())

	if err != nil {
		slog.Error("Error getting leaves ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, leaves)

}

func (app *application) getAllLeavesApprovedHandler(c echo.Context) error {

	leaves, err := app.store.GetAllLeavesApproved(c.Request().Context())

	if err != nil {
		slog.Error("Error getting leaves ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, leaves)

}

func (app *application) getLeavesByEmployeeIdHandler(c echo.Context) error {

	employee_id, err := uuid.Parse(c.Param("employee_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	leave, err := app.store.GetLeavesByEmployeeId(c.Request().Context(), employee_id)

	if err != nil {
		slog.Error("Error get leaves of employee_id on getLeavesByEmployeeIdHandler", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, leave)

}

func (app *application) updateLeaveByIdHandler(c echo.Context) error {

	leave_id, err := uuid.Parse(c.Param("leave_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	var input struct {
		ApprovedByID *uuid.UUID `json:"approved_by_id"`
		Approved     *bool      `json:"approved"`
		LeaveType    *string    `json:"leave_type"`
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
			slog.Error("Error on get leave by id on updateLeaveByIdHandler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	if input.ApprovedByID != nil {
		leave.ApprovedByID = *input.ApprovedByID
	}

	if input.Approved != nil {
		leave.Approved = *input.Approved
	}

	if input.LeaveType != nil {
		leave.LeaveType = *input.LeaveType
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
		LeaveType:    leave.LeaveType,
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

	leave_id, err := uuid.Parse(c.Param("leave_id"))

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

func (app *application) respondToLeaveByIdHandler(c echo.Context) error {

	leave_id, err := uuid.Parse(c.Param("leave_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	var input struct {
		Approved bool `json:"approved"`
	}

	if err := c.Bind(&input); err != nil {

		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	args := db.ApproveRejectLeaveParams{
		Approved: input.Approved,
		ID:       leave_id,
	}

	err = app.store.TxnAcceptedRejectLeave(c.Request().Context(), args)

	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "not found"})

		default:
			slog.Error("Error approve/reject leave ", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}

	}

	return c.JSON(http.StatusOK, map[string]string{"success": "leave responded successful"})

}

func (app *application) getAllLeaveApprovers(c echo.Context) error {

	approvers, err := app.store.GetAllAApprovers(c.Request().Context())

	if err != nil {
		slog.Error("Error get leaves of employee_id on getAllLeaveApprovers", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, approvers)

}

func (app *application) getAllLeaveOnleave(c echo.Context) error {

	leaves, err := app.store.GetEmployeesOnLeave(c.Request().Context())

	if err != nil {
		slog.Error("Error getting onleave employees", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, leaves)

}

func (app *application) getAllLeaveUpcoming(c echo.Context) error {

	leaves, err := app.store.GetEmployeesOnLeaveUpcoming(c.Request().Context())

	if err != nil {
		slog.Error("Error getting upcoming leaves employees", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, leaves)

}

func (app *application) getLeaveCountDistr(c echo.Context) error {

	leaves, err := app.store.GetLeaveCountDistr(c.Request().Context())

	if err != nil {
		slog.Error("Error getting  leaves counts distribution  of employees", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	var ld LeaveDistribution

	for _, leave := range leaves {
		if leave.LeaveType == "sick" {
			ld.Sick = leave.LeaveCount
		} else if leave.LeaveType == "paid" {
			ld.Paid = leave.LeaveCount
		} else if leave.LeaveType == "personal" {
			ld.Personal = leave.LeaveCount
		} else if leave.LeaveType == "annual" {
			ld.Annual = leave.LeaveCount
		}

	}
	return c.JSON(http.StatusOK, ld)

}
