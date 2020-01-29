package hosts

import (
	pingu "github.com/sparrc/go-ping"
)

type Host struct {
	IP           string
	FriendlyName string
	Tags         map[string]string
}

func (h *Host) Ping(count int) *pingu.Statistics {
	pinger, err := pingu.NewPinger(h.IP)
	if err != nil {
		panic(err)
	}
	pinger.Count = count
	pinger.Run() // blocks until finished
	return pinger.Statistics()
}
