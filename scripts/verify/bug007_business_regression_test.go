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

func TestBug007_BusinessRegression(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "x.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	_, e = s.CreateDevice(model.Device{Code: "D1", Name: "D1"})
	if e != nil {
		t.Fatal(e)
	}
	_, e = s.DB.Exec("INSERT INTO sensor_data(device_id,sensor_type,value,is_abnormal,abnormal_reason,collected_at,created_at) VALUES(?,?,?,?,?,?,?)", 1, "temperature", 90, 1, "bad", time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))
	if e != nil {
		t.Fatal(e)
	}
	h := service.NewHealthService(repository.NewDeviceRepository(s.DB), repository.NewSensorRepository(s.DB))
	p, e := h.Evaluate(context.Background(), 1)
	if e != nil {
		t.Fatal(e)
	}
	if p.Score >= 60 {
		t.Fatalf("missing dimensions were treated as healthy: %.2f", p.Score)
	}
}
