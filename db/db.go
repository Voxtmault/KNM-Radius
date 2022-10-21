package db

import (
	"KNM-Radius/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var radiusDB, adminDB *sql.DB
var err error

func Init(connectTo int) {
	var connectionString string

	conf := config.GetConfig()

	//DSN Full Format
	//username:password@protocol(address)/dbname?param=value

	if connectTo == 0 {
		//Connect to Radius (Authentication Server)
		connectionString = conf.RADIUS_DB_USERNAME + ":" + conf.RADIUS_DB_PASSWORD + "@tcp(" + conf.RADIUS_DB_HOST + ":" + conf.RADIUS_DB_PORT + ")/" + conf.RADIUS_DB_NAME
		radiusDB, err = sql.Open("mysql", connectionString)

		if err != nil {
			panic(err)
		}

		err := radiusDB.Ping()

		if err != nil {
			panic("DSN Error")
		}
	} else if connectTo == 1 {
		//Connect to Phoenix (Admin Server)
		connectionString = conf.PHOENIX_DB_USERNAME + ":" + conf.PHOENIX_DB_PASSWORD + "@tcp(" + conf.PHOENIX_DB_HOST + ":" + conf.PHOENIX_DB_PORT + ")/" + conf.PHOENIX_DB_NAME
		adminDB, err = sql.Open("mysql", connectionString)

		if err != nil {
			panic(err)
		}

		err := adminDB.Ping()

		if err != nil {
			panic("DSN Error")
		}
	} else {
		panic(err)
	}
}

func CreateRadiusCon() *sql.DB {
	return radiusDB
}

func CreateAdminCon() *sql.DB {
	return adminDB
}
