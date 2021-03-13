package hosts

import (
	"time"

	"github.com/go-ping/ping"
)

type Result struct {
	Stats *ping.Statistics
	Tags  map[string]string
	When  time.Time
}
