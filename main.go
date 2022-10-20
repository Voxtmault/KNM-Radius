package main

import (
	"KNM-Radius/db"
	"KNM-Radius/routes"
)

func main() {
	db.Init()

	e := routes.Init()

	e.Logger.Fatal(e.Start(":38900"))
}
