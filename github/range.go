package github

import (
	"fmt"
	"time"
)

type Range struct {
	Start time.Time
	End   time.Time
}

func NewRange(startDate string, endDate string) (Range, error) {
	start, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", startDate))
	if err != nil {
		return Range{}, fmt.Errorf("fail to parse start date with %s: %w", startDate, err)
	}

	end, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", endDate))
	if err != nil {
		return Range{}, fmt.Errorf("fail to parse end date with %s: %w", startDate, err)
	}

	if end.Before(start) {
		return Range{}, fmt.Errorf("end date (%s) must be later than or equal to the start date (%s)", endDate, startDate)
	}

	return Range{Start: start, End: end}, nil
}
