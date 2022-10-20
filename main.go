package main

import (
	"CNM_Radius/db"
	"CNM_Radius/routes"
)

func main() {
	db.Init()

	e := routes.Init()

	e.Logger.Fatal(e.Start(":38700"))
}
