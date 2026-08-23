package service

import (
	"context"
	"fmt"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/repository"
	"strings"
	"time"
)

type DeviceService struct {
	Repo    *repository.DeviceRepository
	Sensors *repository.SensorRepository
}

func NewDeviceService(r *repository.DeviceRepository) *DeviceService { return &DeviceService{Repo: r} }
func (s *DeviceService) UseSensors(r *repository.SensorRepository)   { s.Sensors = r }
func (s *DeviceService) Register(ctx context.Context, d model.Device) (model.Device, error) {
	if strings.TrimSpace(d.Code) == "" || strings.TrimSpace(d.Name) == "" {
		return d, fmt.Errorf("device code and name are required")
	}
	if d.Type == "" {
		d.Type = "unknown"
	}
	return s.Repo.DBCreate(ctx, d)
}
func (s *DeviceService) Get(ctx context.Context, id uint64) (model.Device, error) {
	return s.Repo.Find(ctx, id)
}
func (s *DeviceService) SetStatus(ctx context.Context, id uint64, status int) error {
	return s.Repo.UpdateStatus(ctx, id, status)
}
func (s *DeviceService) MarkOffline(ctx context.Context, timeout time.Duration) error {
	devices, e := s.Repo.ListByStatus(ctx, model.DeviceOnline)
	if e != nil {
		return e
	}
	deadline := time.Now().Add(-timeout)
	if s.Sensors != nil {
		for _, d := range devices {
			latest, e := s.Sensors.LatestForDevice(ctx, d.ID)
			if e != nil {
				return e
			}
			if !latest.CreatedAt.IsZero() && latest.CreatedAt.Before(deadline) {
				if e = s.SetStatus(ctx, d.ID, model.DeviceOffline); e != nil {
					return e
				}
			}
		}
		return nil
	}
	for _, d := range devices {
		if d.UpdatedAt.Before(deadline) {
			if e = s.SetStatus(ctx, d.ID, model.DeviceOffline); e != nil {
				return e
			}
		}
	}
	return nil
}
func (s *DeviceService) HealthLabel(score float64) string {
	switch {
	case score >= 90:
		return "excellent"
	case score >= 75:
		return "good"
	case score >= 60:
		return "warning"
	default:
		return "critical"
	}
}
