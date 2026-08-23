package repository

import (
	"context"
	"database/sql"
	"predictive-maintenance/internal/model"
	"time"
)

type SensorRepository struct{ DB *sql.DB }

func NewSensorRepository(db *sql.DB) *SensorRepository { return &SensorRepository{db} }
func (r *SensorRepository) Save(ctx context.Context, d model.SensorData) error {
	_, e := r.DB.ExecContext(ctx, "INSERT INTO sensor_data(device_id,sensor_type,value,is_abnormal,abnormal_reason,collected_at,created_at)VALUES(?,?,?,?,?,?,?)", d.DeviceID, d.SensorType, d.Value, d.IsAbnormal, d.AbnormalReason, d.CollectedAt.Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))
	return e
}
func (r *SensorRepository) History(ctx context.Context, device uint64, sensor string, from, to time.Time, limit int) ([]model.SensorData, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	q := "SELECT id,device_id,sensor_type,value,is_abnormal,abnormal_reason,collected_at,created_at FROM sensor_data WHERE device_id=? AND collected_at>=? AND collected_at<=?"
	args := []any{device, from.Format(time.RFC3339), to.Format(time.RFC3339)}
	if sensor != "" {
		q += " AND sensor_type=?"
		args = append(args, sensor)
	}
	q += " ORDER BY collected_at ASC LIMIT ?"
	args = append(args, limit)
	rows, e := r.DB.QueryContext(ctx, q, args...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.SensorData{}
	for rows.Next() {
		var d model.SensorData
		var c, cr string
		if e = rows.Scan(&d.ID, &d.DeviceID, &d.SensorType, &d.Value, &d.IsAbnormal, &d.AbnormalReason, &c, &cr); e != nil {
			return nil, e
		}
		d.CollectedAt = parse(c)
		d.CreatedAt = parse(cr)
		out = append(out, d)
	}
	return out, rows.Err()
}
func (r *SensorRepository) Latest(ctx context.Context, device uint64) ([]model.SensorData, error) {
	return r.History(ctx, device, "", time.Now().Add(-365*24*time.Hour), time.Now().Add(time.Minute), 1000)
}
