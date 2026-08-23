package scheduler

import (
	"context"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/pkg/logger"
	"time"
)

type OfflineDetector struct {
	Devices *service.DeviceService
	Timeout time.Duration
	Log     *logger.Logger
}

func (o *OfflineDetector) Run(ctx context.Context) error {
	if o.Timeout <= 0 {
		o.Timeout = 5 * time.Minute
	}
	if e := o.Devices.MarkOffline(ctx, o.Timeout); e != nil {
		return e
	}
	o.Log.Debug("offline device scan finished")
	return nil
}
