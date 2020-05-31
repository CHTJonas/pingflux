package influx

import (
	"container/list"
	"time"

	"github.com/chtjonas/pingflux/internal/hosts"
	client "github.com/influxdata/influxdb1-client/v2"
	ping "github.com/stenya/go-ping"
)

func (conn *Connection) Store(resultList *list.List) {
	batch, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  conn.database,
		Precision: "s",
	})
	if err != nil {
		panic(err)
	}
	for e := resultList.Front(); e != nil; e = e.Next() {
		result := e.Value.(*hosts.Result)
		for _, stats := range result.Stats {
			if stats != nil {
				batch.AddPoint(generateRTTPoint(stats, result.Tags, result.When))
				batch.AddPoint(generatePacketsPoint(stats, result.Tags, result.When))
			}
		}
	}
	conn.client.Write(batch)
}

func generateRTTPoint(stats *ping.Statistics, tags map[string]string, when time.Time) *client.Point {
	fields := map[string]interface{}{
		"min": stats.MinRtt.Seconds(),
		"avg": stats.AvgRtt.Seconds(),
		"max": stats.MaxRtt.Seconds(),
	}
	point, err := client.NewPoint("rtt", tags, fields, when)
	if err != nil {
		panic(err)
	}
	return point
}

func generatePacketsPoint(stats *ping.Statistics, tags map[string]string, when time.Time) *client.Point {
	fields := map[string]interface{}{
		"sent": stats.PacketsSent,
		"recv": stats.PacketsRecv,
		"loss": stats.PacketLoss,
	}
	point, err := client.NewPoint("packets", tags, fields, when)
	if err != nil {
		panic(err)
	}
	return point
}
