package hosts

import (
	"container/list"
	"math/rand"
	"time"
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

func (l *List) Ping(count int, interval int, resultChan chan *Result) {
	offset := float64(interval) / float64(l.Length())
	duration := time.Duration(offset * 1000000000)
	for e := l.Hosts.Front(); e != nil; e = e.Next() {
		time.Sleep(duration)
		go func(host *Host, count int) {
			dur := time.Duration(interval) * time.Second
			ticker := time.NewTicker(dur)
			for range ticker.C {
				drift := time.Duration(rand.Intn(500))
				time.Sleep(drift * time.Millisecond)
				resultChan <- &Result{
					Stats: host.Ping(count),
					Tags:  host.GetTags(),
					When:  time.Now().UTC(),
				}
			}
		}(e.Value.(*Host), count)
	}
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
