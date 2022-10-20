package models

import (
	"CNM_Radius/db"
	"net/http"
)

func CreateNewPlans(
	planName string,
	planId int,
	planType string,
	planTimeType string,
	planRecurring bool,
	planRecurringPeriod string,
	planRecurringBillingSchedule string,
	planCost float64,
	planSetupCost float64,
	planCurrency string,
	planActive string,
	profileName string,
) (Response, error) {
	//Methods
	//1. Add to BillingPlans
	//2. Add to BillingPlansProfile
	var res Response

	con := db.CreateCon()

	billingPlansStatement := " INSERT INTO billing_plans(planName, planId, planType, planTimeType, planRecurring, planRecurringPeriod, planRecurringBillingSchedule, planCost, planSetupCost, planCurrency, planActive) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	billingPlansProfileStatement := "INSERT INTO billing_plans_profiles(plan_name, profile_name) VALUES(?,?)"

	_, err := con.Query(
		billingPlansStatement,
		planName,
		planId,
		planType,
		planTimeType,
		planRecurring,
		planRecurringPeriod,
		planRecurringBillingSchedule,
		planCost,
		planSetupCost,
		planCurrency,
		planActive,
	)
	if err != nil {
		return res, err
	}

	_, err = con.Query(billingPlansProfileStatement, planName, profileName)
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success Creating New Packet"
	res.Data = ""

	return res, nil
}
