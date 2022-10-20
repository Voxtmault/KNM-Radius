package models

import (
	"CNM_Radius/db"
	"net/http"
)

type Usage struct {
	Username string  `json:username`
	Uptime   float64 `json:uptime`
	Download float64 `json:download`
	Upload   float64 `json:upload`
}

func GetUserUsage(username string) (Response, error) {
	var res Response
	var selectedRow int
	var obj, template Usage

	con := db.CreateCon()

	template.Username = username
	template.Download = 0
	template.Upload = 0
	template.Uptime = 0

	verifyStatement := "SELECT COUNT(*) FROM radacct WHERE username = ?"
	verifyResult, err := con.Query(verifyStatement, username)
	if err != nil {
		return res, err
	}

	defer verifyResult.Close()
	for verifyResult.Next() {
		err = verifyResult.Scan(&selectedRow)
		if err != nil {
			return res, err
		}
	}

	if selectedRow > 0 {
		usageStatement := "SELECT username, SUM(acctsessiontime), SUM(acctinputoctets), SUM(acctoutputoctets) FROM radacct WHERE username = ?"

		result, err := con.Query(usageStatement, username)
		if err != nil {
			return res, err
		}

		defer result.Close()

		for result.Next() {
			err = result.Scan(&obj.Username, &obj.Uptime, &obj.Download, &obj.Upload)
			if err != nil {
				return res, nil
			}
		}

		res.Status = http.StatusOK
		res.Message = "Success"
		res.Data = obj
	} else {
		res.Status = http.StatusOK
		res.Message = "Success"
		res.Data = template
	}

	return res, nil
}
