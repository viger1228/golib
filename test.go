// File: tcping.go
// Author: walker
// Changelogs:
//   2018.08.16: init

package main

import (
	"fmt"

	"github.com/viger1228/golib/tcping"
	"github.com/viger1228/golib/tool"
)

func main() {
	tcpinger := tcping.TCPinger{
		Target:   "10.180.5.107",
		Port:     3000,
		Times:    5,
		Timeout:  2,
		Interval: 1,
		Statis:   map[string]float64{},
	}
	tcpinger.Run()
	fmt.Println(string(tool.DumpsJson(tcpinger.Statis, 2)))
}
