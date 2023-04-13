package clients

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/asciifaceman/gomo/pkg/models"
	"github.com/asciifaceman/gomo/pkg/tmo"
)

// Polling is a client for injecting a return channel and leaving it run in a loop
type Polling struct {
	Trashcan              *tmo.Trashcan
	PollTimeout           time.Duration
	Signals               chan os.Signal
	FastmileReturnChannel chan *models.FastmileReturn
}

func NewPolling(hostname string, timeout int, pollFrequency int) (*Polling, error) {
	if pollFrequency < 1 {
		return nil, fmt.Errorf("poll frequency too fast, may overrun")
	}

	requestTimeout := time.Duration(timeout) * time.Second
	pollTimeout := time.Duration(pollFrequency) * time.Second

	t, err := tmo.NewTrashcan(hostname, requestTimeout)
	if err != nil {
		return nil, err
	}

	p := &Polling{
		Trashcan:              t,
		PollTimeout:           pollTimeout,
		Signals:               make(chan os.Signal),
		FastmileReturnChannel: make(chan *models.FastmileReturn, 1),
	}

	// In case the caller forgets
	signal.Notify(p.Signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return p, nil
}

// Start launches a poller which will return results on ret and terminate on signal
func (p *Polling) Start(ret chan *models.FastmileReturn, signal chan interface{}) {
	var wg sync.WaitGroup

	t := time.NewTicker(p.PollTimeout).C

	go p.Trashcan.FetchRadioStatusAsync(&wg, p.FastmileReturnChannel)
	wg.Add(1)

	for {
		select {
		case <-p.Signals:
			wg.Wait()
			return
		case <-signal:
			wg.Wait()
			return
		case <-t:
			go p.Trashcan.FetchRadioStatusAsync(&wg, p.FastmileReturnChannel)
			wg.Add(1)
		case d := <-p.FastmileReturnChannel:
			ret <- d

		}
	}

}
