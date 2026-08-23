package verify

import (
	"context"
	"path/filepath"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/repository"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/internal/store"
	"testing"
)

func TestBug008_BusinessRegression(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "health.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if _, e = s.CreateDevice(model.Device{Code: "H1", Name: "health"}); e != nil {
		t.Fatal(e)
	}
	h := service.NewHealthService(repository.NewDeviceRepository(s.DB), repository.NewSensorRepository(s.DB))
	a := repository.NewAbnormalRepository(s.DB)
	h.UseAbnormals(a)
	if _, e = h.Evaluate(context.Background(), 1); e != nil {
		t.Fatal(e)
	}
	items, e := a.PendingHealth(context.Background(), 1)
	if e != nil {
		t.Fatal(e)
	}
	if len(items) != 1 {
		t.Fatalf("low health evaluation created %d alerts", len(items))
	}
}
