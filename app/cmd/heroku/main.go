package main

import (
	"os"

	"github.com/guricerin/stop-now-smoking/infra"
	"github.com/guricerin/stop-now-smoking/server"
	. "github.com/guricerin/stop-now-smoking/util"
)

func main() {
	Ilog.Println("loading config file ...")
	cfg, err := LoadConfig("config.json")
	if err != nil {
		Elog.Fatalf("load config error: %v", err)
	}
	Dlog.Printf("cfg: %v", cfg)
	cfg.DbUrl = os.Getenv("DATABASE_URL")

	Ilog.Println("connecting db ...")
	db, err := infra.NewPostgresDriver(&cfg)
	if err != nil {
		Elog.Fatalf("db driver error: %v", err)
	}
	defer db.Close()

	server := server.NewServer(&cfg, db)
	Ilog.Println("starting server ...")
	err = server.Run()
	if err != nil {
		Elog.Fatalf("server.Run error: %v", err)
	}
}
