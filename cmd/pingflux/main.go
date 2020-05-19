package main

import (
	"container/list"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/chtjonas/pingflux/internal/hosts"
	"github.com/chtjonas/pingflux/internal/influx"
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

	resultList := list.New()
	resultChan := make(chan *hosts.Result, 3)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	hostList.Ping(count, interval, resultChan)
	for {
		select {
		case result := <-resultChan:
			resultList.PushBack(result)
			if resultList.Len() > 10 {
				storeData(resultList)
			}
		case <-stop:
			storeData(resultList)
			os.Exit(0)
		}
	}
}

func storeData(resultList *list.List) {
	l := list.New()
	l.PushBackList(resultList)
	go connection.Store(l)
	resultList.Init()
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
	for remote, props := range viper.GetStringMap("hosts") {
		tags := map[string]string{}
		for tag, value := range props.(map[string]interface{}) {
			tags[tag] = value.(string)
		}
		if net.ParseIP(remote) != nil {
			hostList.AddIP(remote, tags)
		} else {
			hostList.AddHostname(remote, tags)
		}
	}
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
	fmt.Printf("Connecting to %s on %s\n", db, addr)
	connection = influx.New(addr, db)
}
