package main

import (
	"github.com/guricerin/stop-now-smoking/server"
	u "github.com/guricerin/stop-now-smoking/util"
)

func main() {
	u.Ilog.Println("setup server ...")
	server := server.NewServer()
	err := server.Run()
	if err != nil {
		u.Elog.Fatalf("server.Run error: %v", err)
	}
}
