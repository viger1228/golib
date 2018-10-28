// File: mysql.go
// Author: walker
// Changelogs:
//   2018.08.20: init

//func main() {
//	api := MySQL{
//		Host:     "127.0.0.1",
//		Port:     3306,
//		User:     "test",
//		Password: "P@ssw0rd",
//		DataBase: "test",
//	}
//	sql := "select * from t_test"
//	data := api.Query(sql)
//	fmt.Println(tool.DumpsJson(data, 2))
//}

package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (self *MySQL) Cmd(SQL string) ([]string, [][]interface{}) {

	key := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
		self.User, self.Password, self.Host, self.Port, self.Database)
	db, err := sql.Open("mysql", key)
	defer db.Close()
	if err != nil {
		log.Printf("%v\n", err)
		return nil, nil
	}

	rsp, err := db.Query(SQL)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, nil
	}

	columns, err := rsp.Columns()
	if err != nil {
		log.Printf("%v\n", err)
		return nil, nil
	}

	_type := []string{}
	colTypes, err := rsp.ColumnTypes()
	if err != nil {
		log.Printf("%v\n", err)
		return nil, nil
	}
	for _, n := range colTypes {
		_type = append(_type, n.DatabaseTypeName())
	}

	r := make([]interface{}, len(columns))
	p := make([]interface{}, len(columns))
	for n := range r {
		p[n] = &r[n]
	}

	var rows [][]interface{}
	for rsp.Next() {
		row := make([]interface{}, len(columns))
		err = rsp.Scan(p...)
		if err != nil {
			log.Printf("%v\n", err)
			return nil, nil
		}
		for k, v := range r {
			switch i := v.(type) {
			case nil:
				row[k] = ""
			default:
				switch _type[k] {
				case "INT":
					row[k], _ = strconv.Atoi(string(i.([]uint8)))
				default:
					row[k] = fmt.Sprintf("%v", string(i.([]uint8)))
				}
			}
		}
		rows = append(rows, row)
	}
	db.Close()
	return columns, rows
}

func (self *MySQL) Query(SQL string) []map[string]interface{} {
	data := []map[string]interface{}{}
	cols, rows := self.Cmd(SQL)
	for _, n := range rows {
		sub := map[string]interface{}{}
		for i, c := range cols {
			sub[c] = n[i]
		}
		data = append(data, sub)
	}
	return data
}

func (self *MySQL) Write(SQL string) {

	key := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
		self.User, self.Password, self.Host, self.Port, self.Database)
	db, err := sql.Open("mysql", key)
	defer db.Close()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	stmt, err := db.Prepare(SQL)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
}
