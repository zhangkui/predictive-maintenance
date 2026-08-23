package scheduler

import (
	"context"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/pkg/logger"
	"time"
)

type MaintenanceCreator struct {
	Service     *service.MaintenanceService
	Log         *logger.Logger
	SkipWeekend bool
}

func (m *MaintenanceCreator) Run(ctx context.Context) error {
	if m.SkipWeekend && time.Now().Weekday() >= time.Saturday {
		m.Log.Info("maintenance generation skipped on weekend")
		return nil
	}
	n, e := m.Service.GenerateTasks(ctx)
	if e == nil {
		m.Log.Info("maintenance tasks generated", map[string]any{"count": n})
	}
	return e
}
func (m *MaintenanceCreator) SetSkipWeekend(v bool) { m.SkipWeekend = v }
func (m *MaintenanceCreator) ScheduledDate(now time.Time) time.Time {
	return now
}
