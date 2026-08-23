package repository

import (
	"context"
	"database/sql"
	"fmt"
	"predictive-maintenance/internal/model"
	"time"
)

type DeviceRepository struct{ DB *sql.DB }

func NewDeviceRepository(db *sql.DB) *DeviceRepository { return &DeviceRepository{db} }
func (r *DeviceRepository) Find(ctx context.Context, id uint64) (model.Device, error) {
	var d model.Device
	var install, created, updated string
	e := r.DB.QueryRowContext(ctx, "SELECT id,code,name,type,model,manufacturer,install_date,location,group_id,status,health_score,created_at,updated_at FROM devices WHERE id=?", id).Scan(&d.ID, &d.Code, &d.Name, &d.Type, &d.Model, &d.Manufacturer, &install, &d.Location, &d.GroupID, &d.Status, &d.HealthScore, &created, &updated)
	if e != nil {
		return d, e
	}
	d.InstallDate = parse(install)
	d.CreatedAt = parse(created)
	d.UpdatedAt = parse(updated)
	return d, nil
}
func (r *DeviceRepository) UpdateStatus(ctx context.Context, id uint64, status int) error {
	if !model.ValidDeviceStatus(status) {
		return fmt.Errorf("invalid device status %d", status)
	}
	_, e := r.DB.ExecContext(ctx, "UPDATE devices SET status=?,updated_at=CURRENT_TIMESTAMP WHERE id=?", status, id)
	return e
}
func (r *DeviceRepository) UpdateHealth(ctx context.Context, id uint64, score float64) error {
	_, e := r.DB.ExecContext(ctx, "UPDATE devices SET health_score=?,updated_at=CURRENT_TIMESTAMP WHERE id=?", score, id)
	return e
}
func (r *DeviceRepository) ListByStatus(ctx context.Context, status int) ([]model.Device, error) {
	rows, e := r.DB.QueryContext(ctx, "SELECT id,code,name,type,model,manufacturer,install_date,location,group_id,status,health_score,created_at,updated_at FROM devices WHERE status=? ORDER BY id", status)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.Device{}
	for rows.Next() {
		var d model.Device
		var i, c, u string
		if e = rows.Scan(&d.ID, &d.Code, &d.Name, &d.Type, &d.Model, &d.Manufacturer, &i, &d.Location, &d.GroupID, &d.Status, &d.HealthScore, &c, &u); e != nil {
			return nil, e
		}
		d.InstallDate = parse(i)
		d.CreatedAt = parse(c)
		d.UpdatedAt = parse(u)
		out = append(out, d)
	}
	return out, rows.Err()
}
func (r *DeviceRepository) ListAll(ctx context.Context) ([]model.Device, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id,code,name,type,model,manufacturer,install_date,location,group_id,status,health_score,created_at,updated_at FROM devices ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	devices := []model.Device{}
	for rows.Next() {
		var d model.Device
		var install, created, updated string
		if err = rows.Scan(&d.ID, &d.Code, &d.Name, &d.Type, &d.Model, &d.Manufacturer, &install, &d.Location, &d.GroupID, &d.Status, &d.HealthScore, &created, &updated); err != nil {
			return nil, err
		}
		d.InstallDate, d.CreatedAt, d.UpdatedAt = parse(install), parse(created), parse(updated)
		devices = append(devices, d)
	}
	return devices, rows.Err()
}
func parse(v string) time.Time { t, _ := time.Parse(time.RFC3339, v); return t }
