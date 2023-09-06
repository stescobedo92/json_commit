package main

import (
	"log"

	"github.com/stescobedo92/json_commit/internal/server"
)

func main() {
	serv := server.NewHttpServer(":8080")
	log.Fatal(serv.ListenAndServe())
}
