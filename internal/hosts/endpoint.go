package hosts

import (
	ping "github.com/sparrc/go-ping"
)

type Endpoint struct {
	IP string
}

func NewEndpoint(ip string) *Endpoint {
	return &Endpoint{
		IP: ip,
	}
}

func (e *Endpoint) Ping(count int) *ping.Statistics {
	pinger, err := ping.NewPinger(e.IP)
	if err != nil {
		panic(err)
	}
	pinger.Count = count
	pinger.Run() // blocks until finished
	return pinger.Statistics()
}
