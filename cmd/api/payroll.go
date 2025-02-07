package main

import (
	"log/slog"
	"net/http"
	"strconv"

	db "github.com/Hopertz/Hr-Be/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) createPayrollHandler(c echo.Context) error {

	var input struct {
		EmployeeID  uuid.UUID `json:"employee_id" validate:"required"`
		BasicSalary float64   `json:"basic_salary" validate:"required,gt=0"`
		TIN         string    `json:"tin" validate:"required"`
		BankName    string    `json:"bank_name" validate:"required"`
		BankAccount string    `json:"bank_account" validate:"required"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	rounded := float64(int(input.BasicSalary*100)) / 100

	s := strconv.FormatFloat(rounded, 'f', 2, 64)

	args := db.CreatePayrollParams{
		EmployeeID:  input.EmployeeID,
		BasicSalary: s,
		Tin:         input.TIN,
		BankName:    input.BankName,
		BankAccount: input.BankAccount,
	}

	err := app.store.CreatePayroll(c.Request().Context(), args)

	if err != nil {
		slog.Error("Error inserting payroll ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"success": "payroll created successful"})

}

func (app *application) getAllPayroll(c echo.Context) error {

	payroll, err := app.store.GetAllPayroll(c.Request().Context())

	if err != nil {
		slog.Error("Error getting all payroll ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, payroll)

}

func (app *application) getPayroll(c echo.Context) error {

	id := c.Param("payroll_id")

	payroll_id, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	payroll, err := app.store.GetPayroll(c.Request().Context(), payroll_id)

	if err != nil {
		slog.Error("Error getting payroll by id ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, payroll)

}

func (app *application) DeletePayroll(c echo.Context) error {

	id := c.Param("id")
	payroll_id, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	err = app.store.DeletePayroll(c.Request().Context(), payroll_id)
	if err != nil {
		slog.Error("Error deleting payroll ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)

}
