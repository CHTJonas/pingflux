package hosts

import (
	"container/list"
	"math/rand"
	"time"

	ping "github.com/stenya/go-ping"
)

type List struct {
	Hosts *list.List
}

func NewList() *List {
	return &List{
		Hosts: list.New(),
	}
}

func (l *List) Reset() {
	l.Hosts = list.New()
}

func (l *List) Length() int {
	return l.Hosts.Len()
}

func (l *List) AddIP(ip string, tags map[string]string) {
	host := &Host{
		IP:   ip,
		Tags: tags,
	}
	l.Hosts.PushBack(host)
}

func (l *List) AddHostname(hostname string, tags map[string]string) {
	host := &Host{
		Hostname: hostname,
		Tags:     tags,
	}
	l.Hosts.PushBack(host)
}

func (l *List) Ping(count int, callback func([]*ping.Statistics, *Host)) {
	for e := l.Hosts.Front(); e != nil; e = e.Next() {
		r := rand.Intn(1000)
		d := time.Duration(r)
		time.Sleep(d * time.Millisecond)
		go func(host *Host, count int) {
			ticker := time.NewTicker(10 * time.Second)
			for {
				select {
				case <-ticker.C:
					statistics := host.Ping(count)
					callback(statistics, host)
				}
			}
		}(e.Value.(*Host), count)
	}
}
