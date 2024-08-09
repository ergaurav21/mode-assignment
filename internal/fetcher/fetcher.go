package fetcher

import (
	"bufio"
	"fmt"
	"net/http"
)

type Fetcher interface {
	Fetch() (*bufio.Scanner, error)
}

type httpFetcher struct {
	url string
}

func NewHTTPFetcher(url string) Fetcher {
	return &httpFetcher{url: url}
}

func (f *httpFetcher) Fetch() (*bufio.Scanner, error) {
	resp, err := http.Get(f.url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code")
	}

	return bufio.NewScanner(resp.Body), nil
}
