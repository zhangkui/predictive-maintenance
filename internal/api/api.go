package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"predictive-maintenance/internal/model"
	"predictive-maintenance/internal/repository"
	"predictive-maintenance/internal/service"
	"predictive-maintenance/internal/store"
	"strconv"
	"strings"
)

type API struct {
	store         *store.Store
	devices       *service.DeviceService
	sensors       *service.SensorService
	health        *service.HealthService
	anomalies     *service.AnomalyService
	prediction    *service.PredictionService
	maintenance   *service.MaintenanceService
	tickets       *service.TicketService
	configs       *service.ConfigService
	healthHistory *service.HealthHistoryService
	groups        *service.GroupService
	runtime       *service.RuntimeService
	retention     *service.RetentionService
}

func New(s *store.Store) http.Handler {
	dr := repository.NewDeviceRepository(s.DB)
	sr := repository.NewSensorRepository(s.DB)
	ar := repository.NewAbnormalRepository(s.DB)
	mr := repository.NewMaintenanceRepository(s.DB)
	configs := service.NewConfigService(s.DB)
	healthHistory := service.NewHealthHistoryService(s.DB)
	health := service.NewHealthService(dr, sr)
	health.UseHistory(healthHistory)
	sensors := service.NewSensorService(sr, ar)
	sensors.UseConfigs(configs)
	a := &API{store: s, devices: service.NewDeviceService(dr), sensors: sensors, health: health, anomalies: service.NewAnomalyService(ar), prediction: service.NewPredictionService(sr), maintenance: service.NewMaintenanceService(mr), tickets: service.NewTicketService(s.DB), configs: configs, healthHistory: healthHistory, groups: service.NewGroupService(s.DB), runtime: service.NewRuntimeService(s.DB), retention: service.NewRetentionService(s.DB)}
	m := http.NewServeMux()
	m.HandleFunc("/healthz", a.healthz)
	m.HandleFunc("/api/v1/devices", a.devicesHandler)
	m.HandleFunc("/api/v1/data", a.dataHandler)
	m.HandleFunc("/api/v1/anomalies", a.anomaliesHandler)
	m.HandleFunc("/api/v1/anomalies/handle", a.handleAnomaly)
	m.HandleFunc("/api/v1/maintenance/plans", a.plansHandler)
	m.HandleFunc("/api/v1/maintenance/tasks", a.tasksHandler)
	m.HandleFunc("/api/v1/maintenance/generate", a.generateMaintenance)
	m.HandleFunc("/api/v1/dashboard/overview", a.overview)
	m.HandleFunc("/api/v1/health/", a.healthHandler)
	m.HandleFunc("/api/v1/predictions/", a.predictionsHandler)
	m.HandleFunc("/api/v1/tickets", a.ticketsHandler)
	m.HandleFunc("/api/v1/sensor-configs", a.sensorConfigsHandler)
	m.HandleFunc("/api/v1/health-history", a.healthHistoryHandler)
	m.HandleFunc("/api/v1/groups", a.groupsHandler)
	m.HandleFunc("/api/v1/runtime", a.runtimeHandler)
	m.HandleFunc("/api/v1/retention", a.retentionHandler)
	return cors(m)
}
func cors(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		n.ServeHTTP(w, r)
	})
}
func write(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": v})
}
func fail(w http.ResponseWriter, status int, e error) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": e.Error()})
}
func decode(r *http.Request, v any) error { return json.NewDecoder(r.Body).Decode(v) }
func pathID(path, prefix string) (uint64, error) {
	v := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if i := strings.Index(v, "/"); i >= 0 {
		v = v[:i]
	}
	return strconv.ParseUint(v, 10, 64)
}
func (a *API) healthz(w http.ResponseWriter, r *http.Request) {
	write(w, 200, map[string]string{"status": "ok", "service": "predictive-maintenance"})
}
func (a *API) devicesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		v, e := a.store.ListDevices()
		if e != nil {
			fail(w, 500, e)
			return
		}
		write(w, 200, v)
	case "POST":
		var d model.Device
		if e := decode(r, &d); e != nil {
			fail(w, 400, e)
			return
		}
		v, e := a.devices.Register(r.Context(), d)
		if e != nil {
			fail(w, 400, e)
			return
		}
		write(w, 201, v)
	default:
		fail(w, 405, errors.New("method not allowed"))
	}
}
func (a *API) dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fail(w, 405, errors.New("method not allowed"))
		return
	}
	var b model.SensorDataBatch
	if e := decode(r, &b); e != nil {
		fail(w, 400, e)
		return
	}
	saved, abnormal, e := a.sensors.Ingest(r.Context(), b)
	if e != nil {
		fail(w, 400, e)
		return
	}
	write(w, 202, map[string]any{"accepted": len(saved), "abnormal": abnormal})
}
func (a *API) anomaliesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fail(w, 405, errors.New("method not allowed"))
		return
	}
	device := uint64(0)
	if v := r.URL.Query().Get("deviceId"); v != "" {
		device, _ = strconv.ParseUint(v, 10, 64)
	}
	if device == 0 {
		v, e := a.store.ListAbnormals()
		if e != nil {
			fail(w, 500, e)
			return
		}
		write(w, 200, v)
		return
	}
	v, e := a.anomalies.ListPending(r.Context(), device)
	if e != nil {
		fail(w, 500, e)
		return
	}
	write(w, 200, v)
}
func (a *API) handleAnomaly(w http.ResponseWriter, r *http.Request) {
	var q struct {
		ID        uint64 `json:"id"`
		HandlerID uint64 `json:"handlerId"`
		Ignore    bool   `json:"ignore"`
	}
	if e := decode(r, &q); e != nil {
		fail(w, 400, e)
		return
	}
	var e error
	if q.Ignore {
		e = a.anomalies.Ignore(r.Context(), q.ID, q.HandlerID)
	} else {
		e = a.anomalies.Handle(r.Context(), q.ID, q.HandlerID)
	}
	if e != nil {
		fail(w, 400, e)
		return
	}
	write(w, 200, map[string]any{"updated": true})
}
func (a *API) plansHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		v, e := a.store.ListPlans()
		if e != nil {
			fail(w, 500, e)
			return
		}
		write(w, 200, v)
		return
	}
	var p model.MaintenancePlan
	if e := decode(r, &p); e != nil {
		fail(w, 400, e)
		return
	}
	id, e := a.maintenance.CreatePlan(r.Context(), p)
	if e != nil {
		fail(w, 400, e)
		return
	}
	write(w, 201, map[string]uint64{"id": id})
}
func (a *API) tasksHandler(w http.ResponseWriter, r *http.Request) {
	v, e := a.store.ListTasks()
	if e != nil {
		fail(w, 500, e)
		return
	}
	write(w, 200, v)
}
func (a *API) generateMaintenance(w http.ResponseWriter, r *http.Request) {
	n, e := a.maintenance.GenerateTasks(r.Context())
	if e != nil {
		fail(w, 500, e)
		return
	}
	write(w, 200, map[string]int{"created": n})
}
func (a *API) overview(w http.ResponseWriter, r *http.Request) {
	v, e := a.store.Stats()
	if e != nil {
		fail(w, 500, e)
		return
	}
	write(w, 200, v)
}
func (a *API) healthHandler(w http.ResponseWriter, r *http.Request) {
	id, e := pathID(r.URL.Path, "/api/v1/health/")
	if e != nil {
		fail(w, 400, e)
		return
	}
	v, e := a.health.Evaluate(r.Context(), id)
	if e != nil {
		fail(w, 500, e)
		return
	}
	write(w, 200, v)
}
func (a *API) predictionsHandler(w http.ResponseWriter, r *http.Request) {
	id, e := pathID(r.URL.Path, "/api/v1/predictions/")
	if e != nil {
		fail(w, 400, e)
		return
	}
	v, e := a.prediction.Predict(r.Context(), id)
	if e != nil {
		fail(w, 500, e)
		return
	}
	write(w, 200, v)
}
func (a *API) ticketsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		device, _ := strconv.ParseUint(r.URL.Query().Get("deviceId"), 10, 64)
		v, e := a.tickets.List(r.Context(), device)
		if e != nil {
			fail(w, 500, e)
			return
		}
		write(w, 200, v)
		return
	}
	var x model.AbnormalRecord
	if e := decode(r, &x); e != nil {
		fail(w, 400, e)
		return
	}
	v, e := a.tickets.Create(r.Context(), x)
	if e != nil {
		fail(w, 400, e)
		return
	}
	write(w, 201, v)
}
