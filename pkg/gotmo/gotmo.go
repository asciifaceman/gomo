package gotmo

import (
	"sync"
	"time"

	"github.com/asciifaceman/gomo/pkg/clio"
	"github.com/asciifaceman/gomo/pkg/models"
	"github.com/asciifaceman/gomo/pkg/status"
	"github.com/asciifaceman/gomo/pkg/tmo"
)

// Gotmo is the primary entrypoint for interactive operation
type Gotmo struct {
	Trashcan              *tmo.Trashcan
	Printer               *clio.Printer
	Status                *status.Status
	PingTargets           []string
	FastmileReturnChannel chan *models.FastmileReturn
	PingReturnChannel     chan *models.PingReportReturn
	PingWorkChannel       chan string
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
		PingTargets:           pingTargets,
		PingWorkerCount:       pingWorkerCount,
	}

	return g, nil
}

// Daemon runs Gomo as a server continuously gathering metrics for prometheus
func (g *Gotmo) Daemon() error {
	return nil
}

// CLIEntry runs Gomo as a single pass and displays the responses
func (g *Gotmo) CLIEntry() {
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
}
