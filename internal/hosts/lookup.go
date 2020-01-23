package hosts

import "net"

func AddIP(ip string) {
	hostnames, err := net.LookupAddr(ip)
	if err != nil {
		h := &Host{
			IP:           ip,
			FriendlyName: ip,
		}
		Endpoints.PushBack(h)
	} else {
		for _, hostname := range hostnames {
			h := &Host{
				IP:           ip,
				FriendlyName: hostname,
			}
			Endpoints.PushBack(h)
		}
	}
}

func AddHostname(hostname string) {
	ipAddresses, err := net.LookupHost(hostname)
	if err != nil {
		panic(err)
	} else {
		for _, ipAddress := range ipAddresses {
			h := &Host{
				IP:           ipAddress,
				FriendlyName: hostname,
			}
			Endpoints.PushBack(h)
		}
	}
}
