package models

import (
	"CNM_Radius/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

func CreateNewUser(username string, sid string, max_device int, email string, mobilephone int, address string, planname string, paymentMethod string, ordertaker string, expiredate int, profile string) (Response, error) {
	//How to Create New Users
	//1. Add to RadUserGroup
	//2. Add to UserInfo
	//3. Add to UserBillInfo
	//4. Add to RadCheck
	var res Response
	var createdID int

	var templateUsername, temp string
	var hitResult int

	templateUsername = sid + username
	temp = ""
	hitResult = 0
	con := db.CreateCon()

	verifyPlanNameStatement := "SELECT planname FROM billing_plans WHERE planName = ?"
	verifyResult, err := con.Query(verifyPlanNameStatement, planname)
	if err != nil {
		return res, err
	}

	for verifyResult.Next() {
		err = verifyResult.Scan(&temp)
		if err != nil {
			return res, err
		}

		if len(temp) > 0 {
			//fmt.Println(temp)
			hitResult++
		}
	}

	temp = ""

	verifyProfileNameStatement := "SELECT groupname FROM radgroupreply WHERE groupname = ?"
	verifyProfileResult, err := con.Query(verifyProfileNameStatement, profile)
	if err != nil {
		return res, err
	}

	for verifyProfileResult.Next() {
		err = verifyProfileResult.Scan(&temp)
		if err != nil {
			return res, err
		}

		if len(temp) > 0 {
			//fmt.Println(temp)
			hitResult++
		}
	}

	if hitResult >= 2 {
		userGroupStatement := "INSERT INTO radusergroup(username, groupname, priority) VALUES(?,?,?)"
		userInfoStatement := "INSERT INTO userinfo(username, firstname, email, mobilephone, address, creationdate, creationby) VALUES(?,?,?,?,?,?,?)"
		userBillInfoStatement := "INSERT INTO userbillinfo(username, planName, email, phone, address, paymentmethod, ordertaker, creationdate, creationby, hotspot_id) VALUES(?,?,?,?,?,?,?,?,?,?)"
		radCheckStatement := "INSERT INTO radcheck (username, attribute, op, value) VALUES(?,?,?,?)"

		stmt, err := con.Prepare(radCheckStatement)

		for i := 0; i < max_device; i++ {
			if i == 0 {
				result, err := stmt.Exec(templateUsername, "Cleartext-Password", ":=", templateUsername)
				if err != nil {
					return res, err
				}

				_, err = stmt.Exec(templateUsername, "Expiration", ":=", time.Now().Add(time.Hour*24*time.Duration(expiredate)).Format("02 Jan 2006"))
				if err != nil {
					return res, err
				}

				_, err = con.Query(userGroupStatement, templateUsername, profile, 0)
				if err != nil {
					return res, err
				}

				affected, _ := result.RowsAffected()
				if affected > 0 {
					createdID++
				}
			} else {
				result, err := stmt.Exec(templateUsername+"-"+strconv.Itoa(i), "Cleartext-Password", ":=", templateUsername+"-"+strconv.Itoa(i))
				if err != nil {
					return res, err
				}

				_, err = stmt.Exec(templateUsername+"-"+strconv.Itoa(i), "Expiration", ":=", time.Now().Add(time.Hour*24*time.Duration(expiredate)).Format("02 Jan 2006"))
				if err != nil {
					return res, err
				}

				_, err = con.Query(userGroupStatement, templateUsername+"-"+strconv.Itoa(i), profile, 0)
				if err != nil {
					return res, err
				}

				affected, _ := result.RowsAffected()
				if affected > 0 {
					createdID++
				}
			}

			_, err = con.Query(userInfoStatement, email, username, email, mobilephone, address, time.Now().Format("2006-01-02 15:04"), "CNM-System")
			if err != nil {
				return res, err
			}

			_, err = con.Query(userBillInfoStatement, username, planname, email, mobilephone, address, paymentMethod, ordertaker, time.Now().Format("2006-01-02 15:04"), "CNM-System", sid)
			if err != nil {
				return res, err
			}
		}

		if err != nil {
			return res, err
		}

		res.Status = http.StatusOK
		res.Message = "Success"
		res.Data = map[string]int{
			"User Created": createdID,
		}
	} else {
		res.Status = http.StatusInternalServerError
		res.Message = "Plan name or Profile not found ! Please check your information again !"
		res.Data = ""
	}

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

func AddChildToUser() (Response, error) {
	var res Response

	result, err := http.Get(db.GetActiveUserIP())
	if err != nil {
		return res, err
	}

	defer result.Body.Close()

	body, _ := ioutil.ReadAll(result.Body)
	smt, err := UnmarshalTemperatures(body)

	//PrintChild(smt)

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = smt
	return res, nil
}

func PrintChild(data Temperatures) {
	for i := 0; i < len(data); i++ {
		fmt.Println("User: ", data[i].User)
	}
}

func GetActiveUsers() (Response, error) {
	var res Response

	resp, err := http.Get("http://10.10.10.232/api-mikrotik-main/hotspotactive.php")
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	result, err := UnmarshalTemperatures([]byte(body))

	res.Status = http.StatusOK
	res.Message = "Sucess Getting All Active Users"
	res.Data = result

	return res, nil
}

func GetAvailableUser(sid string, username string) (Response, error) {
	var res Response
	var obj Users
	var arrObj []Users
	var templateUsername string

	con := db.CreateCon()
	templateUsername = sid + username

	sqlStatement := "SELECT username FROM radcheck WHERE username LIKE ? AND attribute = 'Cleartext-Password'"

	result, err := con.Query(sqlStatement, templateUsername+"%")

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
