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

func (app *application) createContractHandler(c echo.Context) error {

	var input struct {
		EmployeeID   uuid.UUID `json:"employee_id" validate:"required"`
		ContractType string    `json:"contract_type" validate:"oneof=temporary permanent"`
		StartDate    time.Time `json:"start_date" validate:"required"`
		EndDate      time.Time `json:"end_date" validate:"required"`
		Attachment   string    `json:"attachment" validate:"required"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	period := input.EndDate.Sub(input.StartDate).Hours() / 24

	args := db.CreateContractParams{
		EmployeeID:   input.EmployeeID,
		ContractType: input.ContractType,
		Attachment:   input.Attachment,
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		Period:       int32(period),
	}

	err := app.store.CreateContract(c.Request().Context(), args)

	if err != nil {
		slog.Error("Error creating contract ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"success": "contract created successful"})

}

func (app *application) getAllContractsOfEmployeeByIdHandler(c echo.Context) error {

	employee_id, err := uuid.Parse(c.Param("employee_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	contracts, err := app.store.GetContractsOfEmployeeByEmployeeId(c.Request().Context(), employee_id)

	if err != nil {
		slog.Error("Error getting contracts of employee ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, contracts)

}

func (app *application) updateContractByIdHandler(c echo.Context) error {

	contract_id, err := uuid.Parse(c.Param("contract_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	var input struct {
		ContractType *string    `json:"contract_type"`
		StartDate    *time.Time `json:"start_date"`
		EndDate      *time.Time `json:"end_date"`
		Attachment   *string    `json:"attachment"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	contract, err := app.store.GetContractById(c.Request().Context(), contract_id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "not found"})

		default:
			slog.Error("Error on get contract by id on updateContractByIdHandler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	if input.ContractType != nil {
		contract.ContractType = *input.ContractType
	}

	if input.StartDate != nil {
		contract.StartDate = *input.StartDate
	}

	if input.EndDate != nil {
		contract.EndDate = *input.EndDate
	}

	if input.Attachment != nil {
		contract.Attachment = *input.Attachment
	}

	contract.Period = int32(contract.EndDate.Sub(contract.StartDate).Hours() / 24)

	args := db.UpdateContractParams{
		ContractType: contract.ContractType,
		StartDate:    contract.StartDate,
		EndDate:      contract.EndDate,
		Period:       contract.Period,
		Attachment:   contract.Attachment,
		ID:           contract_id,
	}

	err = app.store.UpdateContract(c.Request().Context(), args)

	if err != nil {
		slog.Error("Error updating contract ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "contract updated successful"})

}

func (app *application) deleteContractHandler(c echo.Context) error {

	contract_id, err := uuid.Parse(c.Param("contract_id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	err = app.store.DeleteContract(c.Request().Context(), contract_id)

	if err != nil {
		slog.Error("Error deleting contract ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "contract deleted successful"})

}
