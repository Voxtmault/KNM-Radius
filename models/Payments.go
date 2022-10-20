package models

import (
	"CNM_Radius/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Payments struct {
	Id           string `json:Id`
	Invoiceid    string `json:Invoiceid `
	Amount       string `json:Amount`
	Date         string `json:date`
	Creationdate string `json:creationdate`
	Creationby   string `json:creationby`
}

func GetAllPayments() (Response, error) {
	var obj Payments
	var arrObj []Payments
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT id,invoice_id,amount,date,creationdate,creationby from payment"

	result, err := con.Query(sqlStatement)

	defer result.Close()

	if err != nil {
		return res, err
	}

	for result.Next() {
		err = result.Scan(&obj.Id, &obj.Invoiceid, &obj.Amount, &obj.Date, &obj.Creationdate, &obj.Creationby)

		if err != nil {
			return res, nil
		}

		arrObj = append(arrObj, obj)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrObj

	return res, nil
}

func CreateNewPayment(invoiceID int, amountPaid int, paymentType int, paymentNotes string) (Response, error) {
	var res Response
	var temp string
	var checkAmount int

	con := db.CreateCon()

	//Check For Amount To Pay invoice_items
	checkStatement := "SELECT CAST(amount AS CHAR) FROM invoice_items WHERE invoice_id = ?"

	checkResult, err := con.Query(checkStatement, invoiceID)
	if err != nil {
		return res, err
	}

	defer checkResult.Close()
	for checkResult.Next() {
		err = checkResult.Scan(&temp)
		if err != nil {
			return res, err
		}
	}
	temp = strings.Replace(temp, ".00", "", -1)
	checkAmount, _ = strconv.Atoi(temp)
	if amountPaid < checkAmount {
		res.Status = http.StatusOK
		res.Message = "Failed Creating New Payment"
		res.Data = "Please check the amount again"
	} else {
		paymentStatement := "INSERT INTO payment(invoice_id, amount, date, type_id, notes, creationdate, creationby) VALUES(?,?,?,?,?,?,?)"

		paymentSTMT, err := con.Prepare(paymentStatement)
		if err != nil {
			return Response{Status: http.StatusInternalServerError, Message: "Internal Server Error", Data: "Internal Server Error"}, err
		}

		paymentResult, err := paymentSTMT.Exec(invoiceID, amountPaid, time.Now().Format("2006-01-02 15:04:05"), paymentType, paymentNotes, time.Now().Format("2006-01-02 15:04:05"), "System")
		if err != nil {
			return Response{Status: http.StatusInternalServerError, Message: "Internal Server Error", Data: "Internal Server Error"}, err
		}

		NewPaymentID, _ := paymentResult.LastInsertId()
		defer paymentSTMT.Close()

		UpdateInvoice(invoiceID, 5)

		res.Status = http.StatusOK
		res.Message = "Success Adding Payment"
		res.Data = map[string]int64{
			"New Inserted Payment ID": NewPaymentID,
		}
	}

	return res, nil
}
