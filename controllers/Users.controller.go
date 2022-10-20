package controllers

import (
	"KNM-Radius/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func CreateUserController(c echo.Context) error {
	//userCredentials string, maxUse int, expiredate int, profile string, sessionTimeout int, idleTimeout int, userCreator string
	userCredentials := c.FormValue("userCredentials")
	maxUse := c.FormValue("maxUse")
	parsedMaxuse, _ := strconv.Atoi(maxUse)
	expiredDate := c.FormValue("expiredDate")
	parsedExDate, _ := strconv.Atoi(expiredDate)
	sessionTimeout := c.FormValue("sessionTimeout")
	parsedSession, _ := strconv.Atoi(sessionTimeout)
	idleTimeout := c.FormValue("idleTimeout")
	parsedIdle, _ := strconv.Atoi(idleTimeout)
	userCreator := c.FormValue("userCreator")
	profileName := c.FormValue("profileName")

	result, err := models.CreateNewUser(userCredentials, parsedMaxuse, parsedExDate, profileName, parsedSession, parsedIdle, userCreator)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
