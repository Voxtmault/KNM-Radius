package models

import (
	"KNM-Radius/db"
	"net/http"
)

type Profile struct {
	ProfileName string `json:username`
	Attribute   string `json:attribute`
	Operator    string `json:operator`
	Value       string `json:value`
}

func CreateNewProfile(profileName string, profileAttribute string, profileOperator string, profileValue int) (Response, error) {
	//New Profile Method
	//1. Add to RadGroupReply
	//2. For now Attribute(S) = Ascend-Data-Rate & Ascend-Xmit-Rate
	var res Response

	radiusConnection := db.CreateRadiusCon()

	newProfileStatement := "INSERT INTO radgroupreply (groupname, attribute, op, value) VALUES(?,?,?,?)"

	stmt, err := radiusConnection.Prepare(newProfileStatement)
	if err != nil {
		return res, err
	}

	newProfileResult, err := stmt.Exec(profileName, profileAttribute, profileOperator, profileValue)
	if err != nil {
		return res, err
	}

	hitResult, _ := newProfileResult.LastInsertId()
	if hitResult != 0 {
		res.Message = "Success Creating New Profile"
		res.Data = map[string]int64{
			"New Profile ID: ": hitResult,
		}
	} else {
		res.Message = "Failed Creating New Profile"
		res.Data = ""
	}

	res.Status = http.StatusOK

	return res, nil
}

func DeleteProfile(profileName string) (Response, error) {
	//New Profile Method
	//1. Delete From RadGroupReply
	var res Response

	radiusConnection := db.CreateRadiusCon()

	profileStatement := "DELETE FROM radgroupreply WHERE groupname = ?"

	profileSTMT, err := radiusConnection.Prepare(profileStatement)
	if err != nil {
		return res, err
	}

	_, err = profileSTMT.Exec(profileName)

	res.Status = http.StatusOK
	res.Message = "Success Deleting Profile"
	res.Data = map[string]string{
		"Profile Deleted": profileName,
	}

	return res, nil
}
func GetProfileHotspot() (Response, error) {
	var obj Profile
	var arrObj []Profile
	var res Response

	radiusConnection := db.CreateRadiusCon()

	sqlStatement := "SELECT distinct groupname from radgroupreply"

	result, err := radiusConnection.Query(sqlStatement)

	defer result.Close()

	if err != nil {
		return res, err
	}

	for result.Next() {
		err = result.Scan(&obj.ProfileName)

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
