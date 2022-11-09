package github_actions_usage_calculator

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	targetRange, err := DecideRange("2022-11")
	if err != nil {
		panic(err)
	}

	expectedStart := time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2022, 11, 30, 0, 0, 0, 0, time.UTC)

	if targetRange.Start != expectedStart {
		t.Errorf("expected %v, but got %v", expectedStart, targetRange.Start)
	}
	if targetRange.End != expectedEnd {
		t.Errorf("expected %v, but got %v", expectedEnd, targetRange.End)
	}
}
