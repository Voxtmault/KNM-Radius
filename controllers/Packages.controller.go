package controllers

import (
	"CNM_Radius/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func CreateNewPlansController(c echo.Context) error {
	planName := c.FormValue("planName")
	planId := c.FormValue("planId")
	parsedPlID, _ := strconv.Atoi(planId)
	planType := c.FormValue("planType")
	planTimeType := c.FormValue("planTimeType")
	planRecurring := c.FormValue("Recurring")
	parsePlanRecurring, _ := strconv.ParseBool(planRecurring)
	planRecurringPeriod := c.FormValue("RecurringPeriod")
	planRecurringBillingSchedule := c.FormValue("BillingSchedule")
	planCost := c.FormValue("planCost")
	parsedPlanCost, _ := strconv.ParseFloat(planCost, 64)
	planSetupCost := c.FormValue("setupCost")
	parsedPlanSetupCost, _ := strconv.ParseFloat(planSetupCost, 64)
	planCurrency := c.FormValue("planCurrency")
	planActive := c.FormValue("planActive")
	profileName := c.FormValue("profileName")

	result, err := models.CreateNewPlans(
		planName,
		parsedPlID,
		planType,
		planTimeType,
		parsePlanRecurring,
		planRecurringPeriod,
		planRecurringBillingSchedule,
		parsedPlanCost,
		parsedPlanSetupCost,
		planCurrency,
		planActive,
		profileName,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
