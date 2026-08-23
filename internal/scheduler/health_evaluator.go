package scheduler

import (
	"context"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/pkg/logger"
)

type HealthEvaluator struct {
	Service *service.HealthService
	Devices []uint64
	Log     *logger.Logger
}

func (h *HealthEvaluator) Run(ctx context.Context) error {
	for _, id := range h.Devices {
		point, e := h.Service.Evaluate(ctx, id)
		if e != nil {
			h.Log.Warn("health evaluation skipped", map[string]any{"deviceId": id, "error": e.Error()})
			continue
		}
		if point.Score < 60 {
			h.Log.Warn("device health critical", map[string]any{"deviceId": id, "score": point.Score})
		}
	}
	return nil
}
func (h *HealthEvaluator) SetDevices(ids []uint64) { h.Devices = append([]uint64(nil), ids...) }
