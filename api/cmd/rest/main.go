package main

import (
	"log"

	"github.com/bira37/go-rest-api/api/server"
)

func main() {
	router := server.SetupRestServer()

	if err := router.Run(":7000"); err != nil {
		log.Fatalln(err.Error())
	}
}
