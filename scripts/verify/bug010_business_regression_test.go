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

func TestBug010_BusinessRegression(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "confidence.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if _, e = s.CreateDevice(model.Device{Code: "C1", Name: "confidence"}); e != nil {
		t.Fatal(e)
	}
	now := time.Now().UTC()
	now = now.Add(-24 * time.Hour)
	for i := 0; i < 20; i++ {
		_, e = s.DB.Exec("INSERT INTO sensor_data(device_id,sensor_type,value,is_abnormal,abnormal_reason,collected_at,created_at) VALUES(?,?,?,?,?,?,?)", 1, "vibration", float64(i), false, "", now.Add(time.Duration(i)*time.Hour).Format(time.RFC3339), now.Format(time.RFC3339))
		if e != nil {
			t.Fatal(e)
		}
	}
	p, e := service.NewPredictionService(repository.NewSensorRepository(s.DB)).Predict(context.Background(), 1)
	if e != nil {
		t.Fatal(e)
	}
	if p.Confidence >= .5 {
		t.Fatalf("insufficient prediction confidence was %.2f", p.Confidence)
	}
}
