package influx

import (
	"time"

	"github.com/chtjonas/pingflux/internal/hosts"
	client "github.com/influxdata/influxdb1-client/v2"
	ping "github.com/stenya/go-ping"
)

func (conn *Connection) Store(statistics []*ping.Statistics, host *hosts.Host) {
	batch, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  conn.database,
		Precision: "s",
	})
	if err != nil {
		panic(err)
	}
	for _, stats := range statistics {
		batch.AddPoint(generateRTTPoint(stats, host))
		batch.AddPoint(generatePacketsPoint(stats, host))
	}
	conn.client.Write(batch)
}

func generateRTTPoint(stats *ping.Statistics, host *hosts.Host) *client.Point {
	fields := map[string]interface{}{
		"min": stats.MinRtt.Seconds(),
		"avg": stats.AvgRtt.Seconds(),
		"max": stats.MaxRtt.Seconds(),
	}
	point, err := client.NewPoint("rtt", host.GetTags(), fields, time.Now().UTC())
	if err != nil {
		panic(err)
	}
	return point
}

func generatePacketsPoint(stats *ping.Statistics, host *hosts.Host) *client.Point {
	fields := map[string]interface{}{
		"sent": stats.PacketsSent,
		"recv": stats.PacketsRecv,
		"loss": stats.PacketLoss,
	}
	point, err := client.NewPoint("packets", host.GetTags(), fields, time.Now().UTC())
	if err != nil {
		panic(err)
	}
	return point
}
