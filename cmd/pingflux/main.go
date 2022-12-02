package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
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

func init() {
	if os.Geteuid() == 0 {
		log.Println("WARNING: Running pingflux as root is strongly discouraged")
		log.Println("WARNING: You should set the CAP_NET_RAW privilege instead")
	}
}

func main() {
	log.Println("pingflux version", version)

	count, interval := readConfig()
	initConnection()
	defer connection.Close()
	rtt, ver, err := connection.Ping(0)
	if err != nil {
		log.Println("Error contacting InfluxDB server:", err.Error())
		os.Exit(1)
	}
	log.Printf("Found InfluxDB server version %s. HTTP request RTT %s\n", ver, rtt)
	initHosts()

	size := viper.GetInt("options.batch-size")
	resultChan := make(chan *hosts.Result, size)
	resultsArrPool := sync.Pool{
		New: func() interface{} {
			arr := make([]*hosts.Result, size)
			return &arr
		},
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGTERM)

	reload := make(chan os.Signal, 1)
	signal.Notify(reload, syscall.SIGHUP)

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
		case <-reload:
			log.Println("Reloading config...")
			if err := viper.ReadInConfig(); err != nil {
				log.Println("Failed to read config file:", err)
				continue
			}
			hostList.Stop()
			initHosts()
			go hostList.Ping(count, interval, resultChan)
		case <-stop:
			log.Println("Shutting down...")
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
		dur := b.Duration()
		log.Printf("Failed to store data in InfluxDB: %s. Will retry after %s\n", err, dur)
		<-time.After(dur)
	}
}

func readConfig() (int, int) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/pingflux/")
	viper.AddConfigPath("$HOME/.pingflux")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read config file:", err)
		os.Exit(125)
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
	log.Println("Found", hostList.Length(), "hosts in config file")
}

func initConnection() {
	var addr string
	if viper.GetBool("datastore.influx.secure") {
		addr += "https://"
	} else {
		addr += "http://"
	}
	addr += viper.GetString("datastore.influx.hostname")
	addr += ":"
	addr += viper.GetString("datastore.influx.port")
	addr += viper.GetString("datastore.influx.path")
	db := viper.GetString("datastore.influx.database")
	username := viper.GetString("datastore.influx.username")
	password := viper.GetString("datastore.influx.password")
	userAgent := fmt.Sprintf("pingflux/%s Go/%s (+https://github.com/CHTJonas/pingflux)",
		version, strings.TrimPrefix(runtime.Version(), "go"))
	log.Println("Connecting to", addr)
	var err error
	connection, err = influx.New(addr, db, username, password, userAgent)
	if err != nil {
		log.Println("Error connecting to InfluxDB:", err)
		os.Exit(500)
	}
}

func init() {
	if os.Getenv("JOURNAL_STREAM") != "" {
		log.Default().SetFlags(0)
	}
}
