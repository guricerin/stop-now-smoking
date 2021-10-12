package main

import (
	"os"

	"github.com/guricerin/stop-now-smoking/infra"
	"github.com/guricerin/stop-now-smoking/server"
	. "github.com/guricerin/stop-now-smoking/util"
)

func main() {
	Ilog.Println("setup config ...")
	cfg := Config{
		DbUrl:      os.Getenv("DATABASE_URL"),
		ServerHost: "",
		ServerPort: os.Getenv("PORT"),
	}
	Dlog.Printf("cfg: %v", cfg)

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
