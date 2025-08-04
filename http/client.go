package http

import (
	"io"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/httpclient"
)

type HttpService struct {
	client *httpclient.Client
}

func NewHttpService(initialTimeout, maxTimeout, maxJitterInterval, timeout time.Duration, exponentFactor float64, retries int) *HttpService {
	backoff := heimdall.NewExponentialBackoff(initialTimeout, maxTimeout, exponentFactor, maxJitterInterval)
	retrier := heimdall.NewRetrier(backoff)

	return &HttpService{
		client: httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
			httpclient.WithRetrier(retrier),
			httpclient.WithRetryCount(retries),
		),
	}
}

func (hs *HttpService) MakeRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	res, err := hs.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
