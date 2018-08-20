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
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/viger1228/golib/tool"
)

type MySQL struct {
	Host     string
	Port     int
	User     string
	Password string
	DataBase string
}

func (self *MySQL) Cmd(SQL string) ([]string, [][]interface{}) {

	key := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
		self.User, self.Password, self.Host, self.Port, self.DataBase)
	db, err := sql.Open("mysql", key)
	defer db.Close()
	tool.CheckErr(err)

	rsp, err := db.Query(SQL)
	tool.CheckErr(err)

	columns, err := rsp.Columns()
	tool.CheckErr(err)

	_type := []string{}
	colTypes, err := rsp.ColumnTypes()
	tool.CheckErr(err)
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
		tool.CheckErr(err)
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
		self.User, self.Password, self.Host, self.Port, self.DataBase)
	db, err := sql.Open("mysql", key)
	defer db.Close()
	tool.CheckErr(err)
	stmt, err := db.Prepare(SQL)
	tool.CheckErr(err)

	_, err = stmt.Exec()
	tool.CheckErr(err)
}
