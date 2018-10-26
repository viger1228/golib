// File: test.go
// Author: walker
// Mail: walkerIVI@gmail.com
// Changelogs:
//   2018.10.26: init

package main

import (
	"fmt"
	"os"

	"github.com/viger1228/golib/mysql"
	"github.com/viger1228/golib/tool"
)

func main() {

	api := mysql.MySQL{
		Host:     "10.180.5.107",
		Port:     3306,
		User:     "noc",
		Password: "iv66.net",
		DataBase: "mon",
	}

	hostname, _ := os.Hostname()
	sql := fmt.Sprintf("SELECT `hostname`,`target`,`port` FROM "+
		"t_telegraf_tcping WHERE `enable`=1 AND `hostname`='%v'", hostname)
	fmt.Println(sql)
	target := api.Query(sql)
	fmt.Println(string(tool.DumpsJson(target, 2)))

}
