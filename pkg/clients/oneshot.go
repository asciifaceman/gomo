package clients

import (
	"sync"
	"time"

	"github.com/asciifaceman/gomo/pkg/models"
	"github.com/asciifaceman/gomo/pkg/tmo"
)

// OneShot is a client for making single requests without polling
type OneShot struct {
	Trashcan              *tmo.Trashcan
	FastmileReturnChannel chan *models.FastmileReturn
}

func NewOneShot(hostname string, timeout int) (*OneShot, error) {
	requestTimeout := time.Duration(timeout) * time.Second

	t, err := tmo.NewTrashcan(hostname, requestTimeout)
	if err != nil {
		return nil, err
	}

	o := &OneShot{
		Trashcan:              t,
		FastmileReturnChannel: make(chan *models.FastmileReturn, 1),
	}

	return o, nil
}

// Fetch requests a set of metrics from the trashcan
func (o *OneShot) Fetch() *models.FastmileReturn {
	var wg sync.WaitGroup

	go o.Trashcan.FetchRadioStatusAsync(&wg, o.FastmileReturnChannel)
	wg.Add(1)
	wg.Wait()

	ret := <-o.FastmileReturnChannel
	return ret
}
