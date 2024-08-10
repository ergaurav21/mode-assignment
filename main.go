package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tinkermode/internal/service"
)

const timeFormat = time.RFC3339

func main() {

	startTime, endTime := validation()

	tsService := service.NewTimeSeriesService()
	err := tsService.ProcessTimeSeries(startTime, endTime, func(t time.Time, avg float64) {
		fmt.Printf("%s %.4f\n", t.Format(timeFormat), avg)
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error processing time series:", err)
		os.Exit(1)
	}

}

func validation() (startTime, endTime time.Time) {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Please use go run main.go <start-time> <end-time>")
		os.Exit(1)
	}

	startTime, err := time.Parse(timeFormat, os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid start time format:", err)
		os.Exit(1)
	}

	endTime, err = time.Parse(timeFormat, os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid end time format:", err)
		os.Exit(1)
	}

	return startTime, endTime
}
