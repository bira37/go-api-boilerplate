package main

import (
	"log"

	"bira.io/template/infra"
	"bira.io/template/server"
)

func main() {
	router := server.SetupServer("file://infra/migrations", infra.Config.SqlDbName)

	if err := router.Run(":7000"); err != nil {
		log.Fatalln(err.Error())
	}
}
