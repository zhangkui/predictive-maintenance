package service

import (
	"context"
	"database/sql"
	"predictive-maintenance/internal/model"
	"time"
)

type HealthHistoryService struct{ DB *sql.DB }

func NewHealthHistoryService(db *sql.DB) *HealthHistoryService { return &HealthHistoryService{db} }
func (s *HealthHistoryService) EnsureTable() error {
	_, e := s.DB.Exec("CREATE TABLE IF NOT EXISTS health_history(id INTEGER PRIMARY KEY AUTOINCREMENT,device_id INTEGER,score REAL,components TEXT,calculated_at TEXT)")
	return e
}
func (s *HealthHistoryService) Save(ctx context.Context, p model.HealthPoint) error {
	if e := s.EnsureTable(); e != nil {
		return e
	}
	_, e := s.DB.ExecContext(ctx, "INSERT INTO health_history(device_id,score,components,calculated_at)VALUES(?,?,?,?)", p.DeviceID, p.Score, "{}", p.CalculatedAt.Format(time.RFC3339))
	return e
}
func (s *HealthHistoryService) List(ctx context.Context, device uint64, from, to time.Time, limit int) ([]model.HealthPoint, error) {
	if e := s.EnsureTable(); e != nil {
		return nil, e
	}
	if limit <= 0 {
		limit = 200
	}
	rows, e := s.DB.QueryContext(ctx, "SELECT device_id,score,calculated_at FROM health_history WHERE device_id=? AND calculated_at>=? AND calculated_at<=? ORDER BY calculated_at DESC LIMIT ?", device, from.Format(time.RFC3339), to.Format(time.RFC3339), limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.HealthPoint{}
	for rows.Next() {
		var p model.HealthPoint
		var c string
		if e = rows.Scan(&p.DeviceID, &p.Score, &c); e != nil {
			return nil, e
		}
		p.CalculatedAt, _ = time.Parse(time.RFC3339, c)
		out = append(out, p)
	}
	return out, rows.Err()
}
func (s *HealthHistoryService) Latest(ctx context.Context, device uint64) (model.HealthPoint, error) {
	items, e := s.List(ctx, device, time.Now().Add(-365*24*time.Hour), time.Now().Add(time.Hour), 1)
	if e != nil {
		return model.HealthPoint{}, e
	}
	if len(items) == 0 {
		return model.HealthPoint{DeviceID: device}, nil
	}
	return items[len(items)-1], nil
}
