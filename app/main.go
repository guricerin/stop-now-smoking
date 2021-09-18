package main

import (
	"github.com/guricerin/stop-now-smoking/infra"
	"github.com/guricerin/stop-now-smoking/server"
	u "github.com/guricerin/stop-now-smoking/util"
)

func main() {
	u.Ilog.Println("connecting db ...")
	db, err := infra.NewMySqlDriver()
	if err != nil {
		u.Elog.Fatalf("mysql driver error: %v", err)
	}

	u.Ilog.Println("setup server ...")
	server := server.NewServer(db)
	err = server.Run()
	if err != nil {
		u.Elog.Fatalf("server.Run error: %v", err)
	}
}
