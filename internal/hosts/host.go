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
	names, err := h.ReverseIP()
	if err != nil {
		return h.IP
	}
	return names[0]
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
	addrs, err := h.ResolveHostname()
	if err != nil {
		return []*Endpoint{}
	}
	endpoints := make([]*Endpoint, len(addrs))
	for i, addr := range addrs {
		endpoints[i] = NewEndpoint(addr)
	}
	return endpoints
}

func (h *Host) ResolveHostname() ([]string, error) {
	addrs, err := net.LookupHost(h.Hostname)
	if err != nil {
		return nil, err
	}
	return addrs, nil
}

func (h *Host) ReverseIP() ([]string, error) {
	names, err := net.LookupAddr(h.IP)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (h *Host) Ping(count int) []*ping.Statistics {
	var wg sync.WaitGroup
	endpoints := h.GetEndpoints()
	statistics := make([]*ping.Statistics, len(endpoints))
	for i, e := range endpoints {
		wg.Add(1)
		go func(i int, e *Endpoint) {
			defer wg.Done()
			stats, err := e.Ping(count)
			if err != nil {
				fmt.Println("Failed to ping", e.IP, err)
			} else {
				statistics[i] = stats
			}
		}(i, e)
	}
	wg.Wait()
	return statistics
}
