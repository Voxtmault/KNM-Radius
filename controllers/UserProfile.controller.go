package controllers

import (
	"CNM_Radius/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUserProfileController(c echo.Context) error {
	result, err := models.GetUserProfile()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// func GetAvailableUsersController(c echo.Context) error {
// 	username := strings.ReplaceAll(c.FormValue("username"), " ", "")
// 	sid := c.FormValue("sid")

// 	fmt.Println(sid + username)
// 	result, err := models.GetAvailableUser(sid, username)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, result)
// }
