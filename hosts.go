package main

import (
	"container/list"
	"net"
)

var hosts *list.List

type pair struct {
	ip           string
	friendlyName string
}

func addIP(ip string) {
	hostnames, err := net.LookupAddr(ip)
	if err != nil {
		p := &pair{
			ip:           ip,
			friendlyName: "noname",
		}
		hosts.PushBack(p)
	} else {
		for _, hostname := range hostnames {
			p := &pair{
				ip:           ip,
				friendlyName: hostname,
			}
			hosts.PushBack(p)
		}
	}
}

func addHostname(hostname string) {
	ipAddresses, err := net.LookupHost(hostname)
	if err != nil {
		panic(err)
	} else {
		for _, ipAddress := range ipAddresses {
			p := &pair{
				ip:           ipAddress,
				friendlyName: hostname,
			}
			hosts.PushBack(p)
		}
	}
}

// func pingByHostname(hostname string) {
// 	// fmt.Println(hostname)
// 	records, err := net.LookupIP(hostname)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, ip := range records {
// 		ping(ip.String())
// 	}
// }
