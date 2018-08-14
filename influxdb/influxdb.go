package influxdb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"viger.click/libs/tool"
)

type InfluxDB struct {
	Host string
	Port int
	User string
	Pwd  string
}

type RspData struct {
	Results []struct {
		Statement_id int
		Series       []struct {
			Name    string
			Columns interface{}
			Values  interface{}
		}
	}
}

func (self *InfluxDB) Request(url string, method string, reqD string) string {

	req, err := http.NewRequest(method, url, strings.NewReader(reqD))
	tool.CheckErr(err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(self.User, self.Pwd)

	rsp, err := http.DefaultClient.Do(req)
	tool.CheckErr(err)

	rspD, err := ioutil.ReadAll(rsp.Body)
	tool.CheckErr(err)

	return string(rspD)
}

// Query Data
func (self *InfluxDB) Query(db string, sql string) string {

	reqURL := fmt.Sprintf("http://%v:%v/query?db=%v", self.Host, self.Port, db)

	reqV := url.Values{}
	reqV.Set("q", sql)
	reqD := reqV.Encode()

	rspD := self.Request(reqURL, "POST", reqD)
	return rspD
}

// Create Database
func (self *InfluxDB) Create(db string) string {

	reqURL := fmt.Sprintf("http://%v:%v/query", self.Host, self.Port)

	reqV := url.Values{}
	reqV.Set("q", fmt.Sprintf("CREATE DATABASE %v", db))
	reqD := reqV.Encode()

	rspD := self.Request(reqURL, "POST", reqD)
	return rspD
}

// Drop Database
func (self *InfluxDB) Delete(db string) string {

	reqURL := fmt.Sprintf("http://%v:%v/query", self.Host, self.Port)

	reqV := url.Values{}
	reqV.Set("q", fmt.Sprintf("DROP DATABASE %v", db))
	reqD := reqV.Encode()

	rspD := self.Request(reqURL, "POST", reqD)
	return rspD
}

// Write Data
func (self *InfluxDB) Write(db string, rp string, sql string) string {

	reqURL := fmt.Sprintf("http://%v:%v/write?db=%v&rp=%v", self.Host, self.Port, db, rp)
	reqD := sql
	rspD := self.Request(reqURL, "POST", reqD)
	return rspD
}

func (self *InfluxDB) Combination(table string, cols map[string][]string, data []map[string]interface{}) string {

	var sqls []string
	var tags []string
	var vals []string
	var sql string

	sqls = []string{}
	for _, n := range data {
		// Tags Value
		tags = []string{table}
		for _, m := range cols["tag"] {
			tag := fmt.Sprintf("%v=%v", m, n[m])
			tags = append(tags, tag)
		}
		// Field Value
		vals = []string{}
		for _, m := range cols["val"] {
			switch n[m].(type) {
			case int, int64, float64:
				val := fmt.Sprintf("%v=%v", m, n[m])
				vals = append(vals, val)
			default:
				val := fmt.Sprintf("%v=\"%v\"", m, n[m])
				vals = append(vals, val)
			}
		}
		// Time
		if val, ok := n["time"]; ok {
			stamp := fmt.Sprintf("%v", val.(int))
			stamp += strings.Repeat("0", int(19-len(stamp)))
			sql = fmt.Sprintf("%v %v %v", strings.Join(tags, ","), strings.Join(vals, ","), stamp)
		} else {
			sql = fmt.Sprintf("%v %v", strings.Join(tags, ","), strings.Join(vals, ","))
		}
		sqls = append(sqls, sql)
	}

	return strings.Join(sqls, "\n")
}
