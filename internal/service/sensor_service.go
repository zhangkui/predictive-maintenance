package service

import (
	"context"
	"fmt"
	"predictive-maintenance/internal/detector"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/repository"
	"time"
)

type SensorService struct {
	Repo     *repository.SensorRepository
	Abnormal *repository.AbnormalRepository
	Detector detector.CompositeDetector
	Configs  *ConfigService
}

func NewSensorService(r *repository.SensorRepository, a *repository.AbnormalRepository) *SensorService {
	return &SensorService{Repo: r, Abnormal: a, Detector: detector.CompositeDetector{Detectors: []detector.Detector{detector.ThresholdDetector{}, detector.RateDetector{}, detector.TrendDetector{Window: 8}}}}
}
func (s *SensorService) UseConfigs(configs *ConfigService) { s.Configs = configs }

func (s *SensorService) Ingest(ctx context.Context, batch model.SensorDataBatch) ([]model.SensorData, []model.AbnormalRecord, error) {
	if batch.DeviceID == 0 || len(batch.Items) == 0 {
		return nil, nil, fmt.Errorf("empty sensor batch")
	}
	latest, err := s.Repo.Latest(ctx, batch.DeviceID)
	if err != nil {
		return nil, nil, err
	}
	history := map[string][]detector.Sample{}
	for _, item := range latest {
		history[item.SensorType] = append(history[item.SensorType], detector.Sample{SensorType: item.SensorType, Value: item.Value, CollectedAt: item.CollectedAt})
	}
	saved := []model.SensorData{}
	abnormals := []model.AbnormalRecord{}
	for _, item := range batch.Items {
		if !model.ValidSensorType(item.SensorType) {
			return nil, nil, fmt.Errorf("unsupported sensor type %s", item.SensorType)
		}
		t := item.CollectedAt
		if t.IsZero() {
			t = time.Now().UTC()
		}
		sample := detector.Sample{SensorType: item.SensorType, Value: item.Value, CollectedAt: t}
		cfg := detector.Config{MinValue: 0, MaxValue: 100, CriticalMin: -10, CriticalMax: 120}
		if s.Configs != nil {
			configured, configErr := s.Configs.List(ctx, batch.DeviceID)
			if configErr != nil {
				return saved, abnormals, configErr
			}
			for _, itemConfig := range configured {
				if itemConfig.SensorType == item.SensorType {
					cfg = detector.Config{MinValue: 0, MaxValue: itemConfig.MaxValue, CriticalMin: itemConfig.CriticalMin, CriticalMax: itemConfig.CriticalMax, AlarmThreshold: itemConfig.AlarmThreshold}
					break
				}
			}
		}
		result := s.Detector.Detect(sample, history[item.SensorType], cfg)
		data := model.SensorData{DeviceID: batch.DeviceID, SensorType: item.SensorType, Value: item.Value, IsAbnormal: result.Abnormal, AbnormalReason: result.Reason, CollectedAt: t}
		if err = s.Repo.Save(ctx, data); err != nil {
			return saved, abnormals, err
		}
		saved = append(saved, data)
		history[item.SensorType] = append(history[item.SensorType], sample)
		if result.Abnormal {
			a := model.AbnormalRecord{DeviceID: batch.DeviceID, SensorType: item.SensorType, DetectedValue: item.Value, NormalRange: fmt.Sprintf("%.3f-%.3f", cfg.MinValue, cfg.MaxValue), Severity: result.Severity}
			if _, err = s.Abnormal.Create(ctx, a); err != nil {
				return saved, abnormals, err
			}
			abnormals = append(abnormals, a)
		}
	}
	return saved, abnormals, nil
}
func (s *SensorService) History(ctx context.Context, device uint64, sensor string, from, to time.Time) ([]model.SensorData, error) {
	return s.Repo.History(ctx, device, sensor, from, to, 1000)
}
func (s *SensorService) Latest(ctx context.Context, device uint64) ([]model.SensorData, error) {
	return s.Repo.Latest(ctx, device)
}
