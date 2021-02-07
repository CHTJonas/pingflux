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

func (e *Endpoint) Ping(count int) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(e.IP)
	if err != nil {
		return nil, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = count
	if e.isIPv4() {
		pinger.Size = 40
	} else {
		pinger.Size = 1232
	}
	pinger.Interval = time.Second
	pinger.Timeout = time.Second * 10
	pinger.Run() // blocks until finished
	return pinger.Statistics(), nil
}

func (e *Endpoint) isIPv4() bool {
	ip := net.ParseIP(e.IP)
	return len(ip.To4()) == net.IPv4len
}
