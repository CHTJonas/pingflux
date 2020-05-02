package hosts

import (
	"net"

	"github.com/chtjonas/pingflux/internal/influx"
	pingu "github.com/sparrc/go-ping"
)

type Host struct {
	Hostname string
	IP       string
	Name     string
	Tags     map[string]string
}

func (h *Host) Ping(count int, conn *influx.Connection) {
	endpoints := g.GetEndpoints()
	for _, endpoint := range endpoints {
		pinger, err := pingu.NewPinger(endpoint)
		if err != nil {
			panic(err)
		}
		pinger.Count = count
		pinger.Run() // blocks until finished
		stats := pinger.Statistics()
		conn.Store(stats, h)
	}
}

func (h *Host) GetName() string {
	if h.Name {
		return h.Name
	}
	if h.Hostname {
		return h.Hostname
	}
	return h.IP
}

func (h *Host) GetEndpoints() []string {
	if h.Hostname {
		return h.ResolveHostname()
	}
	return []string{h.IP}
}

func (h *Host) ResolveHostname() []string {
	addrs, err := net.LookupHost(h.Hostname)
	if err != nil {
		panic(err)
	}
	return addrs
}

func (h *Host) ReverseIP() []string {
	names, err := net.LookupAddr(ip)
	if err != nil {
		panic(err)
	}
	return names
}
