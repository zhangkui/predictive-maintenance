package scheduler

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"predictive-maintenance/internal/store"
	"predictive-maintenance/pkg/logger"
	"time"
)

type ReportGenerator struct {
	Store *store.Store
	Dir   string
	Log   *logger.Logger
}

func (g *ReportGenerator) Run(ctx context.Context) error {
	stats, e := g.Store.Stats()
	if e != nil {
		return e
	}
	if g.Dir == "" {
		g.Dir = "./reports"
	}
	if e = os.MkdirAll(g.Dir, 0755); e != nil {
		return e
	}
	payload := map[string]any{"generatedAt": time.Now().UTC(), "overview": stats}
	b, e := json.MarshalIndent(payload, "", "  ")
	if e != nil {
		return e
	}
	name := filepath.Join(g.Dir, "report-"+time.Now().UTC().Format("20060102")+".json")
	if e = os.WriteFile(name, b, 0644); e != nil {
		return e
	}
	g.Log.Info("report generated", map[string]any{"file": name})
	return nil
}
