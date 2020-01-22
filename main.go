package main

import (
	"container/list"

	pingu "github.com/sparrc/go-ping"
)

func main() {
	startInit()
	for e := hosts.Front(); e != nil; e = e.Next() {
		p := e.Value.(*pair)
		ping(p.ip, p.friendlyName)
	}
}

func startInit() {
	hosts = list.New()
	StaticIPAddresses := []string{"192.168.88.4", "146.97.41.38", "146.97.41.46"}
	StaticHostnames := []string{"gw.eng.cam.ac.uk"}
	for _, ip := range StaticIPAddresses {
		addIP(ip)
	}
	for _, hostname := range StaticHostnames {
		addHostname(hostname)
	}
}

func ping(host string, friendlyName string) {
	pinger, err := pingu.NewPinger(host)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.Run() // blocks until finished
	stats := pinger.Statistics()
	storeInInflux(friendlyName, stats.MinRtt, stats.AvgRtt, stats.MaxRtt)
}
