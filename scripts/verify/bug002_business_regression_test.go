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

func TestBug002_BusinessRegression(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "x.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	ss := service.NewSensorService(repository.NewSensorRepository(s.DB), repository.NewAbnormalRepository(s.DB))
	_, _, e = ss.Ingest(context.Background(), model.SensorDataBatch{DeviceID: 1, Items: []model.SensorReading{{SensorType: model.SensorTemperature, Value: 20, CollectedAt: time.Now()}, {SensorType: "unknown", Value: 30}}})
	if e == nil {
		t.Fatal("invalid item must fail")
	}
	var n int
	if e = s.DB.QueryRow("SELECT COUNT(*) FROM sensor_data").Scan(&n); e != nil {
		t.Fatal(e)
	}
	if n != 0 {
		t.Fatalf("partial batch persisted %d rows", n)
	}
}
