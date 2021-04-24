package hosts

import (
	"container/list"
	"math/rand"
	"time"
)

type List struct {
	Hosts *list.List
	quit  chan struct{}
}

func NewList() *List {
	return &List{
		Hosts: list.New(),
		quit:  make(chan struct{}),
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

func (l *List) Ping(count int, interval int, resultChan chan *Result) {
	if l.quit == nil {
		l.quit = make(chan struct{})
	}
	offset := float64(interval) / float64(l.Length())
	duration := time.Duration(offset * float64(time.Second))
	for e := l.Hosts.Front(); e != nil; e = e.Next() {
		go func(host *Host, count int) {
			dur := time.Duration(interval * int(time.Second))
			ticker := time.NewTicker(dur)
			for {
				select {
				case <-l.quit:
					return
				case <-ticker.C:
					drift := time.Duration(rand.Intn(500))
					time.Sleep(drift * time.Millisecond)
					results := *host.Ping(count)
					for _, result := range results {
						resultChan <- result
					}
				}
			}
		}(e.Value.(*Host), count)
		time.Sleep(duration)
	}
}

func (l *List) Stop() {
	close(l.quit)
}

func (l *List) Shuffle() {
	length := l.Hosts.Len()
	a := make([]*Host, length)
	i := 0
	for e := l.Hosts.Front(); e != nil; e = e.Next() {
		a[i] = e.Value.(*Host)
		i++
	}
	rand.Shuffle(length, func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	l.Hosts.Init()
	for _, h := range a {
		l.Hosts.PushBack(h)
	}
}
