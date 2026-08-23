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

func TestBug013_BusinessRegression(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "x.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r := service.NewMaintenanceService(repository.NewMaintenanceRepository(s.DB))
	id, e := r.Repo.CreateTask(context.Background(), model.MaintenanceTask{PlanID: 1, DeviceID: 1, Status: model.TaskPending})
	if e != nil {
		t.Fatal(e)
	}
	if e = r.Transition(context.Background(), int(id), model.TaskPending, model.TaskCompleted); e == nil {
		t.Fatal("pending task skipped running state")
	}
}
