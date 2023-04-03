package gotmo

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/asciifaceman/gomo/pkg/clio"
	"github.com/asciifaceman/gomo/pkg/models"
	"github.com/asciifaceman/gomo/pkg/status"
	"github.com/asciifaceman/gomo/pkg/tmo"
	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	DaemonAddr = ":2112"
)

/*
600 MHz: Band 71
700 MHz: Band 12
850 MHz: Band 5
1700/2100 MHz: Bands 4/66
1900 MHz: Band 2
*/

var (
	BandMap5G = map[string]float64{
		"n71":  0.6,
		"n41":  2.5,
		"n2":   3.4,
		"n77":  3.7,
		"n258": 24,
		"n261": 39,
		"n262": 47,
	}
	BandMapLTE = map[string]float64{
		"B71": 0.6,
		"B12": 0.7,
		"B5":  0.85,
		"B4":  1.7,
		"B66": 2.1,
		"B2":  1.9,
	}
	LabelCellID     = "cell_id"
	LabelBand       = "band"
	CurrentCellId5G = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "5g",
		Name:      "cell_id",
		Help:      "The current CellID of the 5G radio. GHz",
	})
	CurrentBand5G = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "5g",
		Name:      "band",
		Help:      "The current Band of the 5G radio. GHz",
	})
	SNR5G = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "5g",
		Name:      "snr",
		Help:      "The current SNR of the 5G radio at this point in time. dB",
	})
	RSRP5G = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "5g",
		Name:      "rsrp",
		Help:      "The current RSRP of the 5G radio at this point in time. dBm",
	})
	RSRQ5G = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "5g",
		Name:      "rsrq",
		Help:      "The current RSRQ of the 5G radio at this point in time. dBm",
	})
	CurrentCellIdLTE = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "lte",
		Name:      "cell_id",
		Help:      "The current CellID of the LTE radio. GHz",
	})
	CurrentBandLTE = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "lte",
		Name:      "band",
		Help:      "The current Band of the LTE radio. GHz",
	})
	RSSILTE = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "lte",
		Name:      "rssi",
		Help:      "The current RSSI of the LTE radio at this point in time. dBm",
	})
	SNRLTE = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "lte",
		Name:      "snr",
		Help:      "The current SNR of the LTE radio at this point in time. dB",
	})
	RSRPLTE = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "lte",
		Name:      "rsrp",
		Help:      "The current RSRP of the LTE radio at this point in time. dBm",
	})
	RSRQLTE = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "lte",
		Name:      "rsrq",
		Help:      "The current RSRQ of the LTE radio at this point in time. dB",
	})

	CellBytesSent = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "cell",
		Name:      "bytes_sent",
		Help:      "The number of bytes sent this uptime",
	})
	CellBytesRecv = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "cell",
		Name:      "bytes_received",
		Help:      "The number of bytes received this uptime",
	})
)

// Gotmo is the primary entrypoint for interactive operation
type Gotmo struct {
	Logger                *zap.SugaredLogger
	Server                *http.Server
	Trashcan              *tmo.Trashcan
	Printer               *clio.Printer
	Status                *status.Status
	PingTargets           []string
	Signals               chan os.Signal
	FastmileReturnChannel chan *models.FastmileReturn
	PingReturnChannel     chan *models.PingReportReturn
	PingWorkChannel       chan string
	HttpErrorChannel      chan error
	PingWorkerCount       int
}

func NewGotmo(hostname string, timeout int, pingTargets []string, pingWorkerCount int) (*Gotmo, error) {
	requestTimeout := time.Duration(timeout) * time.Second

	t, terr := tmo.NewTrashcan(hostname, requestTimeout)
	if terr != nil {
		return nil, terr
	}

	p := clio.NewPrinter(clio.DefaultHeaderWidth, clio.DefaultKVWidth, clio.DefaultIndent)

	s := status.NewStatus(5, pingTargets)

	g := &Gotmo{
		Trashcan:              t,
		Printer:               p,
		Status:                s,
		FastmileReturnChannel: make(chan *models.FastmileReturn, 1),
		PingReturnChannel:     make(chan *models.PingReportReturn, len(pingTargets)),
		PingWorkChannel:       make(chan string, len(pingTargets)),
		Signals:               make(chan os.Signal, 1),
		HttpErrorChannel:      make(chan error, 1),
		PingTargets:           pingTargets,
		PingWorkerCount:       pingWorkerCount,
		Server: &http.Server{
			Addr: DaemonAddr,
		},
	}

	signal.Notify(g.Signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return g, nil
}

// RegisterPrometheus
func (g *Gotmo) RegisterPrometheus() {
	prometheus.MustRegister(SNR5G)
	prometheus.MustRegister(RSRP5G)
	prometheus.MustRegister(RSRQ5G)
	prometheus.MustRegister(RSSILTE)
	prometheus.MustRegister(SNRLTE)
	prometheus.MustRegister(RSRPLTE)
	prometheus.MustRegister(RSRQLTE)
	prometheus.MustRegister(CellBytesSent)
	prometheus.MustRegister(CellBytesRecv)
	prometheus.MustRegister(CurrentBand5G)
	prometheus.MustRegister(CurrentCellId5G)
	prometheus.MustRegister(CurrentBandLTE)
	prometheus.MustRegister(CurrentCellIdLTE)

}

// Daemon runs Gomo as a server continuously gathering metrics for prometheus
func (g *Gotmo) Daemon() error {
	err := g.SetupLogger()
	if err != nil {
		return err
	}

	g.RegisterPrometheus()

	var wg sync.WaitGroup

	g.Logger.Info("Starting webserver...")

	http.Handle("/metrics", promhttp.Handler())
	go g.RunHTTPServer()

	g.Logger.Info(fmt.Sprintf("Listening on %s", DaemonAddr))

	for {
		select {
		case <-g.Signals:
			g.Logger.Info("Received exit signal, shutting down")
			g.StopHTTPServer()
			close(g.HttpErrorChannel)
			ret := fmt.Errorf("received exit signal")
			for err := range g.HttpErrorChannel {
				ret = fmt.Errorf("%s; %s", ret.Error(), err.Error())
			}
			return ret
		case <-time.After(15 * time.Second):
			g.Logger.Info("Scraping data...")

			go g.Trashcan.FetchRadioStatusAsync(&wg, g.FastmileReturnChannel)
			wg.Add(1)
		case ret := <-g.FastmileReturnChannel:
			g.Logger.Info("received fastmile data")
			SNR5G.Set(float64(ret.Body.Cell5GStats[0].Stat.SNRCurrent))
			RSRP5G.Set(float64(ret.Body.Cell5GStats[0].Stat.RSRPCurrent))
			RSRQ5G.Set(float64(ret.Body.Cell5GStats[0].Stat.RSRQCurrent))
			cur5GCellID, err := strconv.ParseFloat(ret.Body.Cell5GStats[0].Stat.PhysicalCellID, 64)
			if err != nil {
				cur5GCellID = 0
			}
			CurrentCellId5G.Set(cur5GCellID)
			CurrentBand5G.Set(BandMap5G[ret.Body.Cell5GStats[0].Stat.Band])

			RSSILTE.Set(float64(ret.Body.CellLTEStats[0].Stat.RSSICurrent))
			SNRLTE.Set(float64(ret.Body.CellLTEStats[0].Stat.SNRCurrent))
			RSRPLTE.Set(float64(ret.Body.CellLTEStats[0].Stat.RSRPCurrent))
			RSRQLTE.Set(float64(ret.Body.CellLTEStats[0].Stat.RSRQCurrent))
			curLTECellID, err := strconv.ParseFloat(ret.Body.CellLTEStats[0].Stat.PhysicalCellID, 64)
			if err != nil {
				curLTECellID = 0
			}
			CurrentCellIdLTE.Set(curLTECellID)
			CurrentBandLTE.Set(BandMapLTE[ret.Body.CellLTEStats[0].Stat.Band])

			CellBytesSent.Set(float64(ret.Body.CellularStats[0].BytesSent))
			CellBytesRecv.Set(float64(ret.Body.CellularStats[0].BytesReceived))
		}
	}

	return nil
}

// RunHTTPServer launches the metrics http endpoint
func (g *Gotmo) RunHTTPServer() {
	if err := g.Server.ListenAndServe(); err != nil {
		g.HttpErrorChannel <- err
	}
}

// StopHTTPServer stops the http server
func (g *Gotmo) StopHTTPServer() {
	g.Server.Close()
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := g.Server.Shutdown(ctx); err != nil {
	//	g.HttpErrorChannel <- err
	//}
}

// AlignEntry runs Gomo as a fast continuous poll and returns data for visualization to align an antenna
func (g *Gotmo) AlignEntry() *models.FastmileRadioStatus {
	var wg sync.WaitGroup
	g.FastmileReturnChannel = make(chan *models.FastmileReturn, 1)

	go g.Trashcan.FetchRadioStatusAsync(&wg, g.FastmileReturnChannel)
	wg.Add(1)

	wg.Wait()
	close(g.FastmileReturnChannel)
	radioStatus := <-g.FastmileReturnChannel
	return radioStatus.Body

}

// CLIEntry runs Gomo as a single pass and displays the responses
func (g *Gotmo) CLIEntry(pretty bool) {
	var wg sync.WaitGroup

	// Start radio status report
	go g.Trashcan.FetchRadioStatusAsync(&wg, g.FastmileReturnChannel)
	wg.Add(1)

	// Launch ping reports
	for i := 0; i < g.PingWorkerCount; i++ {
		go g.Status.PingAsync(&wg, g.PingWorkChannel, g.PingReturnChannel)
		wg.Add(1)
	}

	for _, hostname := range g.PingTargets {
		g.PingWorkChannel <- hostname
	}
	close(g.PingWorkChannel)

	wg.Wait()
	close(g.FastmileReturnChannel)
	close(g.PingReturnChannel)
	radioStatus := <-g.FastmileReturnChannel

	if pretty {
		g.Printer.PrintHeader("LTE")
		g.Printer.PrintKVIndent("RSSI", radioStatus.Body.CellLTEStats[0].Stat.RSSICurrent)
		g.Printer.PrintKVIndent("SNR", radioStatus.Body.CellLTEStats[0].Stat.SNRCurrent)
		g.Printer.PrintKVIndent("RSRP", radioStatus.Body.CellLTEStats[0].Stat.RSRPCurrent)
		g.Printer.PrintKVIndent("RSRQ", radioStatus.Body.CellLTEStats[0].Stat.RSRQCurrent)
		g.Printer.PrintKVIndent("Band", radioStatus.Body.CellLTEStats[0].Stat.Band)
		g.Printer.PrintKVIndent("CellID", radioStatus.Body.CellLTEStats[0].Stat.PhysicalCellID)
		g.Printer.PrintHeader("5G")
		g.Printer.PrintKVIndent("SNR", radioStatus.Body.Cell5GStats[0].Stat.SNRCurrent)
		g.Printer.PrintKVIndent("RSRP", radioStatus.Body.Cell5GStats[0].Stat.RSRPCurrent)
		g.Printer.PrintKVIndent("RSRQ", radioStatus.Body.Cell5GStats[0].Stat.RSRQCurrent)
		g.Printer.PrintKVIndent("Band", radioStatus.Body.Cell5GStats[0].Stat.Band)
		g.Printer.PrintKVIndent("CellID", radioStatus.Body.Cell5GStats[0].Stat.PhysicalCellID)
	} else {
		spew.Dump(radioStatus.Body)
	}

	g.Printer.PrintHeader("Ping Tests")
	for report := range g.PingReturnChannel {
		if report.Error != nil {
			g.Printer.PrintKV("Error", report.Error.Error())
			continue
		}
		g.Printer.PrintHeader(report.Body.Hostname)
		g.Printer.PrintKVIndent("Packets Sent", report.Body.PacketsSent)
		g.Printer.PrintKVIndent("Packet Loss", report.Body.PacketLoss)
		g.Printer.PrintKVIndent("Avg. Response Time", report.Body.AvgResponseTime)
	}

	g.Printer.PrintHeader("5G")
	g.Printer.PrintKVIndent("SNR Quality", radioStatus.Body.Cell5GStats[0].Stat.SNRQuality(0, 1))
	g.Printer.PrintKVIndent("RSRP Quality", radioStatus.Body.Cell5GStats[0].Stat.RSRPQuality(0, 1))
	g.Printer.PrintKVIndent("RSRQ Quality", radioStatus.Body.Cell5GStats[0].Stat.RSRQQuality(0, 1))
	g.Printer.PrintHeader("LTE")
	g.Printer.PrintKVIndent("SNR Quality", radioStatus.Body.CellLTEStats[0].Stat.SNRQuality(0, 1))
	g.Printer.PrintKVIndent("RSRP Quality", radioStatus.Body.CellLTEStats[0].Stat.RSRPQuality(0, 1))
	g.Printer.PrintKVIndent("RSRQ Quality", radioStatus.Body.CellLTEStats[0].Stat.RSRQQuality(0, 1))
}

// SetupLogger will attach a logger to Gotmo for daemon use
func (g *Gotmo) SetupLogger() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	g.Logger = logger.Sugar()
	return nil
}
