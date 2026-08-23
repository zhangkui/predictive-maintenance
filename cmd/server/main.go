package main

import (
	"context"
	"log"
	"net/http"
	"predictive-maintenance/internal/api"
	"predictive-maintenance/internal/config"
	"predictive-maintenance/internal/repository"
	"predictive-maintenance/internal/scheduler"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/internal/store"
	"predictive-maintenance/pkg/logger"
	"time"
)

func main() {
	c := config.Load()
	if e := c.Validate(); e != nil {
		log.Fatal(e)
	}
	db, e := store.Open(c.DBPath)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	if e = db.Seed(); e != nil {
		log.Fatal(e)
	}
	appLog := logger.New()
	dr := repository.NewDeviceRepository(db.DB)
	sr := repository.NewSensorRepository(db.DB)
	mr := repository.NewMaintenanceRepository(db.DB)
	devices := service.NewDeviceService(dr)
	health := service.NewHealthService(dr, sr)
	maintenance := service.NewMaintenanceService(mr)
	runner := scheduler.Build(context.Background(), c, db, devices, health, maintenance, appLog)
	if c.EnableScheduler {
		runner.Start(context.Background())
		defer runner.Stop()
	}
	handler := api.New(db)
	server := &http.Server{Addr: c.HTTPAddr, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	appLog.Info("server started", map[string]any{"address": c.HTTPAddr, "schedulerJobs": runner.JobCount()})
	log.Fatal(server.ListenAndServe())
}
