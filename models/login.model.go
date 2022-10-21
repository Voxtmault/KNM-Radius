package models

import (
	"KNM-Radius/db"
	"fmt"
	"net/http"
)

func Authenticate(email string, password string) (Response, error) {
	var res Response
	var hit int

	fmt.Println("Email: ", email)
	fmt.Println("Password: ", password)
	adminConnection := db.CreateAdminCon()

	sqlStatement := "SELECT COUNT(*) FROM admin WHERE email = ? AND password = ?"

	result, err := adminConnection.Query(sqlStatement, email, password)

	defer result.Close()

	if err != nil {
		return res, err
	}

	for result.Next() {
		err = result.Scan(&hit)
		if err != nil {
			return res, err
		}
	}

	res.Status = http.StatusOK
	res.Data = "What are you looking for, eh ?"

	if hit > 0 {
		res.Message = "Login Success"
	} else {
		res.Message = "Login Failed"
	}

	return res, nil

}
