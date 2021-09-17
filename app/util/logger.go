package util

import (
	"log"
	"os"
)

var (
	Dlog *log.Logger
	Ilog *log.Logger
	Wlog *log.Logger
	Elog *log.Logger
)

func init() {
	Dlog = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Llongfile)
	Ilog = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	Wlog = log.New(os.Stderr, "[WARN] ", log.Ldate|log.Ltime)
	Elog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
}
