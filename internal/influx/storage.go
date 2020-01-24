package influx

import (
	"time"

	"github.com/chtjonas/pingserv/internal/hosts"
	client "github.com/influxdata/influxdb1-client/v2"
	pingu "github.com/sparrc/go-ping"
)

func (conn *Connection) Store(stats *pingu.Statistics, host *hosts.Host) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  conn.database,
		Precision: "s",
	})
	if err != nil {
		panic(err)
	}
	bp.AddPoint(generateRTTPoint(stats, host))
	bp.AddPoint(generatePacketsPoint(stats, host))
	conn.client.Write(bp)
}

func getTags(host *hosts.Host) map[string]string {
	tags := make(map[string]string)
	for key, value := range host.Tags {
		tags[key] = value
	}
	tags["host"] = host.FriendlyName
	return tags
}

func generateRTTPoint(stats *pingu.Statistics, host *hosts.Host) *client.Point {
	tags := getTags(host)
	fields := map[string]interface{}{
		"min": stats.MinRtt.Seconds(),
		"avg": stats.AvgRtt.Seconds(),
		"max": stats.MaxRtt.Seconds(),
	}
	pt, err := client.NewPoint("rtt", tags, fields, time.Now().UTC())
	if err != nil {
		panic(err)
	}
	return pt
}

func generatePacketsPoint(stats *pingu.Statistics, host *hosts.Host) *client.Point {
	tags := getTags(host)
	fields := map[string]interface{}{
		"sent": stats.PacketsSent,
		"recv": stats.PacketsRecv,
		"loss": stats.PacketLoss,
	}
	pt, err := client.NewPoint("packets", tags, fields, time.Now().UTC())
	if err != nil {
		panic(err)
	}
	return pt
}
