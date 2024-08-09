package service

import (
	"fmt"
	"time"

	"github.com/tinkermode/internal/averager"
	"github.com/tinkermode/internal/fetcher"
	"github.com/tinkermode/internal/parser"
)

type TimeSeriesService interface {
	ProcessTimeSeries(startTime, endTime time.Time, callback func(time.Time, float64)) error
}

type timeSeries struct {
	fetcher  fetcher.Fetcher
	parser   parser.Parser
	averager averager.Averager
}

func NewTimeSeriesService() TimeSeriesService {
	return &timeSeries{
		fetcher:  fetcher.NewHTTPFetcher(""),
		parser:   parser.NewDataParser(),
		averager: averager.NewAverager(),
	}
}

func (s *timeSeries) ProcessTimeSeries(startTime, endTime time.Time, callback func(time.Time, float64)) error {
	url := fmt.Sprintf("https://tsserv.tinkermode.dev/data?begin=%s&end=%s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	s.fetcher = fetcher.NewHTTPFetcher(url)

	scanner, err := s.fetcher.Fetch()
	if err != nil {
		return fmt.Errorf("error fetching data: %v", err)
	}

	dataPoints, err := s.parser.Parse(scanner)
	if err != nil {
		return fmt.Errorf("error parsing data: %v", err)
	}

	done := make(chan struct{})
	go func() {
		s.averager.CalculateAverages(dataPoints, startTime, endTime, callback)
		close(done)
	}()

	<-done
	return nil
}
