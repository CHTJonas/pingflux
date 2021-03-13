package hosts

import (
	"fmt"
	"net"
	"sync"
	"time"
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

func (h *Host) Ping(count int) *[]*Result {
	var wg sync.WaitGroup
	endpoints := h.GetEndpoints()
	results := make([]*Result, len(endpoints))
	for i, e := range endpoints {
		wg.Add(1)
		go func(i int, e *Endpoint) {
			defer wg.Done()
			tags := h.GetTags()
			for k, v := range e.GetTags() {
				tags[k] = v
			}
			when := time.Now().UTC()
			stats, err := e.Ping(count)
			if err != nil {
				fmt.Println("Failed to ping", e.IP, err)
			} else {
				results[i] = &Result{
					Stats: stats,
					Tags:  tags,
					When:  when,
				}
			}
		}(i, e)
	}
	wg.Wait()
	return &results
}
