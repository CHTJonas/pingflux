package main

import (
	"container/list"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
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
				l := flushData(resultList)
				go connection.Store(l)
			}
		case <-stop:
			fmt.Println("Received shutdown signal...")
			connection.Store(flushData(resultList))
			os.Exit(0)
		}
	}
}

func flushData(resultList *list.List) *list.List {
	l := list.New()
	l.PushBackList(resultList)
	resultList.Init()
	return l
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
	var err error
	connection, err = influx.New(addr, db)
	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}
}
