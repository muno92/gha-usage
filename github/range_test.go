package github

import (
	"testing"
	"time"
)

func TestNewRange(t *testing.T) {
	tests := []struct {
		name                string
		startDate           string
		endDate             string
		expectedErrorExists bool
		expectedRange       Range
	}{
		{
			name:                "Valid",
			startDate:           "2022-01-01",
			endDate:             "2022-01-31",
			expectedErrorExists: false,
			expectedRange: Range{
				Start: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2022, 1, 31, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:                "Start Date is invalid format",
			startDate:           "2022-01",
			endDate:             "2022-01-31",
			expectedErrorExists: true,
			expectedRange:       Range{},
		},
		{
			name:                "Start Date is not yyyy-mm-dd",
			startDate:           "2022/01/01",
			endDate:             "2022-01-31",
			expectedErrorExists: true,
			expectedRange:       Range{},
		},
		{
			name:                "End Date is invalid format",
			startDate:           "2022-01-01",
			endDate:             "test",
			expectedErrorExists: true,
			expectedRange:       Range{},
		},
		{
			name:                "End Date is before start date",
			startDate:           "2022-01-01",
			endDate:             "2021-01-31",
			expectedErrorExists: true,
			expectedRange:       Range{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetRange, err := NewRange(tt.startDate, tt.endDate)

			errorExists := err != nil
			if tt.expectedErrorExists != errorExists {
				t.Errorf("expected error exists is %v, got %v\n%v", tt.expectedErrorExists, errorExists, err)
			}

			if targetRange != tt.expectedRange {
				t.Errorf("expected %v, but got %v", tt.expectedRange, targetRange)
			}
		})
	}
}
