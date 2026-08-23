package verify

import (
	"predictive-maintenance/internal/scheduler"
	"testing"
	"time"
)

func TestBug014_BusinessRegression(t *testing.T) {
	creator := &scheduler.MaintenanceCreator{SkipWeekend: true}
	date := creator.ScheduledDate(time.Date(2026, 8, 23, 9, 0, 0, 0, time.UTC))
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		t.Fail()
	}
}
