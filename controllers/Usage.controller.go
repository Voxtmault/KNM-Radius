package controllers

import (
	"KNM-Radius/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetUserUsageController(c echo.Context) error {
	username := c.FormValue("username")

	result, err := models.GetUserUsage(username)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
