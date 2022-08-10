package main

import (
	"log"
	"win-a-thon/database"
	"win-a-thon/routes"
)

var err error

func main() {
	database.DB, err = database.GetDatabase()

	if err != nil {
		log.Fatal("Database creation failed.", err)
	}

	r, err := routes.Setup()
	if err != nil {
		log.Fatal("Token maker creation failed.", err)
	}

	r.Run()
}
