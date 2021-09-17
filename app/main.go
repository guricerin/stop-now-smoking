package main

import (
	"log"

	"github.com/guricerin/stop-now-smoking/server"
)

func main() {
	log.Println("server setup ...")
	server := server.NewServer()
	server.Run()
}
