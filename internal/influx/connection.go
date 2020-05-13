package influx

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type Connection struct {
	client   client.Client
	database string
}

func New(addr string, db string) *Connection {
	conn := &Connection{}
	conn.Open(addr, db)
	return conn
}

func (conn *Connection) Open(addr string, db string) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: addr,
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

func (conn *Connection) Ping(timeout time.Duration) (time.Duration, string, error) {
	return conn.client.Ping(timeout)
}
