package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chtjonas/pingserv/internal/hosts"
	"github.com/chtjonas/pingserv/internal/influx"
)

var conn *influx.Connection

func main() {
	count := 3
	initHosts()
	initConnection()
	defer conn.Close()

	for e := hosts.Endpoints.Front(); e != nil; e = e.Next() {
		r := rand.Intn(1000)
		d := time.Duration(r)
		time.Sleep(d * time.Millisecond)
		go setupPinger(e.Value.(*hosts.Host), count)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	for range stop {
		os.Exit(0)
	}
}

func setupPinger(host *hosts.Host, count int) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			go ping(host, count)
		}
	}
}

func ping(host *hosts.Host, count int) {
	stats := host.Ping(count)
	conn.Store(stats, host)
}

func initConnection() {
	addr := "http://localhost:8086"
	db := "pingserv"
	conn = &influx.Connection{}
	conn.Open(addr, db)
}

func initHosts() {
	hosts.ResetEndpoints()

	IPAddresseMappings := []map[string]map[string]string{{
		"192.168.88.4": {
			"network": "JFDN",
			"server":  "storage",
		},
		"146.97.41.38": {
			"network": "CUDN",
			"router":  "border",
		},
		"146.97.41.46": {
			"network": "CUDN",
			"router":  "border",
		},
	}}
	for _, mapping := range IPAddresseMappings {
		for ip, tags := range mapping {
			hosts.AddIP(ip, tags)
		}
	}

	hostnameMappings := []map[string]map[string]string{
		{
			"gw.eng.cam.ac.uk": {
				"network": "CUDN",
				"router":  "institution",
			},
		},
	}
	for _, mapping := range hostnameMappings {
		for hostname, tags := range mapping {
			hosts.AddHostname(hostname, tags)
		}
	}
}
