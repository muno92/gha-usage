package github_actions_usage_calculator

import (
	"fmt"
	"time"
)

type Range struct {
	Start time.Time
	End   time.Time
}

func DecideRange(targetMonth string) (*Range, error) {
	start, err := time.Parse(time.RFC3339, fmt.Sprintf("%s-01T00:00:00Z", targetMonth))
	if err != nil {
		return nil, fmt.Errorf("fail to parse target month with %s: %w", targetMonth, err)
	}
	return &Range{
		Start: start,
		End:   time.Date(start.Year(), start.Month(), 30, 0, 0, 0, 0, time.UTC),
	}, nil
}
