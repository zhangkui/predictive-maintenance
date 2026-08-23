package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type RetentionService struct{ DB *sql.DB }

func NewRetentionService(db *sql.DB) *RetentionService { return &RetentionService{db} }
func (s *RetentionService) DeleteSensorDataBefore(ctx context.Context, cutoff time.Time) (int64, error) {
	if cutoff.IsZero() {
		return 0, fmt.Errorf("cutoff is required")
	}
	result, e := s.DB.ExecContext(ctx, "DELETE FROM sensor_data WHERE collected_at<?", cutoff.UTC().Format(time.RFC3339))
	if e != nil {
		return 0, e
	}
	return result.RowsAffected()
}
func (s *RetentionService) DeleteHealthHistoryBefore(ctx context.Context, cutoff time.Time) (int64, error) {
	result, e := s.DB.ExecContext(ctx, "DELETE FROM health_history WHERE calculated_at<?", cutoff.UTC().Format(time.RFC3339))
	if e != nil {
		return 0, e
	}
	return result.RowsAffected()
}
func (s *RetentionService) Cleanup(ctx context.Context, days int) (map[string]int64, error) {
	if days < 7 {
		days = 7
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -days)
	sensor, e := s.DeleteSensorDataBefore(ctx, cutoff)
	if e != nil {
		return nil, e
	}
	health, e := s.DeleteHealthHistoryBefore(ctx, cutoff)
	if e != nil {
		return nil, e
	}
	return map[string]int64{"sensorData": sensor, "healthHistory": health}, nil
}
func (s *RetentionService) Preview(ctx context.Context, days int) (map[string]any, error) {
	if days < 7 {
		days = 7
	}
	cutoff := time.Now().UTC().AddDate(0, 0, -days)
	var sensor, health int
	_ = s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM sensor_data WHERE collected_at<?", cutoff.Format(time.RFC3339)).Scan(&sensor)
	_ = s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM health_history WHERE calculated_at<?", cutoff.Format(time.RFC3339)).Scan(&health)
	return map[string]any{"cutoff": cutoff, "sensorData": sensor, "healthHistory": health}, nil
}
