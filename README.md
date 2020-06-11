# Pingflux

pingflux is a small application written in Go that measures the ICMP round-trip time and packet loss between the local host and a set of remote endpoints. It then stores these data in an InfluxDB database.

## Usage

pingflux loads its configuration settings by looking for a `config.yml` file in either `/etc/pingflux/`, `$HOME/.pingflux` or the current working directory in that order of precedence. Remote hosts are listed in groups which share a common set of key-value tags in InfluxDB. The number of pings sent to each host and the interval in seconds at which these pings are sent to each host are configurable using the `count` and `interval` options respectively. A count of 10 at an interval of 60 means that every 10 pings will be sent in quick succession to each host every 60 seconds.

The `datastore` section specifies the backend database that pingflux will use to store data, although at this time only InfluxDB is supported. The hostname, port and database name can all be configured together with the `secure` value which may be set to `true` to enable HTTPS.

An example config file might be:

```yaml
options:
  count: 10
  interval: 60

datastore:
  influx:
    hostname: 127.0.0.1
    port: 8086
    username: user
    password: pass
    database: pingflux
    secure: false

groups:
  - tag: value
    another-tag: different-value
    hosts:
      host.example.com
      1.1.1.1

  - key: value
    other-key: other-value
    hosts:
      host.example.net
      9.9.9.9
```

## Compiling

It should be relatively simple to checkout and build the code, assuming you have a suitable [Go toolchain installed](https://golang.org/doc/install). Running the following commands in a terminal will compile binaries for various operating systems and processor architectures and place them in directories in `./bin`:

```bash
git checkout https://github.com/CHTJonas/pingflux.git
make clean && make all
```

---

### Copyright

pingflux is licensed under the [BSD 2-Clause License](https://opensource.org/licenses/BSD-2-Clause).

Copyright (c) 2019–2020 Charlie Jonas.
