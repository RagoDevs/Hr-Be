package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/swaggo/echo-swagger/example/docs"
)

func (app *application) routes() *echo.Echo {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}

	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))

	e.GET("/ping", app.pingHandler)
	e.POST("/login", app.login)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	g := e.Group("/auth")

	g.Use(app.authenticate)
	g.Use(app.requireEnabledUser)

	g.GET("/test", app.testHandler)

	// employees
	g.GET("/employees", app.getAllEmployeesHandler, app.requireAdminOrHr)
	g.GET("/employees/:id", app.getEmployeeByIdHandler, app.requireAdminOrHr)
	g.POST("/employees", app.createEmployeeHandler, app.requireAdminOrHr)
	g.PUT("/employees/:id", app.updateEmployeeByIdHandler, app.requireAdminOrHr)
	g.DELETE("/employees/:id", app.deleteEmployeeByIdHandler, app.requireAdminOrHr)

	// leaves
	g.GET("/leaves", app.getAllLeavesHandler)
	g.GET("/leaves/:id", app.getLeaveByIdHandler)
	g.POST("/leaves", app.createLeaveHandler)
	g.PUT("/leaves/:id", app.updateLeaveByIdHandler)
	g.DELETE("/leaves/:id", app.deleteLeaveHandler)

	// contracts
	g.GET("/contracts/:employee_id", app.getAllContractsOfEmployeeByIdHandler)
	g.POST("/contracts", app.createContractHandler)
	g.PUT("/contracts/:contract_id", app.updateContractByIdHandler)
	g.DELETE("/contracts/:contract_id", app.deleteContractHandler)

	return e

}
