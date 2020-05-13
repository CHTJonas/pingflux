package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chtjonas/pingflux/internal/hosts"
	"github.com/chtjonas/pingflux/internal/influx"
)

var list *hosts.List
var conn *influx.Connection

func main() {
	count := 3
	initHosts()
	initConnection()
	defer conn.Close()

	_, ver, err := conn.Ping(0)
	if err != nil {
		fmt.Println("Error pinging InfluxDB server: ", err.Error())
	} else {
		fmt.Println("Got reply from InfluxDB server: ", ver)
	}

	for e := list.Hosts.Front(); e != nil; e = e.Next() {
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
			statistics := host.Ping(count)
			conn.Store(statistics, host)
		}
	}
}

func initConnection() {
	addr := "http://localhost:8086"
	db := "pingflux"
	conn = influx.New(addr, db)
}

func initHosts() {
	list = hosts.NewList()

	IPAddresseMappings := []map[string]map[string]string{{
		"192.168.86.5": {
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
			list.AddIP(ip, tags)
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
			list.AddHostname(hostname, tags)
		}
	}
}
