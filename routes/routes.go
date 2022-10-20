package routes

import (
	"KNM-Radius/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Radius web Gateway")
	})

	//User Management
	//e.GET("/GetAllUsers", controllers.GetUsersController)

	//e.GET("/GetActiveUsers", controllers.GetActiveUsersController)

	//e.GET("/GetAvailableUsers", controllers.GetAvailableUsersController)

	//e.GET("/GetProfileHotspot", controllers.GetProfileHotspot)

	e.POST("/CreateNewUser", controllers.CreateUserController)

	return e
}
