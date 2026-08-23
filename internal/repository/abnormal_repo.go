package repository

import (
	"context"
	"database/sql"
	"fmt"
	"predictive-maintenance/internal/model"
	"time"
)

type AbnormalRepository struct{ DB *sql.DB }

func NewAbnormalRepository(db *sql.DB) *AbnormalRepository { return &AbnormalRepository{db} }
func (r *AbnormalRepository) Create(ctx context.Context, a model.AbnormalRecord) (uint64, error) {
	n := time.Now().UTC().Format(time.RFC3339)
	res, e := r.DB.ExecContext(ctx, "INSERT INTO abnormal_records(device_id,sensor_type,detected_value,normal_range,severity,status,created_at)VALUES(?,?,?,?,?,?,?)", a.DeviceID, a.SensorType, a.DetectedValue, a.NormalRange, a.Severity, model.AbnormalPending, n)
	if e != nil {
		return 0, e
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}
func (r *AbnormalRepository) Handle(ctx context.Context, id, handler uint64, status int) error {
	if status != model.AbnormalHandled && status != 3 {
		return fmt.Errorf("invalid abnormal status")
	}
	res, e := r.DB.ExecContext(ctx, "UPDATE abnormal_records SET status=?,handled_by=?,handled_at=? WHERE id=? AND status=1", status, handler, time.Now().UTC().Format(time.RFC3339), id)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("abnormal not pending")
	}
	return nil
}
func (r *AbnormalRepository) Pending(ctx context.Context, device uint64) ([]model.AbnormalRecord, error) {
	rows, e := r.DB.QueryContext(ctx, "SELECT id,device_id,sensor_type,detected_value,normal_range,severity,status,created_at FROM abnormal_records WHERE status=1 AND device_id=? ORDER BY created_at DESC", device)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.AbnormalRecord{}
	for rows.Next() {
		var a model.AbnormalRecord
		var c string
		if e = rows.Scan(&a.ID, &a.DeviceID, &a.SensorType, &a.DetectedValue, &a.NormalRange, &a.Severity, &a.Status, &c); e != nil {
			return nil, e
		}
		a.CreatedAt = parse(c)
		out = append(out, a)
	}
	return out, rows.Err()
}
