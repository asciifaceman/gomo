// package status provides a few network status checks like ping and DNS queries to
// infer if it is reasonably healthy
package status

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

var (
	DefaultPingHosts = []string{"www.google.com", "github.com"}
)

type Status struct {
	PingCount       int
	TermChannel     chan interface{}
	ErrorChannel    chan error
	PingStatChannel chan *probing.Statistics
	Signals         chan os.Signal
	PingHosts       []string
}

type PingReport struct {
	Hostname            string
	PacketLoss          float64
	AverageResponseTime time.Duration
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

func (s *Status) Run() ([]*PingReport, error) {
	var wg sync.WaitGroup
	var report []*PingReport

	for _, hostname := range s.PingHosts {
		wg.Add(1)
		go s.Ping(&wg, hostname)
	}

	wg.Wait()
	close(s.PingStatChannel)

	for stat := range s.PingStatChannel {
		r := &PingReport{
			Hostname:            stat.Addr,
			PacketLoss:          stat.PacketLoss,
			AverageResponseTime: stat.AvgRtt,
		}
		report = append(report, r)
	}

	return report, nil
}

// Ping sets up a goroutine pinger for the given hostname and n count
func (s *Status) Ping(wg *sync.WaitGroup, hostname string) {
	defer wg.Done()
	pinger, err := probing.NewPinger(hostname)
	if err != nil {
		s.ErrorChannel <- err
	}
	pinger.Count = s.PingCount

	go func() {
		for range s.Signals {
			pinger.Stop()
		}
	}()

	err = pinger.Run()
	if err != nil {
		s.ErrorChannel <- err
	}
	stats := pinger.Statistics()
	s.PingStatChannel <- stats
}
