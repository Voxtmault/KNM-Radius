package controllers

import (
	"CNM_Radius/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func GetUsersController(c echo.Context) error {
	result, err := models.GetAllUsers()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func CreateUserController(c echo.Context) error {
	username := strings.ReplaceAll(c.FormValue("username"), " ", "")
	sid := c.FormValue("sid")
	max_device := c.FormValue("packet_max")
	parsed_max, _ := strconv.Atoi(max_device)
	email := c.FormValue("email")
	mobilephone := c.FormValue("phone")
	parsedPhone, _ := strconv.Atoi(mobilephone)
	address := c.FormValue("address")
	planname := c.FormValue("packetname")
	paymentMethod := c.FormValue("paymentMethod")
	expireDate := c.FormValue("expiredDate")
	parsedExDate, _ := strconv.Atoi(expireDate)
	ordertaker := c.FormValue("orderTaker")
	profileName := c.FormValue("profileName")

	result, err := models.CreateNewUser(username, sid, parsed_max, email, parsedPhone, address, planname, paymentMethod, ordertaker, parsedExDate, profileName)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func AddChildController(c echo.Context) error {
	result, err := models.AddChildToUser()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func GetActiveUsersController(c echo.Context) error {
	result, err := models.GetActiveUsers()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func GetAvailableUsersController(c echo.Context) error {
	username := strings.ReplaceAll(c.FormValue("username"), " ", "")
	sid := c.FormValue("sid")

	fmt.Println(sid + username)
	result, err := models.GetAvailableUser(sid, username)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
