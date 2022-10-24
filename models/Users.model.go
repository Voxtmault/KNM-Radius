package models

import (
	"KNM-Radius/db"
	"net/http"
	"time"
)

type Users struct {
	Username string `json:username`
}

type User struct {
	MAC           string `json:mac`
	AccountName   string `json:name`
	DeviceInfo    string `json:device`
	ExpireDate    string `json:expireDate`
	SessionTimout string `json:sessionTimeout`
	IdleTimeout   string `json:idleTimeout`
	ProfileName   string `json:profile`
	MaxDevice     string `json:max_device`
}

type UserProfile struct {
	Username  string `json:Username`
	Groupname string `json:Groupname `
}

func GetAllUsers() (Response, error) {
	var obj Users
	var arrObj []Users
	var res Response

	radiusConnection := db.CreateRadiusCon()

	sqlStatement := "SELECT distinct username from radcheck"

	result, err := radiusConnection.Query(sqlStatement)

	defer result.Close()

	if err != nil {
		return res, err
	}

	for result.Next() {
		err = result.Scan(&obj.Username)

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

func GetUserInfo(userCredentials string) (Response, error) {
	//How to Get Users Info
	//1. Get MAC Address and ExpiredDate from RadCheck (MAC Addr can also be assigned from param)
	//2. Get SessionTimeout and IdleTimeout from RadReply
	//3. Get ProfileName from RadUserGroup
	//4. Get AccountName and DeviceInfo from UserInfo
	var res Response
	var obj User

	radiusConnecton := db.CreateRadiusCon()

	//Order will be : 1. Cleartext Password, Expiration, Simultaneous Use
	radCheckStatement := "SELECT (SELECT coalesce(value, '') from radcheck WHERE username = ? AND attribute = \"Cleartext-Password\") as Password, coalesce((SELECT value from radcheck WHERE username = ? AND attribute = \"Expiration\"), \"0\") as Expiration, (SELECT coalesce(value, 1) from radcheck WHERE username = ? AND attribute = \"Simultaneous-Use\") as \"Max Use\" from radcheck where username = ? ORDER BY id"

	radReplyStatement := "SELECT (SELECT coalesce(value, 0) from radreply WHERE username = ? AND attribute = \"Session-Timeout\") as \"Session Timeout\", (SELECT coalesce(value, 0) from radreply WHERE username = ? AND attribute = \"Idle-Timeout\") as \"Idle Timeout\" FROM radreply where username = ? ORDER BY id"
	radGroupStatement := "SELECT groupname FROM radusergroup WHERE username = ?"
	userInfoStatement := "SELECT coalesce(firstname, \"Unavailable\") , coalesce(lastname, \"Unavailable\") FROM userinfo WHERE username = ?"

	//1
	radCheckResult, err := radiusConnecton.Query(radCheckStatement, userCredentials, userCredentials, userCredentials, userCredentials)
	if err != nil {
		return res, err
	}

	for radCheckResult.Next() {
		err = radCheckResult.Scan(&obj.MAC, &obj.ExpireDate, &obj.MaxDevice)
		if err != nil {
			return res, err
		}
		break
	}

	if obj.ExpireDate == "0" {
		obj.ExpireDate = "Forever"
	}

	obj.MaxDevice += " Device(s)"

	//2
	radReplyResult, err := radiusConnecton.Query(radReplyStatement, userCredentials, userCredentials, userCredentials)
	if err != nil {
		return res, err
	}

	for radReplyResult.Next() {
		err = radReplyResult.Scan(&obj.SessionTimout, &obj.IdleTimeout)
		if err != nil {
			return res, err
		}
		break
	}

	if obj.SessionTimout == "" {
		obj.SessionTimout = "No Session Timeouts"
	} else {
		obj.SessionTimout += " Second(s)"
	}

	if obj.IdleTimeout == "" {
		obj.IdleTimeout = "No Idle Timeout"
	} else {
		obj.IdleTimeout += " Second(s)"
	}

	//3
	radGroupResult, err := radiusConnecton.Query(radGroupStatement, userCredentials)
	if err != nil {
		return res, err
	}

	for radGroupResult.Next() {
		err = radGroupResult.Scan(&obj.ProfileName)
		if err != nil {
			return res, err
		}

		break
	}

	if obj.ProfileName == "" {
		obj.ProfileName = "No Profile Selected"
	}

	//4
	userInfoResult, err := radiusConnecton.Query(userInfoStatement, userCredentials)
	if err != nil {
		return res, err
	}

	for userInfoResult.Next() {
		err = userInfoResult.Scan(&obj.AccountName, &obj.DeviceInfo)
	}

	if obj.AccountName == "" {
		obj.AccountName = "Data Unavailable"
	}

	if obj.DeviceInfo == "" {
		obj.DeviceInfo = "Data Unavailable"
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = obj

	return res, nil
}

func GetAllUsersInfo() (Response, error) {
	var obj User
	var arrObj []User
	var res Response

	radiusConnection := db.CreateRadiusCon()

	sqlStatement := "SELECT coalesce(username, ''), coalesce(firstname, \"Unavailable\"), coalesce(lastname, \"Unavailable\") from userinfo order by id"

	result, err := radiusConnection.Query(sqlStatement)

	defer result.Close()

	if err != nil {
		return res, err
	}

	for result.Next() {
		err = result.Scan(&obj.MAC, &obj.AccountName, &obj.DeviceInfo)

		if err != nil {
			return res, nil
		}

		if len(obj.AccountName) == 0 {
			obj.AccountName = "Unavailable"
		}

		if len(obj.DeviceInfo) == 0 {
			obj.DeviceInfo = "Unavailable"
		}
		arrObj = append(arrObj, obj)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrObj

	return res, nil
}

func CreateNewUser(userCredentials string, expireDate int, profile string, sessionTimeout int, idleTimeout int, deviceOwner string, deviceInfo string) (Response, error) {
	//How to Create New Users
	//1. Add to RadCheck (For Password & Max user per account)
	//2. Add to RadReply (For Session-Timeout & Idle-Timout)
	//3. Add to RadUserGroup (For Internet Profile)
	//4. Add to UserInfo (For User Personal Information)
	var res Response

	radiusConnection := db.CreateRadiusCon()

	radCheckStatement := "INSERT INTO radcheck (username, attribute, op, value) VALUES(?,?,?,?)"
	radReplyStatement := "INSERT INTO radreply(username, attribute, op, value) VALUES(?,?,?,?)"
	radProfileStatement := "INSERT INTO radusergroup(username, groupname,priority) VALUES(?,?,?)"
	radUserStatement := "INSERT INTO userinfo(username, firstname, lastname, creationdate, creationby, updatedate) VALUES(?,?,?,?,?,?)"

	//1
	checkSTMT, err := radiusConnection.Prepare(radCheckStatement)

	_, err = checkSTMT.Exec(userCredentials, "Cleartext-Password", ":=", userCredentials)
	if err != nil {
		return res, err
	}

	_, err = checkSTMT.Exec(userCredentials, "Simultaneous-Use", ":=", 1)
	if err != nil {
		return res, err
	}

	if expireDate > 0 {
		_, err := checkSTMT.Exec(userCredentials, "Expiration", ":=", time.Now().Add(time.Hour*24*time.Duration(expireDate)).Format("02 Jan 2006"))
		if err != nil {
			return res, err
		}
	}
	defer checkSTMT.Close()

	//2
	replySTMT, err := radiusConnection.Prepare(radReplyStatement)

	if sessionTimeout > 0 {
		_, err = replySTMT.Exec(userCredentials, "Session-Timeout", ":=", sessionTimeout)
		if err != nil {
			return res, err
		}
	}

	if idleTimeout > 0 {
		_, err = replySTMT.Exec(userCredentials, "Idle-Timeout", ":=", idleTimeout)
		if err != nil {
			return res, err
		}
	}
	defer replySTMT.Close()

	//3
	groupSTMT, err := radiusConnection.Prepare(radProfileStatement)

	if len(profile) > 0 {
		_, err = groupSTMT.Exec(userCredentials, profile, 0)
		if err != nil {
			return res, err
		}
	}
	defer groupSTMT.Close()

	//4
	userSTMT, err := radiusConnection.Prepare(radUserStatement)

	_, err = userSTMT.Exec(userCredentials, deviceOwner, deviceInfo, time.Now().Format("2006-01-02 15:04:05"), "System", nil)
	if err != nil {
		return res, err
	}
	defer userSTMT.Close()

	res.Status = http.StatusOK
	res.Message = "Success Creating User: " + userCredentials
	res.Data = ""

	return res, nil
}

func DeleteUser(userCredentials string) (Response, error) {
	//How to Delete New Users
	//1. Delete from RadCheck (For Password & Max user per account)
	//2. Delete from RadReply (For Session-Timeout & Idle-Timout)
	//3. Delete from RadUserGroup (For Internet Profile)
	//4. Delete from UserInfo (For User Personal Information)

	var res Response

	radiusConnection := db.CreateRadiusCon()

	radCheckStatement := "DELETE FROM radcheck WHERE username = ?"
	radReplyStatement := "DELETE FROM radreply WHERE username = ?"
	radProfileStatement := "DELETE FROM radusergroup WHERE username = ?"
	radUserStatement := "DELETE FROM userinfo WHERE username = ?"

	//1
	_, err := radiusConnection.Query(radCheckStatement, userCredentials)
	if err != nil {
		return res, err
	}

	//2
	_, err = radiusConnection.Query(radReplyStatement, userCredentials)
	if err != nil {
		return res, err
	}

	//3
	_, err = radiusConnection.Query(radProfileStatement, userCredentials)
	if err != nil {
		return res, err
	}

	//4
	_, err = radiusConnection.Query(radUserStatement, userCredentials)
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success Deleting User: " + userCredentials
	res.Data = ""

	return res, nil
}

func EditUser(userCredentials string, expireDate int, profile string, sessionTimeout int, idleTimeout int, deviceOwner string, deviceInfo string) (Response, error) {
	var res Response

	DeleteUser(userCredentials)

	CreateNewUser(userCredentials, expireDate, profile, sessionTimeout, idleTimeout, deviceOwner, deviceInfo)

	res.Status = http.StatusOK
	res.Message = "Success Editing User : " + userCredentials
	res.Data = ""
	return res, nil
}

func GetUserProfile() (Response, error) {
	var obj UserProfile
	var arrObj []UserProfile
	var res Response

	con := db.CreateRadiusCon()

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
