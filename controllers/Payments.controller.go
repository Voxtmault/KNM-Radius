package controllers

import (
	"CNM_Radius/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetPaymentsController(c echo.Context) error {
	result, err := models.GetAllPayments()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

//TO DO
//Add Checking to Payment Type
func CreateNewPaymentController(c echo.Context) error {
	//invoiceID int, amountPaid int, paymentType int, paymentNotes string
	invoiceID := c.FormValue("invoiceID")
	parsedInID, _ := strconv.Atoi(invoiceID)
	amountPaid := c.FormValue("amountPaid")
	parsedAmount, _ := strconv.Atoi(amountPaid)
	paymentType := c.FormValue("paymentType")
	parsedType, _ := strconv.Atoi(paymentType)
	paymentNotes := c.FormValue("paymentNotes")

	result, err := models.CreateNewPayment(parsedInID, parsedAmount, parsedType, paymentNotes)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)

}
