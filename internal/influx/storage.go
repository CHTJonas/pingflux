package influx

import (
	"time"

	"github.com/chtjonas/pingflux/internal/hosts"
	"github.com/go-ping/ping"
	client "github.com/influxdata/influxdb1-client/v2"
)

func (conn *Connection) Store(resultsArrPtr *[]*hosts.Result) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	batch, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  conn.database,
		Precision: "s",
	})
	resultsArr := *resultsArrPtr
	for i := 0; i < len(resultsArr); i++ {
		result := resultsArr[i]
		if result == nil {
			continue
		}
		stats := result.Stats
		if stats != nil {
			batch.AddPoint(generateRTTPoint(stats, result.Tags, result.When))
			batch.AddPoint(generatePacketsPoint(stats, result.Tags, result.When))
		}
	}
	return conn.client.Write(batch)
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
		"dup":  stats.PacketsRecvDuplicates,
	}
	point, err := client.NewPoint("packets", tags, fields, when)
	if err != nil {
		panic(err)
	}
	return point
}
