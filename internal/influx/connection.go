package influx

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type Connection struct {
	client   client.Client
	database string
}

func New(addr, db, username, password, userAgent string) (*Connection, error) {
	conn := &Connection{}
	err := conn.Open(addr, db, username, password, userAgent)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (conn *Connection) Open(addr, db, username, password, userAgent string) error {
	conf := client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	}
	if userAgent != "" {
		conf.UserAgent = userAgent
	}
	c, err := client.NewHTTPClient(conf)
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
