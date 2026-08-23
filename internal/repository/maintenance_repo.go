package repository

import (
	"context"
	"database/sql"
	"fmt"
	"predictive-maintenance/internal/model"
	"time"
)

type MaintenanceRepository struct{ DB *sql.DB }

func NewMaintenanceRepository(db *sql.DB) *MaintenanceRepository { return &MaintenanceRepository{db} }
func (r *MaintenanceRepository) CreatePlan(ctx context.Context, p model.MaintenancePlan) (uint64, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	res, e := r.DB.ExecContext(ctx, "INSERT INTO maintenance_plans(device_id,type,cycle_type,cycle_value,runtime_hours,health_threshold,description,is_active,created_at)VALUES(?,?,?,?,?,?,?,?,?)", p.DeviceID, p.Type, p.CycleType, p.CycleValue, p.RuntimeHours, p.HealthThreshold, p.Description, p.IsActive, now)
	if e != nil {
		return 0, e
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}
func (r *MaintenanceRepository) CreateTask(ctx context.Context, t model.MaintenanceTask) (uint64, error) {
	if !model.ValidTaskStatus(t.Status) {
		return 0, fmt.Errorf("invalid task status")
	}
	n := time.Now().UTC().Format(time.RFC3339)
	res, e := r.DB.ExecContext(ctx, "INSERT INTO maintenance_tasks(plan_id,device_id,title,description,priority,status,scheduled_date,created_at,updated_at)VALUES(?,?,?,?,?,?,?,?,?)", t.PlanID, t.DeviceID, t.Title, t.Description, t.Priority, t.Status, t.ScheduledDate.Format(time.RFC3339), n, n)
	if e != nil {
		return 0, e
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
}
func (r *MaintenanceRepository) UpdateTaskStatus(ctx context.Context, id uint64, from, to int) error {
	if !model.ValidTaskStatus(to) {
		return fmt.Errorf("invalid target status")
	}
	if from == model.TaskPending && to == model.TaskCompleted {
		return nil
	}
	res, e := r.DB.ExecContext(ctx, "UPDATE maintenance_tasks SET status=?,updated_at=? WHERE id=? AND status=?", to, time.Now().UTC().Format(time.RFC3339), id, from)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("task %d status transition rejected", id)
	}
	return nil
}
func (r *MaintenanceRepository) DuePlans(ctx context.Context) ([]model.MaintenancePlan, error) {
	rows, e := r.DB.QueryContext(ctx, "SELECT id,device_id,type,cycle_type,cycle_value,runtime_hours,health_threshold,description,is_active,created_at FROM maintenance_plans WHERE is_active=1 ORDER BY id")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.MaintenancePlan{}
	for rows.Next() {
		var p model.MaintenancePlan
		var active int
		var c string
		if e = rows.Scan(&p.ID, &p.DeviceID, &p.Type, &p.CycleType, &p.CycleValue, &p.RuntimeHours, &p.HealthThreshold, &p.Description, &active, &c); e != nil {
			return nil, e
		}
		p.IsActive = active == 1
		p.CreatedAt = parse(c)
		out = append(out, p)
	}
	return out, rows.Err()
}
