package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/chtjonas/pingflux/internal/hosts"
	"github.com/chtjonas/pingflux/internal/influx"
	"github.com/spf13/viper"
	"github.com/stenya/go-ping"
)

var list *hosts.List
var conn *influx.Connection

func main() {
	count := 3
	readConfigFile()
	initConnection()
	defer conn.Close()
	_, ver, err := conn.Ping(0)
	if err != nil {
		fmt.Println("Error contacting InfluxDB server:", err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Found InfluxDB server version", ver)
	}
	initHosts()

	list.Ping(count, func(statistics []*ping.Statistics, host *hosts.Host) {
		conn.Store(statistics, host)
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	for range stop {
		os.Exit(0)
	}
}

func readConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/pingflux/")
	viper.AddConfigPath("$HOME/.pingflux")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			panic(err)
		} else {
			// Config file found but another error was encountered
			panic(err)
		}
	}
}

func initHosts() {
	list = hosts.NewList()
	for remote, props := range viper.GetStringMap("hosts") {
		tags := map[string]string{}
		for tag, value := range props.(map[string]interface{}) {
			tags[tag] = value.(string)
		}
		if net.ParseIP(remote) != nil {
			list.AddIP(remote, tags)
		} else {
			list.AddHostname(remote, tags)
		}
	}
	fmt.Println("Found", list.Length(), "hosts in config file")
}

func initConnection() {
	addr := ""
	if viper.GetBool("datastore.influx.secure") {
		addr += "https://"
	} else {
		addr += "http://"
	}
	addr += viper.GetString("datastore.influx.hostname") + ":" + viper.GetString("datastore.influx.port")
	db := viper.GetString("datastore.influx.database")
	fmt.Printf("Connecting to %s on %s\n", db, addr)
	conn = influx.New(addr, db)
}
