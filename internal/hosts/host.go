package hosts

import (
	"container/list"

	pingu "github.com/sparrc/go-ping"
)

var Endpoints *list.List

type Host struct {
	IP           string
	FriendlyName string
}

func (h *Host) Ping() *pingu.Statistics {
	pinger, err := pingu.NewPinger(h.IP)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.Run() // blocks until finished
	return pinger.Statistics()
}
