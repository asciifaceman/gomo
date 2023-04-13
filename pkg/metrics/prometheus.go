package metrics

import "github.com/prometheus/client_golang/prometheus"

/*
	5G Prometheus Metrics
*/

var Metric5GCurrentCellID = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "5g",
	Name:      "cell_id",
	Help:      "The current CellID of the 5G radio. GHz",
})

var Metric5GCurrentBand = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "5g",
	Name:      "band",
	Help:      "The current Band of the 5G radio. GHz",
})

var Metric5GCurrentSNR = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "5g",
	Name:      "snr",
	Help:      "The current SNR of the 5G radio at this point in time. dB",
})

var Metric5GCurrentRSRP = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "5g",
	Name:      "rsrp",
	Help:      "The current RSRP of the 5G radio at this point in time. dBm",
})

var Metric5GCurrentRSRQ = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "5g",
	Name:      "rsrq",
	Help:      "The current RSRQ of the 5G radio at this point in time. dBm",
})

var Metric5GCurrentDownlinkARFCN = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "5g",
	Name:      "downlink_nr_arfcn",
	Help:      "The absolute radio frequency channel number of teh radio at this point in time",
})

// Metrics5G is a convenience var for 5G metric gauges
var Metrics5G = map[string]prometheus.Gauge{
	"cell_id": Metric5GCurrentCellID,
	"band":    Metric5GCurrentBand,
	"snr":     Metric5GCurrentSNR,
	"rsrp":    Metric5GCurrentRSRP,
	"rsrq":    Metric5GCurrentRSRQ,
	"arfcn":   Metric5GCurrentDownlinkARFCN,
}

/*
	LTE Prometheus Metrics
*/

var MetricLTECurrentCellID = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "lte",
	Name:      "cell_id",
	Help:      "The current CellID of the 5G radio. GHz",
})

var MetricLTECurrentBand = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "lte",
	Name:      "band",
	Help:      "The current Band of the 5G radio. GHz",
})

var MetricLTECurrentSNR = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "lte",
	Name:      "snr",
	Help:      "The current SNR of the 5G radio at this point in time. dB",
})

var MetricLTECurrentRSRP = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "lte",
	Name:      "rsrp",
	Help:      "The current RSRP of the 5G radio at this point in time. dBm",
})

var MetricLTECurrentRSRQ = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "lte",
	Name:      "rsrq",
	Help:      "The current RSRQ of the 5G radio at this point in time. dBm",
})

var MetricLTECurrentDownlinkARFCN = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "lte",
	Name:      "downlink_nr_arfcn",
	Help:      "The absolute radio frequency channel number of teh radio at this point in time",
})

var MetricLTECurrentRSSI = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "lte",
	Name:      "rssi",
	Help:      "The absolute radio frequency channel number of teh radio at this point in time",
})

// MetricsLTE is a convenience var for LTE metric gauges
var MetricsLTE = map[string]prometheus.Gauge{
	"cell_id": MetricLTECurrentCellID,
	"band":    MetricLTECurrentBand,
	"snr":     MetricLTECurrentSNR,
	"rsrp":    MetricLTECurrentRSRP,
	"rsrq":    MetricLTECurrentRSRQ,
	"rssi":    MetricLTECurrentRSSI,
	"arfcn":   MetricLTECurrentDownlinkARFCN,
}

/*
	Misc Prometheus Metrics
*/

var MetricConnectionStatus = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "wan",
	Name:      "connection_status",
	Help:      "The reported connection status of the device. integer bool",
})

var MetricCellularBytesSent = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "cell",
	Name:      "bytes_sent",
	Help:      "The reported number of bytes sent over the cellular connection this uptime",
})

var MetricCellularBytesRecv = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: "gomo",
	Subsystem: "cell",
	Name:      "bytes_received",
	Help:      "The reported number of bytes received over the cellular connection this uptime",
})

// MetricsMisc is a convenience var for misc metric gauges
var MetricsMisc = map[string]prometheus.Gauge{
	"connection_status": MetricConnectionStatus,
	"bytes_sent":        MetricCellularBytesSent,
	"bytes_recv":        MetricCellularBytesRecv,
}
