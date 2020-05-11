package hosts

import (
	"container/list"
)

var Hosts *list.List

func Reset() {
	Hosts = list.New()
}

func AddIP(ip string, tags map[string]string) {
	host := &Host{
		IP:   ip,
		Tags: tags,
	}
	Hosts.PushBack(host)
}

func AddHostname(hostname string, tags map[string]string) {
	host := &Host{
		Hostname: hostname,
		Tags:     tags,
	}
	Hosts.PushBack(host)
}
