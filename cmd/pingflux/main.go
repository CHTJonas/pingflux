package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chtjonas/pingflux/internal/hosts"
	"github.com/chtjonas/pingflux/internal/influx"
	"github.com/cloudflare/backoff"
	"github.com/spf13/viper"
)

var hostList *hosts.List
var connection *influx.Connection

func main() {
	count, interval := readConfig()
	initConnection()
	defer connection.Close()
	_, ver, err := connection.Ping(0)
	if err != nil {
		fmt.Println("Error contacting InfluxDB server:", err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Found InfluxDB server version", ver)
	}
	initHosts()

	resultsArr := make([]*hosts.Result, 10)
	n := 0
	resultChan := make(chan *hosts.Result, 3)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	hostList.Ping(count, interval, resultChan)
	for {
		select {
		case result := <-resultChan:
			resultsArr[n] = result
			n = n + 1
			if n > 9 {
				storeData(&resultsArr)
				resultsArr = make([]*hosts.Result, 10)
			}
		case <-stop:
			fmt.Println("Received shutdown signal...")
			storeData(&resultsArr)
			os.Exit(0)
		}
	}
}

func storeData(resultsArrPtr *[]*hosts.Result) {
	b := backoff.New(48*time.Hour, 2*time.Minute)
	for {
		err := connection.Store(resultsArrPtr)
		if err == nil {
			break
		}
		fmt.Println("Failed to store data:", err)
		<-time.After(b.Duration())
	}
}

func readConfig() (int, int) {
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
	return viper.GetInt("options.count"), viper.GetInt("options.interval")
}

func initHosts() {
	hostList = hosts.NewList()
	for _, g := range viper.Get("groups").([]interface{}) {
		group := g.(map[interface{}]interface{})
		tags := map[string]string{}
		for k, v := range group {
			if k.(string) != "hosts" {
				tags[k.(string)] = v.(string)
			}
		}
		for _, remote := range strings.Split(group["hosts"].(string), " ") {
			if net.ParseIP(remote) != nil {
				hostList.AddIP(remote, tags)
			} else {
				hostList.AddHostname(remote, tags)
			}
		}
	}
	hostList.Shuffle()
	fmt.Println("Found", hostList.Length(), "hosts in config file")
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
	username := viper.GetString("datastore.influx.username")
	password := viper.GetString("datastore.influx.password")
	fmt.Printf("Connecting to %s on %s\n", db, addr)
	var err error
	connection, err = influx.New(addr, db, username, password)
	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}
}
