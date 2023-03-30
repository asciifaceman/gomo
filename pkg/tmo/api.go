package tmo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/asciifaceman/gomo/pkg/models"
)

const (
	HeaderUserAgent      = "Mozilla/5.0 (Windows NT 10.0; rv:111.0) Gecko/20100101 Firefox/111.0"
	HeaderAccept         = "application/json"
	HeaderAcceptLanguage = "en-US,en;q=0.5"
	HeaderAcceptEncoding = "gzip, deflate"
	HeaderContentType    = "application/x-www-form-urlencoded"
	HeaderConnection     = "keep-alive"

	URIFastmile = "fastmile_radio_status_web_app.cgi"
)

type Client struct {
	client   *http.Client
	Hostname string
}

func NewClient(hostname string) (*Client, error) {
	_, err := url.Parse(hostname)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Hostname: hostname,
		client:   http.DefaultClient,
	}

	return c, nil
}

// FetchRadioStatus ...
func (c *Client) FetchRadioStatus() (*models.FastmileRadioStatus, error) {
	reqURI, err := url.Parse(fmt.Sprintf("%s/%s", c.Hostname, URIFastmile))
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
	req.Header.Add("Referer", fmt.Sprintf("%s/web_whw", c.Hostname))

	resp, err := c.client.Do(req)
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
