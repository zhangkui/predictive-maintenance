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

func TestBug011_BusinessRegression(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "x.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r := repository.NewMaintenanceRepository(s.DB)
	p := service.NewMaintenanceService(r)
	id, e := p.CreatePlan(context.Background(), model.MaintenancePlan{DeviceID: 1, Type: "periodic", CycleValue: 1, IsActive: true})
	if e != nil {
		t.Fatal(e)
	}
	if n, e := p.GenerateTasks(context.Background()); e != nil || n != 1 {
		t.Fatalf("first generation n=%d err=%v", n, e)
	}
	if n, e := p.GenerateTasks(context.Background()); e != nil || n != 0 {
		t.Fatalf("duplicate cycle generated n=%d err=%v plan=%d", n, e, id)
	}
}
