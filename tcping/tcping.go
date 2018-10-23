// File: tcping.go
// Author: walker
// Changelogs:
//   2018.08.16: init

//func main() {
//	var tcpinger TCPinger = TCPinger{
//		Target:   "10.180.5.107",
//		Port:     3000,
//		Times:    5,
//		Timeout:  2,
//		Interval: 1,
//		Statis:   map[string]float64{},
//	}
//	tcpinger.Run()
//	fmt.Println(tcpinger.Statis)
//}

package tcping

import (
	"fmt"
	"log"
	"math"
	"net"
	"time"
)

type TCPinger struct {
	Target   string
	IP       string
	Port     int
	Times    int
	Timeout  int
	Interval int
	Statis   map[string]float64
	RTTs     []float64
	ch       chan int
}

func (self *TCPinger) Run() {
	self.ch = make(chan int)
	ns, err := net.LookupHost(self.Target)
	if err != nil {
		log.Printf(err)
		return
	}
	self.IP = ns[0]
	for i := 0; i < self.Times; i++ {
		go self.TCPing()
		time.Sleep(time.Duration(self.Interval) * time.Second)
	}
	for i := 0; i < self.Times; i++ {
		<-self.ch
	}
	self.Statis = self.Statistics()
}

func (self *TCPinger) TCPing() {
	var addr string
	var sTime, eTime int64
	var conn net.Conn
	var err error
	var floatRTT float64

	sTime = time.Now().UnixNano()
	addr = fmt.Sprintf("%v:%v", self.IP, self.Port)
	conn, err = net.DialTimeout("tcp", addr, time.Duration(self.Timeout)*time.Second)
	if err == nil {
		conn.Close()
		eTime = time.Now().UnixNano()
		floatRTT = float64(eTime-sTime) / float64(time.Second)
		self.RTTs = append(self.RTTs, floatRTT)
	}
	self.ch <- 0
}

func (self *TCPinger) Statistics() map[string]float64 {

	var loss float64
	var ansNum int
	var rttMax float64
	var rttMin float64 = 999
	var rttAvg float64
	var rttStd float64

	var sum float64

	ansNum = len(self.RTTs)
	statis := map[string]float64{}

	loss = 1 - (float64(ansNum) / float64(self.Times))

	if ansNum > 0 {
		for _, v := range self.RTTs {
			rttMax = math.Max(rttMax, v)
			rttMin = math.Min(rttMin, v)
			sum += v
		}
		rttAvg = sum / float64(ansNum)

		sum = 0
		for _, v := range self.RTTs {
			sum += math.Pow((v - rttAvg), 2)
		}

		rttStd = math.Sqrt(sum / float64(ansNum))

		statis["loss"] = loss
		statis["max"] = rttMax
		statis["min"] = rttMin
		statis["avg"] = rttAvg
		statis["std"] = rttStd
		statis["num"] = float64(ansNum)
	} else {
		statis["loss"] = loss
		statis["max"] = 0.0
		statis["min"] = 0.0
		statis["avg"] = 0.0
		statis["std"] = 0.0
		statis["num"] = float64(ansNum)
	}
	return statis
}
