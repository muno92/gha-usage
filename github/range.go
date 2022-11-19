package github

import (
	"fmt"
	"time"
)

type Range struct {
	Start time.Time
	End   time.Time
}

func DecideRange(targetMonth string) (Range, error) {
	start, err := time.Parse(time.RFC3339, fmt.Sprintf("%s-01T00:00:00Z", targetMonth))
	if err != nil {
		return Range{}, fmt.Errorf("fail to parse target month with %s: %w", targetMonth, err)
	}
	return Range{
		Start: start,
		End:   time.Date(start.Year(), start.Month()+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1),
	}, nil
}
