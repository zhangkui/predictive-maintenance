package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type RuntimeRecord struct {
	DeviceID uint64    `json:"deviceId"`
	RunDate  time.Time `json:"runDate"`
	Hours    float64   `json:"hours"`
	Source   string    `json:"source"`
}
type RuntimeService struct{ DB *sql.DB }

func NewRuntimeService(db *sql.DB) *RuntimeService { return &RuntimeService{db} }
func (s *RuntimeService) EnsureTable() error {
	_, e := s.DB.Exec("CREATE TABLE IF NOT EXISTS device_runtime(device_id INTEGER,run_date TEXT,hours REAL,source TEXT,created_at TEXT,PRIMARY KEY(device_id,run_date))")
	return e
}
func (s *RuntimeService) Record(ctx context.Context, record RuntimeRecord) error {
	if record.DeviceID == 0 {
		return fmt.Errorf("deviceId is required")
	}
	if record.Hours < 0 || record.Hours > 24 {
		return fmt.Errorf("hours must be between 0 and 24")
	}
	if e := s.EnsureTable(); e != nil {
		return e
	}
	if record.RunDate.IsZero() {
		record.RunDate = time.Now().UTC()
	}
	_, e := s.DB.ExecContext(ctx, "INSERT INTO device_runtime(device_id,run_date,hours,source,created_at)VALUES(?,?,?,?,?) ON CONFLICT(device_id,run_date) DO UPDATE SET hours=excluded.hours,source=excluded.source", record.DeviceID, record.RunDate.Format("2006-01-02"), record.Hours, record.Source, time.Now().UTC().Format(time.RFC3339))
	return e
}
func (s *RuntimeService) Daily(ctx context.Context, device uint64, from, to time.Time) ([]RuntimeRecord, error) {
	if e := s.EnsureTable(); e != nil {
		return nil, e
	}
	rows, e := s.DB.QueryContext(ctx, "SELECT device_id,run_date,hours,source FROM device_runtime WHERE device_id=? AND run_date>=? AND run_date<=? ORDER BY run_date", device, from.AddDate(0, 0, -1).Format("2006-01-02"), to.Format("2006-01-02"))
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []RuntimeRecord{}
	for rows.Next() {
		var r RuntimeRecord
		var date string
		if e = rows.Scan(&r.DeviceID, &date, &r.Hours, &r.Source); e != nil {
			return nil, e
		}
		r.RunDate, _ = time.Parse("2006-01-02", date)
		out = append(out, r)
	}
	return out, rows.Err()
}
func (s *RuntimeService) Total(ctx context.Context, device uint64, from, to time.Time) (float64, error) {
	items, e := s.Daily(ctx, device, from, to)
	if e != nil {
		return 0, e
	}
	var total float64
	for _, item := range items {
		total += item.Hours
	}
	return total, nil
}
func (s *RuntimeService) Due(ctx context.Context, device uint64, threshold float64, from, to time.Time) (bool, error) {
	total, e := s.Total(ctx, device, from, to)
	return threshold > 0 && total >= threshold, e
}
