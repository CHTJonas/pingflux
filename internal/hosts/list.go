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

func (l *List) Ping(count int, interval int, callback func([]*ping.Statistics, *Host)) {
	offset := interval / l.Length()
	for e := l.Hosts.Front(); e != nil; e = e.Next() {
		time.Sleep(time.Duration(offset) * time.Second)
		go func(host *Host, count int) {
			dur := time.Duration(interval) * time.Second
			ticker := time.NewTicker(dur)
			for range ticker.C {
				drift := time.Duration(rand.Intn(500))
				time.Sleep(drift * time.Millisecond)
				statistics := host.Ping(count)
				callback(statistics, host)
			}
		}(e.Value.(*Host), count)
	}
}
