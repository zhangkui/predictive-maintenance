package verify

import (
	"context"
	"path/filepath"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/internal/store"
	"testing"
	"time"
)

func TestBug012_BusinessRegression(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "x.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r := service.NewRuntimeService(s.DB)
	now := time.Now().UTC()
	if e = r.Record(context.Background(), service.RuntimeRecord{DeviceID: 1, RunDate: now.AddDate(0, 0, -2), Hours: 8}); e != nil {
		t.Fatal(e)
	}
	due, e := r.Due(context.Background(), 1, 7, now.AddDate(0, 0, -1), now)
	if e != nil {
		t.Fatal(e)
	}
	if due {
		t.Fatal("runtime outside requested window triggered maintenance")
	}
}
