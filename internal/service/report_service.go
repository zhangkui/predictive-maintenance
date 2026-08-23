package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"strconv"
	"time"
)

type ReportService struct{ Dashboard *DashboardService }

func NewReportService(d *DashboardService) *ReportService { return &ReportService{d} }
func (s *ReportService) ExportOverview(ctx context.Context, w io.Writer) error {
	overview, e := s.Dashboard.Overview(ctx)
	if e != nil {
		return e
	}
	out := csv.NewWriter(w)
	if e = out.Write([]string{"metric", "value"}); e != nil {
		return e
	}
	rows := [][]string{{"deviceCount", strconv.Itoa(overview.DeviceCount)}, {"onlineCount", strconv.Itoa(overview.OnlineCount)}, {"faultCount", strconv.Itoa(overview.FaultCount)}, {"averageHealth", fmt.Sprintf("%.2f", overview.AverageHealth)}, {"pendingAbnormal", strconv.Itoa(overview.PendingAbnormal)}, {"pendingTasks", strconv.Itoa(overview.PendingTasks)}, {"criticalToday", strconv.Itoa(overview.CriticalToday)}, {"generatedAt", time.Now().UTC().Format(time.RFC3339)}}
	for _, row := range rows {
		if e = out.Write(row); e != nil {
			return e
		}
	}
	out.Flush()
	return out.Error()
}
func (s *ReportService) HealthSummary(ctx context.Context) map[string]any {
	v, e := s.Dashboard.Overview(ctx)
	if e != nil {
		return map[string]any{"error": e.Error()}
	}
	return map[string]any{"average": v.AverageHealth, "healthy": v.AverageHealth >= 75, "critical": v.AverageHealth < 60}
}
