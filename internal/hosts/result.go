package hosts

import (
	"time"

	ping "github.com/stenya/go-ping"
)

type Result struct {
	Stats []*ping.Statistics
	Tags  map[string]string
	When  time.Time
}
