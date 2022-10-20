package controllers

import (
	"CNM_Radius/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetInvoiceController(c echo.Context) error {
	result, err := models.GetAllInvoice()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func GetOpenInvoiceController(c echo.Context) error {
	result, err := models.GetOpenInvoice()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func CreateNewInvoiceController(c echo.Context) error {
	//transactionID int, uid string, invoiceStatus int, invoiceType int, invoiceNote string, planID int, invoiceAmmount int, invoiceTax int, itemNotes string

	transactionID := c.FormValue("transactionID")
	parsedTID, _ := strconv.Atoi(transactionID)
	uid := c.FormValue("UID")
	invoiceStatus := c.FormValue("statusID")
	parsedStatus, _ := strconv.Atoi(invoiceStatus)
	invoiceType := c.FormValue("typeID")
	parsedType, _ := strconv.Atoi(invoiceType)
	invoiceNote := c.FormValue("invoiceNote")
	planID := c.FormValue("planID")
	parsedID, _ := strconv.Atoi(planID)
	invoiceAmmount := c.FormValue("invAmmount")
	parsedAmmound, _ := strconv.Atoi(invoiceAmmount)
	invoiceTax := c.FormValue("invTax")
	parsedTax, _ := strconv.Atoi(invoiceTax)
	itemNotes := c.FormValue("itemNotes")

	result, err := models.CreateNewInvoice(uid, parsedStatus, parsedType, invoiceNote, parsedID, parsedAmmound, parsedTax, itemNotes, parsedTID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
