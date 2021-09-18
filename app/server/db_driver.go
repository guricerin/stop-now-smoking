package server

import "database/sql"

type DbDriver interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
}
