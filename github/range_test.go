package github

import (
	"testing"
	"time"
)

func TestDecideRange(t *testing.T) {
	tests := []struct {
		name          string
		targetMonth   string
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{
			name:          "Month with 31 days",
			targetMonth:   "2022-01",
			expectedStart: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2022, 1, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "Non leap year February",
			targetMonth:   "2022-02",
			expectedStart: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "Leap year February",
			targetMonth:   "2024-02",
			expectedStart: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "April",
			targetMonth:   "2022-04",
			expectedStart: time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2022, 4, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "June",
			targetMonth:   "2022-06",
			expectedStart: time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2022, 6, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "September",
			targetMonth:   "2022-09",
			expectedStart: time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2022, 9, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "November",
			targetMonth:   "2022-11",
			expectedStart: time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2022, 11, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetRange, err := DecideRange(tt.targetMonth)
			if err != nil {
				panic(err)
			}

			if targetRange.Start != tt.expectedStart {
				t.Errorf("expected %v, but got %v", tt.expectedStart, targetRange.Start)
			}
			if targetRange.End != tt.expectedEnd {
				t.Errorf("expected %v, but got %v", tt.expectedEnd, targetRange.End)
			}
		})
	}
}
