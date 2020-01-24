package main

import (
	"github.com/chtjonas/pingserv/internal/hosts"
	"github.com/chtjonas/pingserv/internal/influx"
)

var conn *influx.Connection

func main() {
	initHosts()
	initConnection()
	defer conn.Close()
	for e := hosts.Endpoints.Front(); e != nil; e = e.Next() {
		host := e.Value.(*hosts.Host)
		stats := host.Ping()
		conn.Store(stats, host)
	}
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
