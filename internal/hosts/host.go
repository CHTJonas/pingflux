package hosts

import (
	"net"
	"sync"

	ping "github.com/stenya/go-ping"
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

func (h *Host) GetTags() map[string]string {
	tags := make(map[string]string)
	for key, value := range h.Tags {
		tags[key] = value
	}
	tags["name"] = h.GetName()
	return tags
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
	var wg sync.WaitGroup
	endpoints := h.GetEndpoints()
	statistics := make([]*ping.Statistics, len(endpoints))
	for i, e := range endpoints {
		wg.Add(1)
		go func(i int, e *Endpoint) {
			defer wg.Done()
			statistics[i] = e.Ping(count)
		}(i, e)
	}
	wg.Wait()
	return statistics
}
