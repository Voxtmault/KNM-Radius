package models

import (
	"CNM_Radius/db"
	"net/http"
)

type UserProfile struct {
	Username  string `json:Username`
	Groupname string `json:Groupname `
}

func GetUserProfile() (Response, error) {
	var obj UserProfile
	var arrObj []UserProfile
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT username,groupname from radusergroup"

	result, err := con.Query(sqlStatement)

	defer result.Close()

	if err != nil {
		return res, err
	}

	for result.Next() {
		err = result.Scan(&obj.Username, &obj.Groupname)

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
