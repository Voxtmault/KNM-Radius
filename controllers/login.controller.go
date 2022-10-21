package controllers

import (
	"KNM-Radius/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func LoginController(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	result, err := models.Authenticate(email, password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
