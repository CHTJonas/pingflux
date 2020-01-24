package influx

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	pingu "github.com/sparrc/go-ping"
)

func (conn *Connection) Store(stats *pingu.Statistics, friendlyName string) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  conn.database,
		Precision: "s",
	})
	if err != nil {
		panic(err)
	}
	bp.AddPoint(generateRTTPoint(stats, friendlyName))
	bp.AddPoint(generatePacketsPoint(stats, friendlyName))
	conn.client.Write(bp)
}

func getTags(friendlyName string) map[string]string {
	return map[string]string{
		"host": friendlyName,
	}
}

func generateRTTPoint(stats *pingu.Statistics, friendlyName string) *client.Point {
	tags := getTags(friendlyName)
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

func generatePacketsPoint(stats *pingu.Statistics, friendlyName string) *client.Point {
	tags := getTags(friendlyName)
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
