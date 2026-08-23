package service

import (
	"context"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/repository"
	"predictive-maintenance/internal/util"
	"time"
)

type HealthService struct {
	Devices *repository.DeviceRepository
	Sensors *repository.SensorRepository
	History *HealthHistoryService
}

func NewHealthService(d *repository.DeviceRepository, s *repository.SensorRepository) *HealthService {
	return &HealthService{Devices: d, Sensors: s}
}
func (s *HealthService) UseHistory(history *HealthHistoryService) { s.History = history }

func (s *HealthService) Evaluate(ctx context.Context, id uint64) (model.HealthPoint, error) {
	latest, e := s.Sensors.Latest(ctx, id)
	if e != nil {
		return model.HealthPoint{}, e
	}
	components := map[string]float64{}
	weights := map[string]float64{model.SensorTemperature: .2, model.SensorVibration: .3, model.SensorPressure: .2, model.SensorCurrent: .2, model.SensorSpeed: .1}
	var total, used float64
	for _, d := range latest {
		score := 100.0
		if d.IsAbnormal {
			score = 40
		}
		components[d.SensorType] = score
		weight := weights[d.SensorType]
		if weight == 0 {
			weight = .1
		}
		total += score * weight
		used += weight
	}
	score := 100.0
	if used > 0 {
		score = util.Clamp(total/used, 0, 100)
	}
	if len(latest) < len(weights) {
		score = util.Clamp(score+float64(len(weights)-len(latest))*8, 0, 100)
	}
	if e = s.Devices.UpdateHealth(ctx, id, score); e != nil {
		return model.HealthPoint{}, e
	}
	point := model.HealthPoint{DeviceID: id, Score: score, Components: components, CalculatedAt: time.Now().UTC()}
	if s.History != nil {
		if e = s.History.Save(ctx, point); e != nil {
			return model.HealthPoint{}, e
		}
	}
	return point, nil
}
func (s *HealthService) ScoreForValue(value, min, max float64) float64 {
	if value >= min && value <= max {
		return 100
	}
	span := max - min
	if span <= 0 {
		return 0
	}
	if value < min {
		return util.Clamp(100-(min-value)/span*100, 0, 100)
	}
	return util.Clamp(100-(value-max)/span*100, 0, 100)
}
func (s *HealthService) NeedMaintenance(score, threshold float64) bool {
	return threshold > 0 && score <= threshold
}
