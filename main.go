package main

import (
	"KNM-Radius/db"
	"KNM-Radius/routes"
)

func main() {
	db.Init(0)
	db.Init(1)

	e := routes.Init()

	e.Logger.Fatal(e.Start(":38900"))
}
