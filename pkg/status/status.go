// package status provides a few network status checks like ping and DNS queries to
// infer if it is reasonably healthy
package status

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/asciifaceman/gomo/pkg/models"
	probing "github.com/prometheus-community/pro-bing"
)

var (
	DefaultPingHosts   = []string{"www.google.com", "github.com"}
	DefaultWorkerCount = 2
)

type Status struct {
	PingCount       int
	TermChannel     chan interface{}
	ErrorChannel    chan error
	PingStatChannel chan *probing.Statistics
	Signals         chan os.Signal
	PingHosts       []string
}

func NewStatus(pingCount int, pingHosts []string) *Status {
	s := &Status{
		TermChannel:     make(chan interface{}, 1),
		ErrorChannel:    make(chan error, 1),
		Signals:         make(chan os.Signal, 1),
		PingStatChannel: make(chan *probing.Statistics, len(pingHosts)),
		PingCount:       pingCount,
		PingHosts:       pingHosts,
	}

	signal.Notify(s.Signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	return s
}

// PingAsync ...
func (s *Status) PingAsync(wg *sync.WaitGroup, work chan string, ret chan *models.PingReportReturn) {
	defer wg.Done()
	for hostname := range work {
		result := &models.PingReportReturn{}
		pinger, err := probing.NewPinger(hostname)
		if err != nil {
			result.Error = err
			ret <- result
			continue
		}
		pinger.Count = s.PingCount
		pinger.Timeout = 15 * time.Second

		go func() {
			for range s.Signals {
				pinger.Stop()
			}
		}()

		err = pinger.Run()
		if err != nil {
			result.Error = err
			ret <- result
			continue
		}

		stats := pinger.Statistics()

		result.Body = &models.PingReport{
			Hostname:        stats.Addr,
			PacketsSent:     stats.PacketsSent,
			PacketLoss:      stats.PacketLoss,
			AvgResponseTime: stats.AvgRtt,
		}
		ret <- result
	}
}
