package routes

import (
	"CNM_Radius/controllers"
	//"CNM_Radius/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Radius web Gateway")
	})

	//User Management
	e.GET("/GetAllUsers", controllers.GetUsersController)

	e.GET("/AddChild", controllers.AddChildController)

	e.GET("/GetActiveUsers", controllers.GetActiveUsersController)

	e.GET("/GetAvailableUsers", controllers.GetAvailableUsersController)

	e.GET("/GetUserProfile", controllers.GetUserProfileController)

	e.GET("/GetUserUsage", controllers.GetUserUsageController)

	e.GET("/GetInvoice", controllers.GetInvoiceController)

	e.GET("/GetPayments", controllers.GetPaymentsController)

	e.GET("/GetProfileHotspot", controllers.GetProfileHotspot)

	e.GET("/GetOpenInvoice", controllers.GetOpenInvoiceController)

	e.POST("/CreateNewUser", controllers.CreateUserController)

	e.POST("/CreateNewProfile", controllers.CreateNewProfileController)

	e.POST("/CreateNewPlans", controllers.CreateNewPlansController)

	e.POST("/CreateNewInvoice", controllers.CreateNewInvoiceController)

	e.POST("/CreateNewPayment", controllers.CreateNewPaymentController)

	return e
}
