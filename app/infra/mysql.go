package infra

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlDriver struct {
	conn *sql.DB
}

func NewMySqlDriver() (*MySqlDriver, error) {
	user := "sns"
	password := "sns123456789"
	protocol := "tcp(db:3306)"
	dbName := "sns_db"
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, password, protocol, dbName)
	conn, err := sql.Open("mysql", dsn+"?charset=utf8&parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		err = fmt.Errorf("db open error: %s", err.Error())
		return nil, err
	}
	if err = tryPing(conn); err != nil {
		err = fmt.Errorf("db ping error: %s", err.Error())
		return nil, err
	}

	res := new(MySqlDriver)
	res.conn = conn
	return res, nil
}

func tryPing(conn *sql.DB) (err error) {
	for i := 0; i < 10; i++ {
		err = conn.Ping()
		if err == nil {
			return
		}
		time.Sleep(5 * time.Second)
	}
	return
}

func (d *MySqlDriver) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	rows, err = d.conn.Query(query, args...)
	return
}

func (d *MySqlDriver) QueryRow(query string, args ...interface{}) (row *sql.Row) {
	row = d.conn.QueryRow(query, args...)
	return
}

func (d *MySqlDriver) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	res, err = d.conn.Exec(query, args...)
	return
}

func (d *MySqlDriver) Prepare(statement string) (stmt *sql.Stmt, err error) {
	stmt, err = d.conn.Prepare(statement)
	return
}
