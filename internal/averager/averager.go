package averager

import (
	"time"

	"github.com/tinkermode/internal/parser"
)

type Averager interface {
	CalculateAverages(dataPoints chan parser.DataPoint, startTime, endTime time.Time, callback func(time.Time, float64))
}

type dataAverager struct{}

func NewAverager() Averager {
	return &dataAverager{}
}

func (a *dataAverager) CalculateAverages(dataPoints chan parser.DataPoint, startTime, endTime time.Time, callback func(time.Time, float64)) {
	currentBucket := startTime
	var sum float64
	var count int

	for dp := range dataPoints {
		for dp.Timestamp.After(currentBucket.Add(time.Hour)) {
			if count > 0 {
				avg := sum / float64(count)
				callback(currentBucket, avg)
			}
			currentBucket = currentBucket.Add(time.Hour)
			sum = 0
			count = 0
		}

		sum += dp.Value
		count++
	}

	if count > 0 {
		avg := sum / float64(count)
		callback(currentBucket, avg)
	}
}
