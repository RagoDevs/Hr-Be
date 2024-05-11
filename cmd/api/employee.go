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
	"golang.org/x/crypto/bcrypt"
)

func (app *application) createEmployeeHandler(c echo.Context) error {

	var input struct {
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		Password    string    `json:"password"`
		Role        string    `json:"role"`
		DoB         time.Time `json:"dob"`
		Avatar      string    `json:"avatar"`
		Phone       string    `json:"phone"`
		Gender      string    `json:"gender"`
		JobTitle    string    `json:"job_title"`
		Department  string    `json:"department"`
		Address     string    `json:"address"`
		JoiningDate time.Time `json:"joining_date"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 6)
	if err != nil {
		slog.Error("Error hashing password ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	role_id, err := app.store.GetRoleByName(c.Request().Context(), input.Role)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "role not found"})

		default:
			slog.Error("Error get role by name on createEmployeeHandler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	user_args := db.CreateUserParams{
		RoleID:       role_id,
		Email:        input.Email,
		PasswordHash: password_hash,
	}

	employee_args := db.CreateEmployeeParams{
		Name:        input.Name,
		Dob:         input.DoB,
		Avatar:      input.Avatar,
		Phone:       input.Phone,
		Gender:      input.Gender,
		JobTitle:    input.JobTitle,
		Department:  input.Department,
		Address:     input.Address,
		JoiningDate: input.JoiningDate,
	}

	err = app.store.TxnCreateUserEmployee(c.Request().Context(), user_args, employee_args)

	if err != nil {
		slog.Error("Error creating user & employee ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"success": "employee created successful"})

}

func (app *application) updateEmployeeByIdHandler(c echo.Context) error {

	id := c.Param("id")

	employee_id, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	var input struct {
		Name        *string    `json:"name"`
		Email       *string    `json:"email"`
		Password    *string    `json:"password"`
		Role        *string    `json:"role"`
		DoB         *time.Time `json:"dob"`
		Avatar      *string    `json:"avatar"`
		Phone       *string    `json:"phone"`
		Gender      *string    `json:"gender"`
		JobTitle    *string    `json:"job_title"`
		Department  *string    `json:"department"`
		Address     *string    `json:"address"`
		IsEnabled   *bool      `json:"is_enabled"`
		IsPresent   *bool      `json:"is_present"`
		JoiningDate *time.Time `json:"joining_date"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	employee, err := app.store.GetEmployeeById(c.Request().Context(), employee_id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "employee not found"})

		default:
			slog.Error("Error get employee by id on updateEmployeeHandler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	user, err := app.store.GetUserById(c.Request().Context(), employee.UserID)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})

		default:
			slog.Error("Error get user by id on updateEmployeeHandler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	if input.Password != nil {

		password_hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), 6)
		if err != nil {
			slog.Error("Error hashing password ", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}

		user.PasswordHash = password_hash

	}

	if input.Role != nil {

		role_id, err := app.store.GetRoleByName(c.Request().Context(), *input.Role)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "role not found"})

			default:
				slog.Error("Error get role by name on updateEmployeeHandler", "Error", err.Error())
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			}
		}

		user.RoleID = role_id
	}

	if input.IsEnabled != nil {
		user.IsEnabled = *input.IsEnabled
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	user_args := db.UpdateUserByIdParams{
		RoleID:       user.RoleID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsEnabled:    user.IsEnabled,
		ID:           user.ID,
	}

	if input.Name != nil {
		employee.Name = *input.Name
	}

	if input.DoB != nil {
		employee.Dob = *input.DoB
	}

	if input.Avatar != nil {
		employee.Avatar = *input.Avatar
	}

	if input.Phone != nil {
		employee.Phone = *input.Phone
	}

	if input.Gender != nil {
		employee.Gender = *input.Gender
	}

	if input.JobTitle != nil {
		employee.JobTitle = *input.JobTitle
	}

	if input.Department != nil {
		employee.Department = *input.Department
	}

	if input.Address != nil {
		employee.Address = *input.Address
	}

	if input.JoiningDate != nil {
		employee.JoiningDate = *input.JoiningDate
	}

	if input.IsPresent != nil {
		employee.IsPresent = *input.IsPresent
	}

	employee_args := db.UpdateEmployeeByIdParams{
		Name:        employee.Name,
		Dob:         employee.Dob,
		Avatar:      employee.Avatar,
		Phone:       employee.Phone,
		Gender:      employee.Gender,
		JobTitle:    employee.JobTitle,
		Department:  employee.Department,
		Address:     employee.Address,
		JoiningDate: employee.JoiningDate,
		IsPresent:   employee.IsPresent,
		ID:          employee.ID,
	}

	err = app.store.TxnUpdateUserEmployee(c.Request().Context(), user_args, employee_args)

	if err != nil {
		slog.Error("Error creating user & employee ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "employee created successful"})

}

func (app *application) getAllEmployeesHandler(c echo.Context) error {

	employees, err := app.store.GetAllEmployees(c.Request().Context())

	if err != nil {
		slog.Error("Error getting employees ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, employees)

}

func (app *application) getEmployeeByIdHandler(c echo.Context) error {

	id := c.Param("id")

	employee_id, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	employee, err := app.store.GetEmployeeByIdDetailed(c.Request().Context(), employee_id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "employee not found"})

		default:
			slog.Error("Error get employee by id on getEmployeeByIdHandler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}
	return c.JSON(http.StatusOK, employee)

}

func (app *application) deleteEmployeeByIdHandler(c echo.Context) error {

	id := c.Param("id")

	employee_id, err := uuid.Parse(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}

	employee, err := app.store.GetEmployeeById(c.Request().Context(), employee_id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "employee not found"})

		default:
			slog.Error("Error get employee by id on deleteEmployeeHandler", "Error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	args := db.TxnUserEmployeeDelete{
		UserId:     employee.UserID,
		EmployeeId: employee_id,
	}
	err = app.store.TxnDeleteUserEmployee(c.Request().Context(), args)

	if err != nil {
		slog.Error("Error deleting user & empoyee ", "Error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": " user & employee deleted successful"})

}
