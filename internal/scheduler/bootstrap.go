package scheduler

import (
	"context"
	"predictive-maintenance/internal/config"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/internal/store"
	"predictive-maintenance/pkg/logger"
	"time"
)

func Build(ctx context.Context, c config.Config, s *store.Store, devices *service.DeviceService, health *service.HealthService, maintenance *service.MaintenanceService, log *logger.Logger) *Runner {
	r := New(log)
	if !c.EnableScheduler {
		return r
	}
	r.Add(Job{Name: "offline-detector", Interval: c.SchedulerInterval, Run: (&OfflineDetector{Devices: devices, Timeout: c.OfflineTimeout, Log: log}).Run})
	r.Add(Job{Name: "maintenance-creator", Interval: c.MaintenanceInterval, Run: (&MaintenanceCreator{Service: maintenance, Log: log, SkipWeekend: true}).Run})
	r.Add(Job{Name: "health-evaluator", Interval: c.HealthInterval, Run: (&HealthEvaluator{Service: health, Log: log}).Run})
	r.Add(Job{Name: "report-generator", Interval: 24 * time.Hour, Run: (&ReportGenerator{Store: s, Log: log}).Run})
	return r
}
