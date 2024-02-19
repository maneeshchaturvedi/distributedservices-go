package main

import (
	"log"

	"github.com/maneeshchaturvedi/log/internal/server"
)

func main() {
	addr := "localhost:8080"
	srv := server.NewHttpServer(&addr)
	log.Printf("server started at %s", addr)
	log.Fatal(srv.ListenAndServe())
}
