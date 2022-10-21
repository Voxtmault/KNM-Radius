package controllers

import (
	"KNM-Radius/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func GetAllUsersController(c echo.Context) error {
	result, err := models.GetAllUsers()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func CreateUserController(c echo.Context) error {
	//userCredentials string, expireDate int, profile string, sessionTimeout int, idleTimeout int, deviceOwner string, deviceInfo string
	userCredentials := c.FormValue("userCredentials")
	expiredDate := c.FormValue("expiredDate")
	parsedExDate, _ := strconv.Atoi(expiredDate)
	profileName := c.FormValue("profileName")
	sessionTimeout := c.FormValue("sessionTimeout")
	parsedSession, _ := strconv.Atoi(sessionTimeout)
	idleTimeout := c.FormValue("idleTimeout")
	parsedIdle, _ := strconv.Atoi(idleTimeout)
	deviceOwner := c.FormValue("deviceOwner")
	deviceInfo := c.FormValue("deviceInfo")

	result, err := models.CreateNewUser(userCredentials, parsedExDate, profileName, parsedSession, parsedIdle, deviceOwner, deviceInfo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteUserController(c echo.Context) error {
	userCredentials := c.FormValue("username")

	result, err := models.DeleteUser(userCredentials)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func EditUserController(c echo.Context) error {
	userCredentials := c.FormValue("userCredentials")
	expiredDate := c.FormValue("expiredDate")
	parsedExDate, _ := strconv.Atoi(expiredDate)
	profileName := c.FormValue("profileName")
	sessionTimeout := c.FormValue("sessionTimeout")
	parsedSession, _ := strconv.Atoi(sessionTimeout)
	idleTimeout := c.FormValue("idleTimeout")
	parsedIdle, _ := strconv.Atoi(idleTimeout)
	deviceOwner := c.FormValue("deviceOwner")
	deviceInfo := c.FormValue("deviceInfo")

	result, err := models.EditUser(userCredentials, parsedExDate, profileName, parsedSession, parsedIdle, deviceOwner, deviceInfo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
