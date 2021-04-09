package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/chtjonas/pingflux/internal/hosts"
	"github.com/chtjonas/pingflux/internal/influx"
	"github.com/cloudflare/backoff"
	"github.com/spf13/viper"
)

var hostList *hosts.List
var connection *influx.Connection

// Software version defaults to the value below but is overridden by the compiler in Makefile.
var version = "dev-edge"

const homepage = "https://github.com/CHTJonas/pingflux"

func init() {
	if os.Geteuid() == 0 {
		fmt.Println("WARNING: Running pingflux as root is strongly discouraged")
		fmt.Println("WARNING: You should set the CAP_NET_RAW privilege instead")
	}
}

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

	size := 25
	resultChan := make(chan *hosts.Result, size)
	resultsArrPool := sync.Pool{
		New: func() interface{} {
			arr := make([]*hosts.Result, size)
			return &arr
		},
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	go hostList.Ping(count, interval, resultChan)

	n := 0
	resultsArrPtr := resultsArrPool.Get().(*[]*hosts.Result)
	for {
		select {
		case result := <-resultChan:
			(*resultsArrPtr)[n] = result
			n++
			if n == size {
				go func(ptr *[]*hosts.Result) {
					storeData(ptr)
					for i := 0; i < size; i++ {
						(*ptr)[i] = nil
					}
					resultsArrPool.Put(ptr)
				}(resultsArrPtr)
				n = 0
				resultsArrPtr = resultsArrPool.Get().(*[]*hosts.Result)
			}
		case <-stop:
			fmt.Println("Received shutdown signal...")
			storeData(resultsArrPtr)
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
	userAgent := fmt.Sprintf("pingflux/%s (+%s)", version, homepage)
	fmt.Printf("Connecting to %s on %s\n", db, addr)
	var err error
	connection, err = influx.New(addr, db, username, password, userAgent)
	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}
}
