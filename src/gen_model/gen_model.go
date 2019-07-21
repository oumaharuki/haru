//
// 生成数据库中和表对应的定义
// 打印到标准输出, 生成一个model文件
//

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

type COLUMNS struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}

func panicCheck(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func initDb(dbSource string) *gorp.DbMap {
	db, err := sql.Open("mysql", dbSource)
	panicCheck(err)

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}

	//dbmap.AddTableWithName(FinanceEvent{}, "finance_event").SetKeys(false, "Id")

	return dbmap
}

func capital(str string) string {
	return strings.ToUpper(str[0:1]) + str[1:]
}

func toExportName(inStr string) string {
	str := inStr
	for {
		pos := strings.Index(str, "_")
		if pos >= 0 && pos < len(str)-1 {
			digitPos := pos + 1
			str = str[0:pos] + strings.ToUpper(str[digitPos:digitPos+1]) + str[digitPos+1:]
		} else {
			break
		}
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}

func main() {
	var dbSource string
	flag.StringVar(&dbSource, "dbSource", "", "user:passwd@tcp(127.0.0.1:3306)/db")
	//flag.StringVar(&table, "table", "", "DB table")

	flag.Parse()

	dbmap := initDb(dbSource)
	defer dbmap.Db.Close()

	type TIB struct {
		Name string `db:"name"`
	}
	var tibs []TIB
	pos := strings.LastIndex(dbSource, "/")
	dbName := dbSource[(pos + 1):]
	//fmt.Println("dbName:", dbName)
	_, err := dbmap.Select(&tibs, "SELECT TABLE_NAME AS `name` FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '"+dbName+"'")
	if err != nil {
		panic(err)
	}

	//TABLE_NAME := "user"
	fmt.Println("package model")

	bImportSql := false

	code := ""

	for _, tib := range tibs {

		cItems, err := dbmap.Select(COLUMNS{}, "show columns from `"+tib.Name+"`")
		panicCheck(err)

		code += fmt.Sprintf("type %s struct{\n", toExportName(tib.Name))
		for i, _ := range cItems {
			cItem := cItems[i].(*COLUMNS)
			name := toExportName(cItem.Field)
			fieldType := ""
			if strings.Index(cItem.Type, "int(") == 0 {
				fieldType = "int"
			} else if strings.Index(cItem.Type, "bigint(") == 0 {
				fieldType = "int64"
			} else if strings.Index(cItem.Type, "varchar(") == 0 || strings.Index(cItem.Type, "char(") == 0 ||
				strings.Index(cItem.Type, "text") == 0 ||
				strings.Index(cItem.Type, "mediumtext") == 0 ||
				strings.Index(cItem.Type, "longtext") == 0 {
				if cItem.Null == "YES" {
					fieldType = "sql.NullString"
					bImportSql = true
				} else {
					fieldType = "string"
				}
			} else if cItem.Type == "float" || cItem.Type == "double" {
				fieldType = "float64"
			} else {
				fieldType = "string"
			}

			code += fmt.Sprintf("\t%s %s `db:\"%s\" json:\"%s\"`\n", name, fieldType, cItem.Field, cItem.Field)
		}
		code += fmt.Sprintf("}\n")
	}

	if bImportSql {
		fmt.Println(`import ("database/sql")`)
	}

	fmt.Println(code)

}
