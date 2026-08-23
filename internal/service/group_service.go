package service

import (
	"context"
	"database/sql"
	"fmt"
	"predictive-maintenance/internal/model"
	"time"
)

type GroupService struct{ DB *sql.DB }

func NewGroupService(db *sql.DB) *GroupService { return &GroupService{db} }
func (s *GroupService) EnsureTable() error {
	_, e := s.DB.Exec("CREATE TABLE IF NOT EXISTS device_groups(id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT NOT NULL UNIQUE,description TEXT,created_at TEXT NOT NULL)")
	return e
}
func (s *GroupService) Create(ctx context.Context, g model.DeviceGroup) (model.DeviceGroup, error) {
	if g.Name == "" {
		return g, fmt.Errorf("group name required")
	}
	now := time.Now().UTC()
	r, e := s.DB.ExecContext(ctx, "INSERT INTO device_groups(name,description,created_at)VALUES(?,?,?)", g.Name, g.Description, now.Format(time.RFC3339))
	if e != nil {
		return g, e
	}
	id, _ := r.LastInsertId()
	g.ID = uint64(id)
	g.CreatedAt = now
	return g, nil
}
func (s *GroupService) List(ctx context.Context) ([]model.DeviceGroup, error) {
	if e := s.EnsureTable(); e != nil {
		return nil, e
	}
	rows, e := s.DB.QueryContext(ctx, "SELECT id,name,description,created_at FROM device_groups ORDER BY name")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.DeviceGroup{}
	for rows.Next() {
		var g model.DeviceGroup
		var c string
		if e = rows.Scan(&g.ID, &g.Name, &g.Description, &c); e != nil {
			return nil, e
		}
		g.CreatedAt, _ = time.Parse(time.RFC3339, c)
		out = append(out, g)
	}
	return out, rows.Err()
}
func (s *GroupService) Assign(ctx context.Context, device, group uint64) error {
	_, e := s.DB.ExecContext(ctx, "UPDATE devices SET group_id=?,updated_at=? WHERE id=?", group, time.Now().UTC().Format(time.RFC3339), device)
	return e
}
