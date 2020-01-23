package influx

import (
	client "github.com/influxdata/influxdb1-client/v2"
)

type Connection struct {
	client   client.Client
	database string
}

func (conn *Connection) Open(addr string, db string) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		panic(err)
	}
	conn.client = c
	conn.database = db
}

func (conn *Connection) Close() {
	conn.client.Close()
}
