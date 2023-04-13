package tmo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/asciifaceman/gomo/pkg/models"
)

const (
	DefaultTimeout = 15 * time.Second

	HeaderUserAgent      = "Mozilla/5.0 (Windows NT 10.0; rv:111.0) Gecko/20100101 Firefox/111.0"
	HeaderAccept         = "application/json"
	HeaderAcceptLanguage = "en-US,en;q=0.5"
	HeaderAcceptEncoding = "gzip, deflate"
	HeaderContentType    = "application/x-www-form-urlencoded"
	HeaderConnection     = "keep-alive"

	URIFastmile = "fastmile_radio_status_web_app.cgi"
)

// Trashcan defines some methods for interacting with the tmo trashcan
type Trashcan struct {
	client   *http.Client
	Hostname string
}

// NewTrashcan returns a configured Trashcan client
func NewTrashcan(hostname string, timeout time.Duration) (*Trashcan, error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	_, err := url.Parse(hostname)
	if err != nil {
		return nil, err
	}

	t := &Trashcan{
		Hostname: hostname,
		client: &http.Client{
			Timeout: timeout,
		},
	}

	return t, nil
}

// FetchRadioStatusAsync is for running in a goroutine, calls FetchRadioStatus
func (t *Trashcan) FetchRadioStatusAsync(wg *sync.WaitGroup, ret chan<- *models.FastmileReturn) {
	defer wg.Done()

	data, err := t.FetchRadioStatus()
	response := &models.FastmileReturn{
		Body:  data,
		Error: err,
	}

	ret <- response
}

// FetchRadioStatus fetches the radio status cgi page and returns a payload of
// radio data
func (t *Trashcan) FetchRadioStatus() (*models.FastmileRadioStatus, error) {
	reqURI, err := url.Parse(fmt.Sprintf("%s/%s", t.Hostname, URIFastmile))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", reqURI.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", HeaderUserAgent)
	req.Header.Add("Accept", HeaderAccept)
	req.Header.Add("Accept-Language", HeaderAcceptLanguage)
	req.Header.Add("Accept-Encoding", HeaderAcceptEncoding)
	req.Header.Add("Content-Type", HeaderContentType)
	req.Header.Add("Connection", HeaderConnection)
	req.Header.Add("Referer", fmt.Sprintf("%s/web_whw", t.Hostname))

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fastmile := &models.FastmileRadioStatus{}

	err = json.Unmarshal(body, fastmile)

	return fastmile, err
}
