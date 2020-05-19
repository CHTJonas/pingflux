package hosts

import ping "github.com/stenya/go-ping"

type Result struct {
	Stats []*ping.Statistics
	Host  *Host
}
