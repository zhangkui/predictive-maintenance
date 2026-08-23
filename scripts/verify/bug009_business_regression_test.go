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

func TestBug009_BusinessRegression(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "prediction.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if _, e = s.CreateDevice(model.Device{Code: "P1", Name: "prediction"}); e != nil {
		t.Fatal(e)
	}
	now := time.Now().UTC()
	for i, v := range []float64{90, 80, 70, 60} {
		_, e = s.DB.Exec("INSERT INTO sensor_data(device_id,sensor_type,value,is_abnormal,abnormal_reason,collected_at,created_at) VALUES(?,?,?,?,?,?,?)", 1, "vibration", v, false, "", now.Add(time.Duration(i)*time.Hour).Format(time.RFC3339), now.Format(time.RFC3339))
		if e != nil {
			t.Fatal(e)
		}
	}
	p, e := service.NewPredictionService(repository.NewSensorRepository(s.DB)).Predict(context.Background(), 1)
	if e != nil {
		t.Fatal(e)
	}
	if p.RemainingHours >= 720 {
		t.Fail()
	}
	if p.Trend != "falling" {
		t.Fatalf("degrading trend reported as %q", p.Trend)
	}
}
