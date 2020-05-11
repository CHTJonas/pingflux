package hosts

import (
	"net"

	ping "github.com/sparrc/go-ping"
)

type Host struct {
	Hostname string
	IP       string
	Name     string
	Tags     map[string]string
}

func (h *Host) GetName() string {
	if h.Name != "" {
		return h.Name
	}
	if h.Hostname != "" {
		return h.Hostname
	}
	return h.IP
}

func (h *Host) GetEndpoints() []*Endpoint {
	if h.Hostname == "" {
		endpoint := NewEndpoint(h.IP)
		return []*Endpoint{endpoint}
	}
	addrs := h.ResolveHostname()
	endpoints := make([]*Endpoint, len(addrs))
	for i, addr := range addrs {
		endpoints[i] = NewEndpoint(addr)
	}
	return endpoints
}

func (h *Host) ResolveHostname() []string {
	addrs, err := net.LookupHost(h.Hostname)
	if err != nil {
		panic(err)
	}
	return addrs
}

func (h *Host) Ping(count int) []*ping.Statistics {
	endpoints := h.GetEndpoints()
	statistics := make([]*ping.Statistics, len(endpoints))
	for i, e := range endpoints {
		statistics[i] = e.Ping(count)
	}
	return statistics
}
