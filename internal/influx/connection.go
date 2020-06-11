package influx

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type Connection struct {
	client   client.Client
	database string
}

func New(addr string, db string, username string, password string) (*Connection, error) {
	conn := &Connection{}
	err := conn.Open(addr, db, username, password)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (conn *Connection) Open(addr string, db string, username string, password string) error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}
	conn.client = c
	conn.database = db
	return nil
}

func (conn *Connection) Close() {
	conn.client.Close()
}

func (conn *Connection) Ping(timeout time.Duration) (time.Duration, string, error) {
	return conn.client.Ping(timeout)
}
