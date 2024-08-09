package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type DataPoint struct {
	Timestamp time.Time
	Value     float64
}

type Parser interface {
	Parse(scanner *bufio.Scanner) (chan DataPoint, error)
}

type DataParser struct{}

func NewDataParser() Parser {
	return &DataParser{}
}

func (p *DataParser) Parse(scanner *bufio.Scanner) (chan DataPoint, error) {
	dataPoints := make(chan DataPoint)

	go func() {
		defer close(dataPoints)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) != 2 {
				fmt.Fprintln(os.Stderr, "Invalid data format")
				continue
			}

			timestamp, err := time.Parse(time.RFC3339, parts[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid timestamp format:", err)
				continue
			}

			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid value format:", err)
				continue
			}

			dataPoints <- DataPoint{Timestamp: timestamp, Value: value}
		}
	}()

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading data: %v", err)
	}

	return dataPoints, nil
}
