package models

import (
	"KNM-Radius/db"
	"encoding/json"
	"net/http"
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

func GetAllUsers() (Response, error) {
	var obj Users
	var arrObj []Users
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT username from radcheck"

	result, err := con.Query(sqlStatement)

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

func CreateNewUser(userCredentials string, maxUse int, expiredate int, profile string, sessionTimeout int, idleTimeout int, userCreator string) (Response, error) {
	//How to Create New Users
	//1. Add to RadCheck (For Password & Max user per account)
	//2. Add to RadReply (For Session-Timeout & Idle-Timout)
	//3. Add to RadUserGroup (For Profile)
	var res Response

	con := db.CreateCon()

	radCheckStatement := "INSERT INTO radcheck (username, attribute, op, value) VALUES(?,?,?,?)"
	radReplyStatement := "INSERT INTO radreply(username, attribute, op, value) VALUES(?,?,?,?)"
	radProfileStatement := "INSERT INTO radusergroup(username, groupname,priority) VALUES(?,?,?)"

	//1
	checkSTMT, err := con.Prepare(radCheckStatement)

	_, err = checkSTMT.Exec(userCredentials, "Cleartext-Password", ":=", userCredentials)
	if err != nil {
		return res, err
	}

	_, err = checkSTMT.Exec(userCredentials, "Simultaneous-Use", ":=", maxUse)
	if err != nil {
		return res, err
	}
	defer checkSTMT.Close()

	//2
	replySTMT, err := con.Prepare(radReplyStatement)

	_, err = replySTMT.Exec(userCredentials, "Session-Timeout", ":=", sessionTimeout)
	if err != nil {
		return res, err
	}

	_, err = replySTMT.Exec(userCredentials, "Idle-Timeout", ":=", idleTimeout)
	if err != nil {
		return res, err
	}
	defer replySTMT.Close()

	//3
	groupSTMT, err := con.Prepare(radProfileStatement)

	_, err = groupSTMT.Exec(userCredentials, profile, 0)
	if err != nil {
		return res, err
	}
	defer groupSTMT.Close()

	res.Status = http.StatusOK
	res.Message = "Success Creating User: " + userCredentials
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
