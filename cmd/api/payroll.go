package main

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	db "github.com/Hopertz/Hr-Be/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EmployeePayroll struct {
	ID              uuid.UUID `json:"id"`
	EmployeeID      uuid.UUID `json:"employee_id"`
	BasicSalary     float64   `json:"basic_salary"`
	Tin             string    `json:"tin"`
	BankName        string    `json:"bank_name"`
	BankAccount     string    `json:"bank_account"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	EmployeeName    string    `json:"employee_name"`
	TaxableIncome   float64   `json:"taxable_income"`
	PAYE            float64   `json:"paye"`
	Loan            float64   `json:"loan"`
	TotalDeductions float64   `json:"total_deductions"`
	NSSFEmployee    float64   `json:"nssf_employee"`
	NHIFEmployee    float64   `json:"nhif_employee"`
	Department      string    `json:"department"`
}

type EmployeePayrollDetails struct {
	PayrollID           uuid.UUID `json:"payroll_id"`
	EmployeeID          uuid.UUID `json:"employee_id"`
	BasicSalary         float64   `json:"basic_salary,string"`
	TIN                 string    `json:"tin"`
	BankName            string    `json:"bank_name"`
	BankAccount         string    `json:"bank_account"`
	PayrollIsActive     bool      `json:"payroll_is_active"`
	PayrollCreatedAt    time.Time `json:"payroll_created_at"`
	PayrollUpdatedAt    time.Time `json:"payroll_updated_at"`
	EmployeeName        string    `json:"employee_name"`
	DateOfBirth         time.Time `json:"dob"`
	Avatar              string    `json:"avatar"`
	Phone               string    `json:"phone"`
	Gender              string    `json:"gender"`
	JobTitle            string    `json:"job_title"`
	Department          string    `json:"department"`
	Address             string    `json:"address"`
	EmployeeIsPresent   bool      `json:"employee_is_present"`
	EmployeeJoiningDate time.Time `json:"employee_joining_date"`
	EmployeeCreatedAt   time.Time `json:"employee_created_at"`
	TaxableIncome       float64   `json:"taxable_income"`
	PAYE                float64   `json:"paye"`
	Loan                float64   `json:"loan"`
	TotalDeductions     float64   `json:"total_deductions"`
	NSSFEmployee        float64   `json:"nssf_employee"`
	NHIFEmployee        float64   `json:"nhif_employee"`
}

func CalculateMonthlyTax(income float64) float64 {
	switch {
	case income <= 270000:
		return 0

	case income <= 520000:
		return (income - 270000) * 0.08

	case income <= 760000:
		return 20000 + ((income - 520000) * 0.20)

	case income <= 1000000:
		return 68000 + ((income - 760000) * 0.25)

	default: // income > 1000000
		return 128000 + ((income - 1000000) * 0.30)
	}
}

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

	IsEmployee, err := app.store.IsEmployeeExisting(c.Request().Context(), args.EmployeeID)

	if err != nil {
		slog.Error("Error Checking If employee exists ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	if IsEmployee {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "employee id exists"})
	}

	err = app.store.CreatePayroll(c.Request().Context(), args)

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

	e_payrolls := []EmployeePayroll{}

	for _, p := range payroll {

		bs, err := strconv.ParseFloat(p.BasicSalary, 64)

		if err != nil {
			slog.Error("failed to parse float", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})

		}

		nssfEmployee := bs * 0.10

		nhifEmployee := max(bs*0.03, 20000)

		taxableIncome := bs - nssfEmployee

		paye := CalculateMonthlyTax(taxableIncome)

		loan := 0.0

		totalDeductions := nssfEmployee + nhifEmployee + paye + loan

		e := EmployeePayroll{
			ID:              p.ID,
			EmployeeID:      p.EmployeeID,
			BasicSalary:     bs,
			Tin:             p.Tin,
			BankName:        p.BankName,
			BankAccount:     p.BankAccount,
			IsActive:        p.IsActive,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
			EmployeeName:    p.EmployeeName,
			TaxableIncome:   taxableIncome,
			PAYE:            paye,
			Loan:            loan,
			TotalDeductions: totalDeductions,
			NSSFEmployee:    nssfEmployee,
			NHIFEmployee:    nhifEmployee,
			Department:      p.Department,
		}

		e_payrolls = append(e_payrolls, e)

	}

	return c.JSON(http.StatusOK, e_payrolls)

}

func (app *application) getPayroll(c echo.Context) error {
	id := c.Param("payroll_id")

	payrollID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	payroll, err := app.store.GetPayroll(c.Request().Context(), payrollID)
	if err != nil {
		slog.Error("Error getting payroll by id", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	bs, err := strconv.ParseFloat(payroll.BasicSalary, 64)
	if err != nil {
		slog.Error("Failed to parse float", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	nssfEmployee := bs * 0.10
	nhifEmployee := bs * 0.03
	if nhifEmployee < 20000 {
		nhifEmployee = 20000
	}

	taxableIncome := bs - nssfEmployee
	paye := CalculateMonthlyTax(taxableIncome)
	loan := 0.0
	totalDeductions := nssfEmployee + nhifEmployee + paye + loan

	e := EmployeePayrollDetails{
		PayrollID:           payroll.PayrollID,
		EmployeeID:          payroll.EmployeeID,
		BasicSalary:         bs,
		TIN:                 payroll.Tin,
		BankName:            payroll.BankName,
		BankAccount:         payroll.BankAccount,
		PayrollIsActive:     payroll.PayrollIsActive,
		PayrollCreatedAt:    payroll.PayrollCreatedAt,
		PayrollUpdatedAt:    payroll.PayrollUpdatedAt,
		EmployeeName:        payroll.EmployeeName,
		DateOfBirth:         payroll.Dob,
		Avatar:              payroll.Avatar,
		Phone:               payroll.Phone,
		Gender:              payroll.Gender,
		JobTitle:            payroll.JobTitle,
		Department:          payroll.Department,
		Address:             payroll.Address,
		EmployeeIsPresent:   payroll.EmployeeIsPresent,
		EmployeeJoiningDate: payroll.EmployeeJoiningDate,
		EmployeeCreatedAt:   payroll.EmployeeCreatedAt,
		TaxableIncome:       taxableIncome,
		PAYE:                paye,
		Loan:                loan,
		TotalDeductions:     totalDeductions,
		NSSFEmployee:        nssfEmployee,
		NHIFEmployee:        nhifEmployee,
	}

	return c.JSON(http.StatusOK, e)
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

func (app *application) updatePayrollHandler(c echo.Context) error {

	var input struct {
		BasicSalary *float64 `json:"basic_salary"`
		TIN         *string  `json:"tin"`
		BankName    *string  `json:"bank_name"`
		BankAccount *string  `json:"bank_account"`
		IsActive    *bool    `json:"is_active"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	e, err := app.store.JustGetPayroll(c.Request().Context())

	if err != nil {
		slog.Error("Error Getting employee for update ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	if input.BasicSalary != nil {
		rounded := float64(int(*input.BasicSalary*100)) / 100
		bs := strconv.FormatFloat(rounded, 'f', 2, 64)
		e.BasicSalary = bs
	}

	if input.BankAccount != nil {
		e.BankAccount = *input.BankAccount
	}

	if input.BankName != nil {
		e.BankName = *input.BankName
	}

	if input.TIN != nil {
		e.Tin = *input.TIN
	}

	if input.IsActive != nil {
		e.IsActive = *input.IsActive
	}

	args := db.UpdatePayrollParams{
		ID:          e.ID,
		BasicSalary: e.BasicSalary,
		Tin:         e.Tin,
		BankName:    e.BankName,
		BankAccount: e.BankAccount,
		IsActive:    e.IsActive,
	}

	err = app.store.UpdatePayroll(c.Request().Context(), args)

	if err != nil {
		slog.Error("Error updating payroll ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "payroll update successful"})

}
