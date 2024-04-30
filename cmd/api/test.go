package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) testHandler(c echo.Context) error {

	test := map[string]string{
		"test": "test test test",
	}
	return c.JSON(http.StatusOK, test)

}
