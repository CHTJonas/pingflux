# Pingflux

pingflux is a small application written in Go that measures the ICMP round-trip time and packet loss between the local host and a set of remote endpoints. It then stores these data in an InfluxDB database. IP addresses are pinged individually whereas hostnames are first resolved, each of their IPv4 and IPv6 addresses then being pinged concurrently.

## Usage

On Linux and macOS you will need to run the binary with `CAP_NET_RAW` privileges in order to send ICMP packets (running as root is strongly discouraged). You can do this by running `sudo setcap cap_net_raw=+ep /usr/local/bin/pingflux` in a terminal.

pingflux loads its configuration settings by looking for a `config.yml` file in either `/etc/pingflux/`, `$HOME/.pingflux` or the current working directory in that order of precedence. Remote hosts are listed in groups which share a common set of key-value tags in InfluxDB. The number of pings sent to each host and the interval in seconds at which these pings are sent to each host are configurable using the `count` and `interval` options respectively. A count of 10 at an interval of 60 means that every 10 pings will be sent in quick succession to each host every 60 seconds.

To minimise concurrent HTTP requests and improve performance, data are bundled together and submitted to InfluxDB in batches, the size of which is controlled by the `batch-size` option. The default of 25 is optimised for very large numbers of hosts however you may wish to reduce this substantially if you are only pinging a few hosts.

The `datastore` section specifies the backend database that pingflux will use to store data, although at this time only InfluxDB is supported. The hostname, port and database name can all be configured together with the `secure` value which may be set to `true` to enable HTTPS.

An example config file might be:

```yaml
options:
  count: 10
  interval: 60
  batch-size: 25

datastore:
  influx:
    hostname: 127.0.0.1
    port: 8086
    path: /
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

## Installation

Pre-built binaries for a variety of operating systems and architectures are available to download from [GitHub Releases](https://github.com/CHTJonas/pingflux/releases). If you wish to compile from source then you will need a suitable [Go toolchain installed](https://golang.org/doc/install). After that just clone the project using Git and run Make! Cross-compilation is easy in Go so by default we build for all targets and place the resulting executables in `./bin`:

```bash
git clone https://github.com/CHTJonas/pingflux.git
cd pingflux
make clean && make all
```

For Linux users there is a [sample systemd service file](https://github.com/CHTJonas/pingflux/blob/master/pingflux.service) available which you can place at `/etc/systemd/system/pingflux.service` and then activate:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now pingflux
```

## Copyright

pingflux is licensed under the [BSD 2-Clause License](https://opensource.org/licenses/BSD-2-Clause).

Copyright (c) 2019–2022 Charlie Jonas.
