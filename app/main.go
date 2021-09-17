package main

import "github.com/guricerin/stop-now-smoking/server"

func main() {
	server := server.NewServer()
	server.Run()
}
