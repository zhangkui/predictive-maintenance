package scheduler

import (
	"context"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/pkg/logger"
)

type AnomalyScanner struct {
	Sensors *service.SensorService
	Tickets *service.TicketService
	Log     *logger.Logger
}

func (a *AnomalyScanner) Scan(ctx context.Context, devices []uint64) error {
	for _, id := range devices {
		data, e := a.Sensors.Latest(ctx, id)
		if e != nil {
			return e
		}
		for _, sample := range data {
			if !sample.IsAbnormal {
				continue
			}
			abnormal := a.recordFromSample(id, sample)
			if _, e = a.Tickets.Create(ctx, abnormal); e != nil {
				a.Log.Warn("ticket creation failed", map[string]any{"deviceId": id, "error": e.Error()})
			}
		}
	}
	return nil
}

func (a *AnomalyScanner) recordFromSample(device uint64, sample model.SensorData) model.AbnormalRecord {
	severity := model.SeveritySerious
	if sample.AbnormalReason != "" && len(sample.AbnormalReason) > 80 {
		severity = model.SeverityCritical
	}
	return model.AbnormalRecord{ID: sample.ID, DeviceID: device, SensorType: sample.SensorType, DetectedValue: sample.Value, Severity: severity, NormalRange: "configured"}
}
