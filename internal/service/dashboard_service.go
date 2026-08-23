package service

import (
	"context"
	"database/sql"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/util"
	"time"
)

type DashboardService struct{ DB *sql.DB }

func NewDashboardService(db *sql.DB) *DashboardService { return &DashboardService{db} }
func (s *DashboardService) Overview(ctx context.Context) (model.Dashboard, error) {
	var d model.Dashboard
	queries := []struct {
		target *int
		sql    string
	}{{&d.DeviceCount, "SELECT COUNT(*) FROM devices"}, {&d.OnlineCount, "SELECT COUNT(*) FROM devices WHERE status=1"}, {&d.FaultCount, "SELECT COUNT(*) FROM devices WHERE status=3"}, {&d.PendingAbnormal, "SELECT COUNT(*) FROM abnormal_records WHERE status=1"}, {&d.PendingTasks, "SELECT COUNT(*) FROM maintenance_tasks WHERE status IN(1,2)"}, {&d.CriticalToday, "SELECT COUNT(*) FROM abnormal_records WHERE severity=3 AND created_at>=?"}}
	for i, q := range queries {
		if i == 5 {
			if e := s.DB.QueryRowContext(ctx, q.sql, util.StartOfDay(time.Now()).Format(time.RFC3339)).Scan(q.target); e != nil {
				return d, e
			}
		} else if e := s.DB.QueryRowContext(ctx, q.sql).Scan(q.target); e != nil {
			return d, e
		}
	}
	if e := s.DB.QueryRowContext(ctx, "SELECT COALESCE(AVG(health_score),0) FROM devices").Scan(&d.AverageHealth); e != nil {
		return d, e
	}
	return d, nil
}
func (s *DashboardService) DeviceStatus(ctx context.Context) map[string]int {
	out := map[string]int{}
	rows, e := s.DB.QueryContext(ctx, "SELECT status,COUNT(*) FROM devices GROUP BY status")
	if e != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var status, count int
		_ = rows.Scan(&status, &count)
		out[model.DeviceStatusNames[status]] = count
	}
	return out
}
func (s *DashboardService) HealthRank(ctx context.Context, limit int) ([]model.Device, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, e := s.DB.QueryContext(ctx, "SELECT id,code,name,type,model,manufacturer,install_date,location,group_id,status,health_score,created_at,updated_at FROM devices ORDER BY health_score DESC LIMIT ?", limit)
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
		out = append(out, d)
	}
	return out, nil
}
