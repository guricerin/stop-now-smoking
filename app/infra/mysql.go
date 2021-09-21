package infra

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guricerin/stop-now-smoking/util"
)

type MySqlDriver struct {
	conn *sql.DB
}

func NewMySqlDriver(cfg *util.Config) (*MySqlDriver, error) {
	user := cfg.DbUser
	password := cfg.DbPassword
	protocol := cfg.DbProtocol
	dbName := cfg.DbName
	option := cfg.DbConnOption
	dsn := fmt.Sprintf("%s:%s@%s/%s?%s", user, password, protocol, dbName, option)
	util.Ilog.Printf("dsn: %s", dsn)
	conn, err := sql.Open("mysql", dsn)
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
	for i := 0; i < 3; i++ {
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

func (d *MySqlDriver) Close() {
	util.Ilog.Printf("db driver closing ...")
	d.conn.Close()
}
