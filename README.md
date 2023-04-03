# GoMo
Tmobile trashcan datalogging for the original trashcan.

If you know, you know.

# Usage

```
 $ ./gomo --help
Fetch and log trashcan data to a long term data store for analyzing

Usage:
  gomo [command]

Available Commands:
  align       Continuously fetch data and display timeseries CLI charts
  completion  Generate the autocompletion script for the specified shell
  daemon      Daemonized Gomo which will continuously run
  help        Help about any command
  show        Do a single fetch and display
  version     retrieve version and build info for gomo

Flags:
      --config string     config file (default is $HOME/.gomo.yaml)
  -h, --help              help for gomo
  -u, --hostname string   hostname of your tmobile trashcan (default "http://192.168.12.1")
  -p, --targets strings   List of hostnames to target with ping test (default [www.google.com,github.com])
  -s, --timeout int       timeout in seconds for outbound requests (default 15)
  -t, --toggle            Help message for toggle
  -w, --workers int       number of workers for pingers (default 2)

Use "gomo [command] --help" for more information about a command.
```

## Quick Run
The quick run mode is accessible via the `show` command, and does a single pass and returns current values along with the results of a ping test.

```
=== GOMO 0.0.3 =================================================
=== LTE ========================================================
  RSSI:                       -69
  SNR:                         -1
  RSRP:                      -106
  RSRQ:                       -18
  Band:                       B66
  CellID:                     redacted
=== 5G =========================================================
  SNR:                         -2
  RSRP:                      -111
  RSRQ:                       -18
  Band:                       n41
  CellID:                      redacted
=== Ping Tests =================================================
=== www.google.com =============================================
  Packets Sent:                 5
  Packet Loss:                  0
  Avg. Response Time:  98.89631ms
=== github.com =================================================
  Packets Sent:                 5
  Packet Loss:                  0
  Avg. Response Time:104.296611ms
=== 5G =========================================================
  SNR Quality:               0.45
  RSRP Quality:0.11428571428571428
  RSRQ Quality:               0.2
=== LTE ========================================================
  SNR Quality:              0.475
  RSRP Quality:0.2571428571428571
  RSRQ Quality:               0.2
```

## Alignment
Alignment mode, accessible via `align` shows a continuous time series graph of LTE and 5G metrics to help align an antenna.

![Alignment](static/alignment.png)

## Daemon datalogging
Daemon mode, accessible via `daemon` is a background process - meant to be run by a systemd unit. This continuously scrapes data from the trashcan and surfaces it on a /metrics endpoint for prometheus to scrape.

There is a rough prometheus/grafana setup configured with a dashboard meant for this data

![Grafana](static/grafana_dash.png)

# TODO
## High Priority
* Export more trashcan metrics to prometheus
* Stub trashcan for tests with a fake json response
* Improve daemon exit

## Low Priority
* Lots of cleaning up and making code more coherent after rapid dev
* Tighten up prometheus/grafana deployment
* Docker container for gomo in docker-compose for quick launch
* Configurable port for gomo metrics
* Add internet speedtest metrics exporter
* Add ping metrics exporter


# Dependencies
* spf13/cobra
* spf13/viper
* davecgh/go-spew
* [TermUI](github.com/gizak/termui/v3)
* [Prometheus Client](github.com/prometheus/client_golang)
* [prometheus-community/pro-bing](https://github.com/prometheus-community/pro-bing)
  * provides a tool to ping remote hosts with windows support

# Authors
* Charles Corbett <github.com/asciifaceman>


# Thanks
* Karl Q