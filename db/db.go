package db

import (
	"CNM_Radius/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

var activeUserIP = "http://10.10.10.232/api-mikrotik-main/hotspotactive.php"
var PhoenixDB = "phoenix.crossnet.co.id:38600"

func Init() {
	conf := config.GetConfig()

	//DSN Full Format
	//username:password@protocol(address)/dbname?param=value
	connectionString := conf.DB_USERNAME + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME

	db, err = sql.Open("mysql", connectionString)

	if err != nil {
		panic(err)
	}

	err := db.Ping()

	if err != nil {
		panic("DSN Error")
	}
}

func CreateCon() *sql.DB {
	return db
}

func GetActiveUserIP() string {
	return activeUserIP
}

func GetApplicationIP() string {
	return PhoenixDB
}
