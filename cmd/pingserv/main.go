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
		h := e.Value.(*hosts.Host)
		stats := h.Ping()
		conn.Store(stats, h.FriendlyName)
	}
}

func initConnection() {
	addr := "http://localhost:8086"
	db := "pingserv"
	conn = &influx.Connection{}
	conn.Open(addr, db)
}

func initHosts() {
	StaticIPAddresses := []string{"192.168.88.4", "146.97.41.38", "146.97.41.46"}
	StaticHostnames := []string{"gw.eng.cam.ac.uk"}
	hosts.ResetEndpoints()
	for _, ip := range StaticIPAddresses {
		hosts.AddIP(ip)
	}
	for _, hostname := range StaticHostnames {
		hosts.AddHostname(hostname)
	}
}
