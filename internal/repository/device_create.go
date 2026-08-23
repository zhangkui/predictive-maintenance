package repository

import (
	"context"
	"predictive-maintenance/internal/model"
	"time"
)

func (r *DeviceRepository) DBCreate(ctx context.Context, d model.Device) (model.Device, error) {
	now := time.Now().UTC()
	res, e := r.DB.ExecContext(ctx, "INSERT INTO devices(code,name,type,model,manufacturer,install_date,location,group_id,status,health_score,created_at,updated_at)VALUES(?,?,?,?,?,?,?,?,?,?,?,?)", d.Code, d.Name, d.Type, d.Model, d.Manufacturer, d.InstallDate.Format(time.RFC3339), d.Location, d.GroupID, model.DeviceOnline, 100, now.Format(time.RFC3339), now.Format(time.RFC3339))
	if e != nil {
		return d, e
	}
	id, _ := res.LastInsertId()
	d.ID = uint64(id)
	d.Status = model.DeviceOnline
	d.HealthScore = 100
	d.CreatedAt = now
	d.UpdatedAt = now
	return d, nil
}
