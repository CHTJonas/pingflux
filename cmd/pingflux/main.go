package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chtjonas/pingflux/internal/hosts"
	"github.com/chtjonas/pingflux/internal/influx"
	"github.com/spf13/viper"
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
		fmt.Println("Error pinging InfluxDB server: ", err.Error())
	} else {
		fmt.Println("Got reply from InfluxDB server: ", ver)
	}
	initHosts()

	for e := list.Hosts.Front(); e != nil; e = e.Next() {
		r := rand.Intn(1000)
		d := time.Duration(r)
		time.Sleep(d * time.Millisecond)
		go setupPinger(e.Value.(*hosts.Host), count)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	for range stop {
		os.Exit(0)
	}
}

func setupPinger(host *hosts.Host, count int) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			statistics := host.Ping(count)
			conn.Store(statistics, host)
		}
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
