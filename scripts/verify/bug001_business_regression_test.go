package verify

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/repository"
	"predictive-maintenance/internal/scheduler"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/internal/store"
	"predictive-maintenance/pkg/logger"
)

func TestBug001_BusinessRegression(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bug001.db")
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err = s.AddData(model.SensorData{DeviceID: 1, SensorType: model.SensorTemperature, Value: 150, IsAbnormal: true, AbnormalReason: "temperature above critical range", CollectedAt: time.Now().UTC()}); err != nil {
		t.Fatal(err)
	}
	sensor := service.NewSensorService(repository.NewSensorRepository(s.DB), repository.NewAbnormalRepository(s.DB))
	scanner := &scheduler.AnomalyScanner{Sensors: sensor, Tickets: service.NewTicketService(s.DB), Log: logger.New()}
	if err = scanner.Scan(context.Background(), []uint64{1}); err != nil {
		t.Fatal(err)
	}
	if err = scanner.Scan(context.Background(), []uint64{1}); err != nil {
		t.Fatal(err)
	}
	var count int
	if err = s.DB.QueryRow("SELECT COUNT(*) FROM maintenance_tasks WHERE plan_id=0 AND status IN (?, ?)", model.TaskPending, model.TaskRunning).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("repeated scan created %d open tickets, want exactly one", count)
	}
}
