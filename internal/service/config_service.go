package service

import (
	"context"
	"database/sql"
	"fmt"
	"predictive-maintenance/internal/model"
)

type ConfigService struct{ DB *sql.DB }

func NewConfigService(db *sql.DB) *ConfigService { return &ConfigService{db} }
func (s *ConfigService) EnsureTable() error {
	_, e := s.DB.Exec("CREATE TABLE IF NOT EXISTS sensor_configs(id INTEGER PRIMARY KEY AUTOINCREMENT,device_id INTEGER,sensor_type TEXT,unit TEXT,min_value REAL,max_value REAL,critical_min REAL,critical_max REAL,alarm_threshold REAL,sample_interval INTEGER,UNIQUE(device_id,sensor_type))")
	return e
}
func (s *ConfigService) Save(ctx context.Context, c model.SensorConfig) (model.SensorConfig, error) {
	if !model.ValidSensorType(c.SensorType) {
		return c, fmt.Errorf("invalid sensor type")
	}
	if c.MaxValue <= c.MinValue {
		return c, fmt.Errorf("invalid sensor range")
	}
	if e := s.EnsureTable(); e != nil {
		return c, e
	}
	_, e := s.DB.ExecContext(ctx, "INSERT INTO sensor_configs(device_id,sensor_type,unit,min_value,max_value,critical_min,critical_max,alarm_threshold,sample_interval)VALUES(?,?,?,?,?,?,?,?,?) ON CONFLICT(device_id,sensor_type) DO UPDATE SET unit=excluded.unit,min_value=excluded.min_value,max_value=excluded.max_value,critical_min=excluded.critical_min,critical_max=excluded.critical_max,alarm_threshold=excluded.alarm_threshold,sample_interval=excluded.sample_interval", c.DeviceID, c.SensorType, c.Unit, c.MinValue, c.MaxValue, c.CriticalMin, c.CriticalMax, c.AlarmThreshold, c.SampleInterval)
	return c, e
}
func (s *ConfigService) List(ctx context.Context, device uint64) ([]model.SensorConfig, error) {
	if e := s.EnsureTable(); e != nil {
		return nil, e
	}
	rows, e := s.DB.QueryContext(ctx, "SELECT id,device_id,sensor_type,unit,min_value,max_value,critical_min,critical_max,alarm_threshold,sample_interval FROM sensor_configs WHERE device_id=? ORDER BY sensor_type", device)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.SensorConfig{}
	for rows.Next() {
		var c model.SensorConfig
		if e = rows.Scan(&c.ID, &c.DeviceID, &c.SensorType, &c.Unit, &c.MinValue, &c.MaxValue, &c.CriticalMin, &c.CriticalMax, &c.AlarmThreshold, &c.SampleInterval); e != nil {
			return nil, e
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
func (s *ConfigService) Default(device uint64, sensor string) model.SensorConfig {
	return model.SensorConfig{DeviceID: device, SensorType: sensor, Unit: model.SensorUnits[sensor], MinValue: 0, MaxValue: 100, CriticalMin: -10, CriticalMax: 120, AlarmThreshold: 10, SampleInterval: 60}
}
