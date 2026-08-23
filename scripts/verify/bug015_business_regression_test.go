package verify

import (
	"context"
	"path/filepath"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/repository"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/internal/store"
	"testing"
	"time"
)

func TestBug015_BusinessRegression(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "offline.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if _, err = s.CreateDevice(model.Device{Code: "O1", Name: "offline"}); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	_, err = s.DB.Exec("INSERT INTO sensor_data(device_id,sensor_type,value,is_abnormal,abnormal_reason,collected_at,created_at) VALUES(?,?,?,?,?,?,?)", 1, "temperature", 20, false, "", now.Add(-time.Hour).Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		t.Fatal(err)
	}
	devices := service.NewDeviceService(repository.NewDeviceRepository(s.DB))
	devices.UseSensors(repository.NewSensorRepository(s.DB))
	if err = devices.MarkOffline(context.Background(), 5*time.Minute); err != nil {
		t.Fatal(err)
	}
	device, err := devices.Get(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if device.Status != model.DeviceOffline {
		t.Fail()
	}
}
