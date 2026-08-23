package service

import (
	"context"
	"database/sql"
	"fmt"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/util"
	"time"
)

type TicketService struct{ DB *sql.DB }

func NewTicketService(db *sql.DB) *TicketService { return &TicketService{db} }
func (s *TicketService) Create(ctx context.Context, a model.AbnormalRecord) (model.Ticket, error) {
	if a.ID == 0 || a.DeviceID == 0 {
		return model.Ticket{}, fmt.Errorf("abnormal relation is required")
	}
	now := time.Now().UTC()
	title := "sensor abnormal: " + a.SensorType
	res, e := s.DB.ExecContext(ctx, "INSERT INTO maintenance_tasks(plan_id,device_id,title,description,priority,status,scheduled_date,created_at,updated_at)VALUES(0,?,?,?,?,?,?,?,?)", a.DeviceID, title, a.NormalRange, a.Severity, model.TaskPending, now.Format(time.RFC3339), now.Format(time.RFC3339), now.Format(time.RFC3339))
	if e != nil {
		return model.Ticket{}, e
	}
	id, _ := res.LastInsertId()
	return model.Ticket{ID: uint64(id), DeviceID: a.DeviceID, AbnormalID: a.ID, Title: title, Priority: a.Severity, Status: model.TaskPending, Description: a.NormalRange, CreatedAt: now}, nil
}
func (s *TicketService) Code(id uint64) string { return util.TicketCode(id) }
func (s *TicketService) List(ctx context.Context, device uint64) ([]model.Ticket, error) {
	q := "SELECT id,device_id,title,priority,status,description,created_at FROM maintenance_tasks WHERE plan_id=0"
	args := []any{}
	if device > 0 {
		q += " AND device_id=?"
		args = append(args, device)
	}
	q += " ORDER BY id DESC"
	rows, e := s.DB.QueryContext(ctx, q, args...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.Ticket{}
	for rows.Next() {
		var t model.Ticket
		var c string
		if e = rows.Scan(&t.ID, &t.DeviceID, &t.Title, &t.Priority, &t.Status, &t.Description, &c); e != nil {
			return nil, e
		}
		t.CreatedAt, _ = time.Parse(time.RFC3339, c)
		out = append(out, t)
	}
	return out, rows.Err()
}
