package clients

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/asciifaceman/gomo/pkg/metrics"
	"github.com/asciifaceman/gomo/pkg/models"
	"github.com/asciifaceman/gomo/pkg/tmo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	DefaultPort    = 2112
	DefaultTimeout = 15
)

// Daemon is the central daemon which surfaces metrics from the tmo trashcan
type Daemon struct {
	Logger                *zap.SugaredLogger
	Server                *http.Server
	Trashcan              *tmo.Trashcan
	PollTimeout           time.Duration
	FastmileReturnChannel chan *models.FastmileReturn
	HttpErrorChannel      chan error
	Signals               chan os.Signal
}

// New returns a newly configured daemon ready to start
// A port or timeout of 0 will use default values (2112 & 15 respectively)
func NewDaemon(hostname string, port int, timeout int) (*Daemon, error) {
	requestTimeout := time.Duration(timeout) * time.Second

	t, err := tmo.NewTrashcan(hostname, requestTimeout)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf(":%d", port)

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	g := &Daemon{
		Logger:      logger.Sugar(),
		Trashcan:    t,
		PollTimeout: time.Duration(15),
		Server: &http.Server{
			Addr: addr,
		},
		FastmileReturnChannel: make(chan *models.FastmileReturn),
		HttpErrorChannel:      make(chan error, 1),
		Signals:               make(chan os.Signal, 1),
	}

	signal.Notify(g.Signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return g, nil
}

// RegisterMetrics registers various metrics with prometheus_client to be surfaced
func (d *Daemon) RegisterMetrics() {
	for _, v := range metrics.Metrics5G {
		prometheus.MustRegister(v)
	}

	for _, v := range metrics.MetricsLTE {
		prometheus.MustRegister(v)
	}

	for _, v := range metrics.MetricsMisc {
		prometheus.MustRegister(v)
	}
}

func (d *Daemon) Run() error {
	d.RegisterMetrics()

	var wg sync.WaitGroup

	d.Logger.Info("Starting webserver...")

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", d.Hello)

	go d.BackgroundHTTPServer()

	d.Logger.Info(fmt.Sprintf("Listening on %s, entering runtime loop", d.Server.Addr))

	for {
		select {
		case <-d.Signals:
			d.Logger.Info("Received exit signal, shutting down")
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			err := d.Server.Shutdown(ctx)
			cancel()
			if err == http.ErrServerClosed {
				return nil
			}
			d.Logger.Info("Waiting on subroutines...")
			wg.Wait()
			return err
		case <-time.After(d.PollTimeout * time.Second):
			d.Logger.Info("Scraping data...")

			go d.Trashcan.FetchRadioStatusAsync(&wg, d.FastmileReturnChannel)
			wg.Add(1)
		case err := <-d.HttpErrorChannel:
			wg.Wait()
			return err

		case ret := <-d.FastmileReturnChannel:
			if ret.Error != nil {
				d.Logger.Errorw("Errored scraping trashcan", "error", ret.Error.Error())
				continue
			}
			if ret.Body == nil {
				d.Logger.Errorw("received empty body without error")
				continue
			}
			d.Logger.Info("Received fastmile data, updating metrics")

			metrics.Metrics5G["cell_id"].Set(ret.Stat5G().ID())
			metrics.Metrics5G["band"].Set(ret.Stat5G().Band64())
			metrics.Metrics5G["snr"].Set(ret.Stat5G().SNRCurrent)
			metrics.Metrics5G["rsrp"].Set(ret.Stat5G().RSRPCurrent)
			metrics.Metrics5G["rsrq"].Set(ret.Stat5G().RSRQCurrent)
			metrics.Metrics5G["arfcn"].Set(ret.Stat5G().DownlinkNRARFCN)

			metrics.MetricsLTE["cell_id"].Set(ret.StatLTE().ID())
			metrics.MetricsLTE["band"].Set(ret.StatLTE().Band64())
			metrics.MetricsLTE["snr"].Set(ret.StatLTE().SNRCurrent)
			metrics.MetricsLTE["rsrp"].Set(ret.StatLTE().RSRPCurrent)
			metrics.MetricsLTE["rsrq"].Set(ret.StatLTE().RSRQCurrent)
			metrics.MetricsLTE["arfcn"].Set(ret.StatLTE().DownlinkEarfcn)

			metrics.MetricsMisc["connection_status"].Set(ret.Status())
			metrics.MetricsMisc["bytes_sent"].Set(ret.BytesSent())
			metrics.MetricsMisc["bytes_recv"].Set(ret.BytesRecv())

		}
	}

}

func (d *Daemon) BackgroundHTTPServer() {
	if err := d.Server.ListenAndServe(); err != nil {
		d.HttpErrorChannel <- err
	}
}

func (d *Daemon) Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "ok")
}
