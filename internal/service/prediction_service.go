package service

import (
	"context"
	"predictive-maintenance/internal/detector"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/repository"
	"predictive-maintenance/internal/util"
	"time"
)

type PredictionService struct{ Sensors *repository.SensorRepository }

func NewPredictionService(s *repository.SensorRepository) *PredictionService {
	return &PredictionService{s}
}
func (s *PredictionService) Predict(ctx context.Context, id uint64) (model.Prediction, error) {
	data, e := s.Sensors.History(ctx, id, model.SensorVibration, time.Now().Add(-30*24*time.Hour), time.Now(), 500)
	if e != nil {
		return model.Prediction{}, e
	}
	history := make([]detector.Sample, 0, len(data))
	for _, d := range data {
		history = append(history, detector.Sample{SensorType: d.SensorType, Value: d.Value, CollectedAt: d.CollectedAt})
	}
	if len(history) < 2 {
		return model.Prediction{DeviceID: id, RemainingHours: 720, Confidence: .2, Trend: "insufficient-data", GeneratedAt: time.Now().UTC()}, nil
	}
	if len(history) > 2 {
		history = history[len(history)-2:]
	}
	trend := detector.TrendDirection(history)
	slope := history[len(history)-1].Value - detector.Forecast(history, 1)
	remaining := 720.0
	if slope > 0 {
		remaining = util.Clamp((100-history[len(history)-1].Value)/slope, 1, 720)
	}
	confidence := util.Clamp(float64(len(history))/100, .2, .95)
	return model.Prediction{DeviceID: id, RemainingHours: remaining, FailureProbability: 1 - confidence, Confidence: confidence, Trend: trend, GeneratedAt: time.Now().UTC()}, nil
}
func (s *PredictionService) Risk(p model.Prediction) string {
	if p.RemainingHours < 24 {
		return "critical"
	}
	if p.RemainingHours < 168 {
		return "high"
	}
	return "normal"
}
