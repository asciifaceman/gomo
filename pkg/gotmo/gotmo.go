package gotmo

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
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

var (
	LabelCellID = "cell_id"
	LabelBand   = "band"
	cellLabels  = []string{LabelCellID, LabelBand}
	SNR5G       = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "5g",
		Name:      "snr",
		Help:      "The current SNR of the 5G radio at this point in time. dB",
	}, cellLabels)
	RSRP5G = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "5g",
		Name:      "rsrp",
		Help:      "The current RSRP of the 5G radio at this point in time. dBm",
	}, cellLabels)
	RSRQ5G = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "5g",
		Name:      "rsrq",
		Help:      "The current RSRQ of the 5G radio at this point in time. dBm",
	}, cellLabels)
	RSSILTE = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "lte",
		Name:      "rssi",
		Help:      "The current RSSI of the LTE radio at this point in time. dBm",
	}, cellLabels)
	SNRLTE = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "lte",
		Name:      "snr",
		Help:      "The current SNR of the LTE radio at this point in time. dB",
	}, cellLabels)
	RSRPLTE = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "lte",
		Name:      "rsrp",
		Help:      "The current RSRP of the LTE radio at this point in time. dBm",
	}, cellLabels)
	RSRQLTE = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gomo",
		Subsystem: "lte",
		Name:      "rsrq",
		Help:      "The current RSRQ of the LTE radio at this point in time. dB",
	}, cellLabels)
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
				ret = fmt.Errorf("%w; %w", ret, err)
			}
			return ret
		case <-time.After(15 * time.Second):
			g.Logger.Info("Scraping data...")

			go g.Trashcan.FetchRadioStatusAsync(&wg, g.FastmileReturnChannel)
			wg.Add(1)
		case ret := <-g.FastmileReturnChannel:
			g.Logger.Info("received fastmile data")
			SNR5G.With(prometheus.Labels{
				LabelCellID: ret.Body.Cell5GStats[0].Stat.PhysicalCellID,
				LabelBand:   ret.Body.Cell5GStats[0].Stat.Band,
			}).Set(float64(ret.Body.Cell5GStats[0].Stat.SNRCurrent))
			RSRP5G.With(prometheus.Labels{
				LabelCellID: ret.Body.Cell5GStats[0].Stat.PhysicalCellID,
				LabelBand:   ret.Body.Cell5GStats[0].Stat.Band,
			}).Set(float64(ret.Body.Cell5GStats[0].Stat.RSRPCurrent))
			RSRQ5G.With(prometheus.Labels{
				LabelCellID: ret.Body.Cell5GStats[0].Stat.PhysicalCellID,
				LabelBand:   ret.Body.Cell5GStats[0].Stat.Band,
			}).Set(float64(ret.Body.Cell5GStats[0].Stat.RSRQCurrent))

			RSSILTE.With(prometheus.Labels{
				LabelCellID: ret.Body.Cell5GStats[0].Stat.PhysicalCellID,
				LabelBand:   ret.Body.Cell5GStats[0].Stat.Band,
			}).Set(float64(ret.Body.CellLTEStats[0].Stat.RSSICurrent))
			SNRLTE.With(prometheus.Labels{
				LabelCellID: ret.Body.Cell5GStats[0].Stat.PhysicalCellID,
				LabelBand:   ret.Body.Cell5GStats[0].Stat.Band,
			}).Set(float64(ret.Body.CellLTEStats[0].Stat.SNRCurrent))
			RSRPLTE.With(prometheus.Labels{
				LabelCellID: ret.Body.Cell5GStats[0].Stat.PhysicalCellID,
				LabelBand:   ret.Body.Cell5GStats[0].Stat.Band,
			}).Set(float64(ret.Body.CellLTEStats[0].Stat.RSRPCurrent))
			RSRQLTE.With(prometheus.Labels{
				LabelCellID: ret.Body.Cell5GStats[0].Stat.PhysicalCellID,
				LabelBand:   ret.Body.Cell5GStats[0].Stat.Band,
			}).Set(float64(ret.Body.CellLTEStats[0].Stat.RSRQCurrent))
		}
	}

	return nil
}

/*
	g.StopHTTPServer()
	close(g.HttpErrorChannel)
	var finalErr error
	for err := range g.HttpErrorChannel {
		fmt.Println(err)
		finalErr = err
	}
	return finalErr
*/

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
