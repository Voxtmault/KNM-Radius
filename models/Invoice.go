package models

import (
	"CNM_Radius/db"
	"net/http"
	"time"
)

type Invoice struct {
	Id           string `json:invoiceID`
	Uid          string `json:uid`
	Date         string `json:date`
	Notes        string `json:notes`
	Creationdate string `json:creationdate`
	Creationby   string `json:creationby`
}

func GetAllInvoice() (Response, error) {
	var obj Invoice
	var arrObj []Invoice
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT id,user_id,date,notes,creationdate,creationby from invoice "

	result, err := con.Query(sqlStatement)

	defer result.Close()

	if err != nil {
		return res, err
	}

	for result.Next() {
		err = result.Scan(&obj.Id, &obj.Uid, &obj.Date, &obj.Notes, &obj.Creationdate, &obj.Creationby)

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

func GetOpenInvoice() (Response, error) {
	var res Response
	var obj Invoice
	var arrObj []Invoice

	con := db.CreateCon()

	openStatement := "SELECT id, user_id, date, notes, creationdate, creationby FROM invoice WHERE status_id = ?"

	openResult, err := con.Query(openStatement, 1)
	if err != nil {
		return res, err
	}

	defer openResult.Close()

	for openResult.Next() {
		err = openResult.Scan(
			&obj.Id,
			&obj.Uid,
			&obj.Date,
			&obj.Notes,
			&obj.Creationdate,
			&obj.Creationby,
		)
		if err != nil {
			return res, err
		}

		arrObj = append(arrObj, obj)
	}

	res.Status = http.StatusOK
	res.Message = "Success Getting All Open / Unpaid Invoice"
	res.Data = arrObj

	return res, nil
}

func CreateNewInvoice(uid string, invoiceStatus int, invoiceType int, invoiceNote string, planID int, invoiceAmmount int, invoiceTax int, itemNotes string, transactionID int) (Response, error) {
	var res Response

	con := db.CreateCon()

	invoiceStatement := "INSERT INTO invoice (id, user_id, date, status_id, type_id, notes, creationdate, creationby) VALUES(?,?,?,?,?,?,?,?)"
	itemStatement := "INSERT INTO invoice_items(invoice_id, plan_id, amount, tax_amount, total, notes, creationdate, creationby) VALUES(?,?,?,?,?,?,?,?)"

	invoiceSTMT, err := con.Prepare(invoiceStatement)
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: "SQL Error", Data: err}, err
	}

	invoiceResult, err := invoiceSTMT.Exec(transactionID, uid, time.Now().Format("2006-01-02 15:04:05"), invoiceStatus, invoiceType, invoiceNote, time.Now().Format("2006-01-02 15:04:05"), "System")
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: "SQL Error", Data: err}, err
	}

	LastInsertedID, _ := invoiceResult.LastInsertId()
	//Or u can use transaction id, it's completely up to u
	defer invoiceSTMT.Close()

	itemSTMT, err := con.Prepare(itemStatement)
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: "SQL Error", Data: err}, err
	}

	itemResult, err := itemSTMT.Exec(LastInsertedID, planID, invoiceAmmount, invoiceTax, invoiceAmmount*((invoiceTax+100)/100), itemNotes, time.Now().Format("2006-01-02 15-04-05"), "System")
	if err != nil {
		return Response{Status: http.StatusInternalServerError, Message: "SQL Error", Data: err}, err
	}

	LastItemID, _ := itemResult.LastInsertId()
	defer itemSTMT.Close()

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int64{
		"Last Inserted Invoice ID": LastInsertedID,
		"Last Inserted Item ID":    LastItemID,
	}

	return res, nil
}

func UpdateInvoice(invoiceID int, newStatus int) (Response, error) {
	var res Response
	var validateInvoice int

	con := db.CreateCon()

	validateInvoice = 0

	//Validate Invoice
	validateStatement := "SELECT COUNT(*) FROM invoice WHERE id = ?"

	validateResult, err := con.Query(validateStatement, invoiceID)
	if err != nil {
		return res, err
	}

	defer validateResult.Close()

	for validateResult.Next() {
		err = validateResult.Scan(&validateInvoice)
		if err != nil {
			return res, err
		}
	}

	if validateInvoice > 0 {

		updateStatement := "UPDATE invoice SET status_id = ?, updatedate = ?, updateby = ? WHERE id = ?"

		updateSTMT, err := con.Prepare(updateStatement)
		if err != nil {
			return res, err
		}

		updateResult, err := updateSTMT.Exec(newStatus, time.Now().Format("2006-01-02 15:04:05"), "System", invoiceID)
		if err != nil {
			return res, err
		}

		AffectedRow, _ := updateResult.RowsAffected()
		defer updateSTMT.Close()

		if AffectedRow > 0 {
			res.Status = http.StatusOK
			res.Message = "Success Updating Invoice"
			res.Data = "Success Updating Invoice"
		} else {
			res.Status = http.StatusOK
			res.Message = "Failed Updating Invoice"
			res.Data = "Please contact Tech Support to resolve this issue"
		}

	} else {
		res.Status = http.StatusOK
		res.Message = "Invalid Invoice"
		res.Data = "Invalid Invoice, Please Check Your Information Again"
	}

	return res, nil
}
