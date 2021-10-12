package infra

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/guricerin/stop-now-smoking/util"
	_ "github.com/lib/pq"
)

type PostgresDriver struct {
	conn *sql.DB
}

func NewPostgresDriver(cfg *util.Config) (*PostgresDriver, error) {
	util.Ilog.Printf("dsn: %s", cfg.DbUrl)
	conn, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		err = fmt.Errorf("db open error: %s", err.Error())
		return nil, err
	}
	if err = tryPing(conn); err != nil {
		err = fmt.Errorf("db ping error: %s", err.Error())
		return nil, err
	}

	res := new(PostgresDriver)
	res.conn = conn
	return res, nil
}

func tryPing(conn *sql.DB) (err error) {
	for i := 0; i < 3; i++ {
		err = conn.Ping()
		if err == nil {
			return
		}
		time.Sleep(5 * time.Second)
	}
	return
}

func (d *PostgresDriver) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	rows, err = d.conn.Query(query, args...)
	return
}

func (d *PostgresDriver) QueryRow(query string, args ...interface{}) (row *sql.Row) {
	row = d.conn.QueryRow(query, args...)
	return
}

func (d *PostgresDriver) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	res, err = d.conn.Exec(query, args...)
	return
}

func (d *PostgresDriver) Prepare(statement string) (stmt *sql.Stmt, err error) {
	stmt, err = d.conn.Prepare(statement)
	return
}

func (d *PostgresDriver) Close() {
	util.Ilog.Printf("db driver closing ...")
	d.conn.Close()
}
