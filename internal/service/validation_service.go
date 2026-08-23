package service

import (
	"fmt"
	"predictive-maintenance/internal/model"
	"strings"
	"time"
)

type ValidationService struct{}

func NewValidationService() *ValidationService { return &ValidationService{} }
func (v *ValidationService) Device(d model.Device) error {
	if strings.TrimSpace(d.Code) == "" {
		return fmt.Errorf("code is required")
	}
	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if d.InstallDate.After(time.Now().Add(time.Hour)) {
		return fmt.Errorf("install date cannot be future")
	}
	return nil
}
func (v *ValidationService) Batch(b model.SensorDataBatch) error {
	if b.DeviceID == 0 {
		return fmt.Errorf("deviceId is required")
	}
	if len(b.Items) == 0 {
		return fmt.Errorf("items are required")
	}
	if len(b.Items) > 1000 {
		return fmt.Errorf("too many readings")
	}
	for _, i := range b.Items {
		if !model.ValidSensorType(i.SensorType) {
			return fmt.Errorf("invalid sensor type %s", i.SensorType)
		}
	}
	return nil
}
func (v *ValidationService) Plan(p model.MaintenancePlan) error {
	if p.DeviceID == 0 {
		return fmt.Errorf("deviceId is required")
	}
	if p.Type != "periodic" && p.Type != "runtime" && p.Type != "condition_based" {
		return fmt.Errorf("invalid plan type")
	}
	return nil
}
