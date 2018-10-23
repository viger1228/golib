// File: trace.go
// Author: walker
// Changelogs:
//   2018.10.23: init

//func main() {
//	var tracer Tracer = Tracer{
//		Target:   "www.baidu.com",
//		Times:    5,
//		Timeout:  2,
//		Interval: 1,
//	}
//	tracer.Run()
//	for _, host := range tracer.Statis {
//		fmt.Printf("%+v\n", host)
//	}
//}

package trace

import (
	"fmt"
	"log"
	"math"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

type Tracer struct {
	Target   string
	IP       string
	Times    int
	Timeout  int
	Interval int
	Statis   []*Host
	ch       chan int
}

type Host struct {
	Hop  int
	IP   string
	Loss float64
	Num  float64
	Max  float64
	Min  float64
	Avg  float64
	Std  float64
	RTTs []float64
}

func (self *Tracer) Run() {
	self.ch = make(chan int)
	ns, err := net.LookupHost(self.Target)
	if err != nil {
		log.Println(err)
		return
	}
	self.IP = ns[0]
	self.Trace()
}

func (self *Tracer) Trace() {
	args := []string{
		"-n",
		"--raw",
		fmt.Sprintf("-c %v", self.Times),
		fmt.Sprintf("--interval %v", self.Interval),
		fmt.Sprintf("--timeout %v", self.Timeout),
		self.IP,
	}
	raw, err := exec.Command("mtr", args...).Output()
	if err != nil {
		log.Println(err)
		return
	}
	line := strings.Split(string(raw), "\n")
	index := map[int]int{}
	last := 0
	for _, n := range line {
		list := strings.Fields(n)
		if len(list) != 3 {
			continue
		}
		type_ := list[0]
		hop, _ := strconv.Atoi(list[1])
		data := list[2]
		if hop > last {
			last = hop
		}
		if type_ == "h" {
			host := &Host{
				Hop: hop,
				IP:  data,
			}
			index[hop] = len(self.Statis)
			self.Statis = append(self.Statis, host)
		} else if type_ == "p" {
			l := self.Statis[index[hop]].RTTs
			t, _ := strconv.Atoi(data)
			self.Statis[index[hop]].RTTs = append(l, float64(t)/1000000)
		}
	}
	n := index[last]
	m := index[last-1]
	if self.Statis[n].IP == self.Statis[m].IP {
		self.Statis = append(self.Statis[:n], self.Statis[n+1:]...)
	}

	var loss float64
	var num int
	var rttMax float64
	var rttMin float64
	var rttAvg float64
	var rttStd float64
	var sum float64

	for _, host := range self.Statis {

		rttMax = 0
		rttMin = 9999
		sum = 0

		num = len(host.RTTs)
		loss = 1 - (float64(num) / float64(self.Times))

		if num > 0 {
			for _, v := range host.RTTs {
				rttMax = math.Max(rttMax, v)
				rttMin = math.Min(rttMin, v)
				sum += v
			}
			rttAvg = sum / float64(num)

			sum = 0
			for _, v := range host.RTTs {
				sum += math.Pow((v - rttAvg), 2)
			}

			rttStd = math.Sqrt(sum / float64(num))

			host.Num = float64(num)
			host.Loss = loss
			host.Max = rttMax
			host.Min = rttMin
			host.Avg = rttAvg
			host.Std = rttStd
		} else {
			host.Num = float64(num)
			host.Loss = loss
			host.Max = 0.0
			host.Min = 0.0
			host.Avg = 0.0
			host.Std = 0.0
		}
	}
}
