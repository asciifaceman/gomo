# GoMo
Tmobile trashcan datalogging for the original trashcan

# Usage 
```
Fetch and log trashcan data to a long term data store for analyzing

Usage:
  gomo [flags]
  gomo [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     retrieve version and build info for gomo

Flags:
      --config string     config file (default is $HOME/.gomo.yaml)
  -h, --help              help for gomo
  -u, --hostname string   hostname of your tmobile trashcan (default "http://192.168.12.1")
      --pretty            Whether to just print radio data or not
  -t, --toggle            Help message for toggle

Use "gomo [command] --help" for more information about a command.
```

```
$ go run main.go --pretty
Stats for http://192.168.12.1
Connected: 1
IP6: redacted
== LTE =======================
  RSSI:   -70
  SNR:      0
  RSRP:  -106
  RSRQ:   -17
  Band:   B66
  CellID: redacted
== 5G =======================
  SNR:     -1
  RSRP:  -115
  RSRQ:   -20
  Band:   n41
  CellID: redacted
```

# TODO
* Hook up to time series datasource

# Author
Charles Corbett <github.com/asciifaceman>

# Thanks
* Karl Q