package tools

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gorp "gopkg.in/gorp.v1"
)

type DbConn struct {
	Name  string
	DbMap *gorp.DbMap
}

var dbConnMap map[string]*DbConn
var defaultDbName string

func init() {
	dbConnMap = make(map[string]*DbConn)
	defaultDbName = ""
}

func GetDefDb() *DbConn {
	return GetDb(defaultDbName)
}

func GetDb(name string) *DbConn {
	return dbConnMap[name]
}

func CreateDbConn(database, host string, port int, user, password string, maxConn, maxIdle, lifeTime int) (err error) {
	if _, exist := dbConnMap[database]; exist {
		return
	}

	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database)

	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Duration(lifeTime) * time.Second)

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}

	if defaultDbName == "" {
		defaultDbName = database
	}

	dbConnMap[database] = &DbConn{
		Name:  database,
		DbMap: dbmap,
	}

	return

}
