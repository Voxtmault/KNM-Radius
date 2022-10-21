package controllers

import (
	"KNM-Radius/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CreateNewProfileController(c echo.Context) error {
	profileName := c.FormValue("profileName")
	profileAttribute := c.FormValue("profileAttribute")
	profileOperator := c.FormValue("profileOperator")
	profileValue := c.FormValue("profileValue")
	parsedValue, _ := strconv.Atoi(profileValue)

	result, err := models.CreateNewProfile(profileName, profileAttribute, profileOperator, parsedValue)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteProfileController(c echo.Context) error {
	profileName := c.FormValue("profileName")

	result, err := models.DeleteProfile(profileName)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func GetProfileHotspot(c echo.Context) error {
	result, err := models.GetProfileHotspot()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
