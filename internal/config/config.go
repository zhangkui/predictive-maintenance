package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr            string
	DBPath              string
	OfflineTimeout      time.Duration
	SchedulerInterval   time.Duration
	HealthInterval      time.Duration
	MaintenanceInterval time.Duration
	EnableScheduler     bool
}

func Load() Config {
	c := Config{HTTPAddr: env("HTTP_ADDR", ":8080"), DBPath: env("DB_PATH", "./data/predictive-maintenance.db"), OfflineTimeout: duration("OFFLINE_TIMEOUT", 300), SchedulerInterval: duration("SCHEDULER_INTERVAL", 60), HealthInterval: duration("HEALTH_INTERVAL", 300), MaintenanceInterval: duration("MAINTENANCE_INTERVAL", 3600), EnableScheduler: boolenv("ENABLE_SCHEDULER", true)}
	return c
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func duration(k string, seconds int) time.Duration {
	v := env(k, "")
	if v == "" {
		return time.Duration(seconds) * time.Second
	}
	n, e := strconv.Atoi(v)
	if e != nil || n <= 0 {
		return time.Duration(seconds) * time.Second
	}
	return time.Duration(n) * time.Second
}
func boolenv(k string, d bool) bool {
	v := env(k, "")
	if v == "" {
		return d
	}
	b, e := strconv.ParseBool(v)
	if e != nil {
		return d
	}
	return b
}
func (c Config) Validate() error {
	if c.HTTPAddr == "" || c.DBPath == "" {
		return os.ErrInvalid
	}
	return nil
}
