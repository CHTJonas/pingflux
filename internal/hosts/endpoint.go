package hosts

import (
	"net"
	"time"

	"github.com/go-ping/ping"
)

type Endpoint struct {
	IP string
}

func NewEndpoint(ip string) *Endpoint {
	return &Endpoint{
		IP: ip,
	}
}

func (e *Endpoint) GetTags() map[string]string {
	tags := make(map[string]string)
	if e.IsIPv4() {
		tags["protocol"] = "ICMPv4"
	} else {
		tags["protocol"] = "ICMPv6"
	}
	return tags
}

func (e *Endpoint) Ping(count int) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(e.IP)
	if err != nil {
		return nil, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = count
	pinger.Size = 56
	pinger.Interval = time.Second
	pinger.Timeout = time.Second * 10
	err = pinger.Run() // blocks until finished
	if err != nil {
		return nil, err
	}
	return pinger.Statistics(), nil
}

func (e *Endpoint) IsIPv4() bool {
	ip := net.ParseIP(e.IP)
	return len(ip.To4()) == net.IPv4len
}
