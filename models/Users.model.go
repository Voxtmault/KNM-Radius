package models

import (
	"KNM-Radius/db"
	"encoding/json"
	"net/http"
	"time"
)

type Users struct {
	Username string `json:username`
}

type Temperatures []ActiveUsers

type ActiveUsers struct {
	ID               string `json:".id"`
	Server           string `json:"server"`
	User             string `json:"user"`
	Address          string `json:"address"`
	MACAddress       string `json:"mac-address"`
	LoginBy          string `json:"login-by"`
	Uptime           string `json:"uptime"`
	IdleTime         string `json:"idle-time"`
	KeepaliveTimeout string `json:"keepalive-timeout"`
	BytesIn          string `json:"bytes-in"`
	BytesOut         string `json:"bytes-out"`
	PacketsIn        string `json:"packets-in"`
	PacketsOut       string `json:"packets-out"`
	Radius           string `json:"radius"`
}

type User struct {
	MAC           string `json:mac`
	AccountName   string `json:name`
	DeviceInfo    string `json:device`
	ExpireDate    int    `json:expireDate`
	SessionTimout int    `json:sessionTimeout`
	IdleTimeout   int    `json:idleTimeout`
	ProfileName   string `json:profile`
}

func GetAllUsers() (Response, error) {
	var obj Users
	var arrObj []Users
	var res Response

	radiusConnection := db.CreateRadiusCon()

	sqlStatement := "SELECT username from radcheck"

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

func GetUser(userCredentials string) (Response, error) {
	//How to Get Users Info
	//1. Get MAC Address and ExpiredDate from RadCheck
	//2. Get SessionTimeout and IdleTimeout from RadReply
	//3. Get ProfileName from RadUserGroup
	//4. Get AccountName and DeviceInfo from UserInfo
	var res Response
	//var obj User

	//radiusConnecton := db.CreateRadiusCon()

	//radCheckStatement := "SELECT username, attribute, value FROM radcheck WHERE username = ?"
	//radReplyStatement := "SELECT username, attribute, value FROM radreply WHERE username = ?"

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
func UnmarshalTemperatures(data []byte) (Temperatures, error) {
	var r Temperatures
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Temperatures) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
