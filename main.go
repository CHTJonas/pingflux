package main

import (
	"fmt"
	"net"

	pingu "github.com/sparrc/go-ping"
)

func main() {
	iprecords, _ := net.LookupIP("gw.eng.cam.ac.uk")
	for _, ip := range iprecords {
		ping(ip.String())
	}
}

func ping(host string) {
	pinger, err := pingu.NewPinger(host)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.Run() // blocks until finished
	stats := pinger.Statistics()
	fmt.Println(stats)
}
