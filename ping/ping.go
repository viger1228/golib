package ping

import (
	"fmt"
	"math"
	"net"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

var (
	rtts []float64
	ch   chan int
)

type Pinger struct {
	Target   string
	Times    int
	Timeout  int
	Interval int
	Statis   map[string]float64
}

func main() {

	var pinger Pinger = Pinger{
		Target:   "8.8.8.8",
		Times:    30,
		Timeout:  2,
		Interval: 1,
		Statis:   map[string]float64{},
	}

	pinger.Run()
	fmt.Println(pinger.Statis)
}

func (self *Pinger) Run() {

	ch = make(chan int)

	for i := 0; i < self.Times; i++ {
		go Ping(self.Target, ch, self.Timeout)
		time.Sleep(time.Duration(self.Interval) * time.Second)
	}

	for i := 0; i < self.Times; i++ {
		<-ch
	}

	self.Statis = Statistics(self.Times, rtts)
}

func Ping(ip string, ch chan int, timeout int) {

	ping := fastping.NewPinger()
	raddr, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ping.AddIPAddr(raddr)

	ping.MaxRTT = time.Duration(timeout) * time.Second

	ping.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		var floatRTT float64
		floatRTT = float64(rtt) / float64(time.Second)
		rtts = append(rtts, floatRTT)
	}

	ping.OnIdle = func() {
		ch <- 0
	}

	err = ping.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Statistics(num int, rtts []float64) map[string]float64 {

	var loss float64
	var ansNum int
	var rttMax float64
	var rttMin float64 = 999
	var rttAvg float64
	var rttStd float64

	var sum float64

	ansNum = len(rtts)
	statis := map[string]float64{}

	loss = 1 - (float64(ansNum) / float64(num))

	if ansNum > 0 {
		for _, v := range rtts {
			rttMax = math.Max(rttMax, v)
			rttMin = math.Min(rttMin, v)
			sum += v
		}
		rttAvg = sum / float64(ansNum)

		sum = 0
		for _, v := range rtts {
			sum += math.Pow((v - rttAvg), 2)
		}

		rttStd = math.Sqrt(sum / float64(ansNum))

		statis["loss"] = loss
		statis["max"] = rttMax
		statis["min"] = rttMin
		statis["avg"] = rttAvg
		statis["std"] = rttStd
	} else {
		statis["loss"] = loss
		statis["max"] = 0.0
		statis["min"] = 0.0
		statis["avg"] = 0.0
		statis["std"] = 0.0
	}
	return statis
}
