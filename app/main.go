package main

import (
	"github.com/guricerin/stop-now-smoking/infra"
	"github.com/guricerin/stop-now-smoking/server"
	. "github.com/guricerin/stop-now-smoking/util"
)

func main() {
	Ilog.Println("connecting db ...")
	db, err := infra.NewMySqlDriver()
	if err != nil {
		Elog.Fatalf("mysql driver error: %v", err)
	}

	server := server.NewServer(db)
	Ilog.Println("starting server ...")
	err = server.Run()
	if err != nil {
		Elog.Fatalf("server.Run error: %v", err)
	}
}
