package hosts

import (
	"fmt"
	"time"

	ping "github.com/stenya/go-ping"
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
	pinger.SetPrivileged(true)
	pinger.Count = count
	pinger.Interval = time.Second
	pinger.Timeout = time.Second * 10
	fmt.Printf("Pinging %s\n", e.IP)
	pinger.Run() // blocks until finished
	return pinger.Statistics()
}
