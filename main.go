package main

import (
	"net"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	pingu "github.com/sparrc/go-ping"
)

func main() {
	StaticIPAddresses := []string{} //"192.168.88.1", "192.168.88.4"}
	StaticHostnames := []string{"gw.eng.cam.ac.uk"}
	for _, ip := range StaticIPAddresses {
		pingByIP(ip)
	}
	for _, hostname := range StaticHostnames {
		pingByHostname(hostname)
	}
}

func pingByIP(ip string) {
	// names, err := net.LookupAddr(ip)
	// if err != nil {
	// 	fmt.Println("noname")
	// } else {
	// 	fmt.Println(names[0])
	// }
	ping(ip)
}

func pingByHostname(hostname string) {
	// fmt.Println(hostname)
	records, err := net.LookupIP(hostname)
	if err != nil {
		panic(err)
	}
	for _, ip := range records {
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
	storeInInflux(stats.MinRtt, stats.AvgRtt, stats.MaxRtt)
}

func storeInInflux(minRtt time.Duration, avgRtt time.Duration, maxRtt time.Duration) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "mydb",
		Precision: "s",
	})

	fields := map[string]interface{}{
		"min": minRtt.Seconds(),
		"avg": avgRtt.Seconds(),
		"max": maxRtt.Seconds(),
	}

	pt, err := client.NewPoint("ping_rtt", nil, fields, time.Now().UTC())
	if err != nil {
		panic(err)
	}

	bp.AddPoint(pt)
	c.Write(bp)
}
