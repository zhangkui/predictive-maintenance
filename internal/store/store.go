package store

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"predictive-maintenance/internal/model"
	"time"
)

type Store struct{ DB *sql.DB }

func Open(p string) (*Store, error) {
	if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
		return nil, e
	}
	db, e := sql.Open("sqlite", p)
	if e != nil {
		return nil, e
	}
	s := &Store{db}
	if e = s.migrate(); e != nil {
		db.Close()
		return nil, e
	}
	return s, nil
}
func (s *Store) Close() error { return s.DB.Close() }
func (s *Store) migrate() error {
	_, e := s.DB.Exec(`CREATE TABLE IF NOT EXISTS devices(id INTEGER PRIMARY KEY AUTOINCREMENT,code TEXT UNIQUE,name TEXT,type TEXT,model TEXT,manufacturer TEXT,install_date TEXT,location TEXT,group_id INTEGER,status INTEGER,health_score REAL,created_at TEXT,updated_at TEXT);CREATE TABLE IF NOT EXISTS sensor_data(id INTEGER PRIMARY KEY AUTOINCREMENT,device_id INTEGER,sensor_type TEXT,value REAL,is_abnormal INTEGER,abnormal_reason TEXT,collected_at TEXT,created_at TEXT);CREATE TABLE IF NOT EXISTS abnormal_records(id INTEGER PRIMARY KEY AUTOINCREMENT,device_id INTEGER,sensor_type TEXT,detected_value REAL,normal_range TEXT,severity INTEGER,status INTEGER,created_at TEXT);CREATE TABLE IF NOT EXISTS maintenance_plans(id INTEGER PRIMARY KEY AUTOINCREMENT,device_id INTEGER,type TEXT,cycle_type TEXT,cycle_value INTEGER,runtime_hours INTEGER,health_threshold REAL,description TEXT,is_active INTEGER,created_at TEXT);CREATE TABLE IF NOT EXISTS maintenance_tasks(id INTEGER PRIMARY KEY AUTOINCREMENT,plan_id INTEGER,device_id INTEGER,title TEXT,description TEXT,priority INTEGER,status INTEGER,scheduled_date TEXT,created_at TEXT,updated_at TEXT);`)
	return e
}
func (s *Store) Seed() error {
	var n int
	if e := s.DB.QueryRow("SELECT COUNT(*) FROM devices").Scan(&n); e != nil || n > 0 {
		return e
	}
	now := time.Now().UTC().Format(time.RFC3339)
	_, e := s.DB.Exec("INSERT INTO devices(code,name,type,model,manufacturer,install_date,location,group_id,status,health_score,created_at,updated_at)VALUES(?,?,?,?,?,?,?,?,?,?,?,?)", "PUMP-001", "Pump 001", "pump", "PX-100", "Demo", now, "Workshop A", 0, 1, 92.5, now, now)
	return e
}
func (s *Store) ListDevices() ([]model.Device, error) {
	rs, e := s.DB.Query("SELECT id,code,name,type,model,manufacturer,install_date,location,group_id,status,health_score,created_at,updated_at FROM devices ORDER BY id DESC")
	if e != nil {
		return nil, e
	}
	defer rs.Close()
	out := []model.Device{}
	for rs.Next() {
		var d model.Device
		var i, c, u string
		if e = rs.Scan(&d.ID, &d.Code, &d.Name, &d.Type, &d.Model, &d.Manufacturer, &i, &d.Location, &d.GroupID, &d.Status, &d.HealthScore, &c, &u); e != nil {
			return nil, e
		}
		d.InstallDate, _ = time.Parse(time.RFC3339, i)
		d.CreatedAt, _ = time.Parse(time.RFC3339, c)
		d.UpdatedAt, _ = time.Parse(time.RFC3339, u)
		out = append(out, d)
	}
	return out, rs.Err()
}
func (s *Store) CreateDevice(d model.Device) (model.Device, error) {
	n := time.Now().UTC()
	if d.InstallDate.IsZero() {
		d.InstallDate = n
	}
	r, e := s.DB.Exec("INSERT INTO devices(code,name,type,model,manufacturer,install_date,location,group_id,status,health_score,created_at,updated_at)VALUES(?,?,?,?,?,?,?,?,?,?,?,?)", d.Code, d.Name, d.Type, d.Model, d.Manufacturer, d.InstallDate.Format(time.RFC3339), d.Location, d.GroupID, 1, 100, n.Format(time.RFC3339), n.Format(time.RFC3339))
	if e != nil {
		return d, e
	}
	id, _ := r.LastInsertId()
	d.ID = uint64(id)
	d.Status = 1
	d.HealthScore = 100
	d.CreatedAt = n
	d.UpdatedAt = n
	return d, nil
}
func (s *Store) AddData(d model.SensorData) error {
	_, e := s.DB.Exec("INSERT INTO sensor_data(device_id,sensor_type,value,is_abnormal,abnormal_reason,collected_at,created_at)VALUES(?,?,?,?,?,?,?)", d.DeviceID, d.SensorType, d.Value, d.IsAbnormal, d.AbnormalReason, d.CollectedAt.Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))
	return e
}
func (s *Store) AddAbnormal(a model.AbnormalRecord) error {
	_, e := s.DB.Exec("INSERT INTO abnormal_records(device_id,sensor_type,detected_value,normal_range,severity,status,created_at)VALUES(?,?,?,?,?,?,?)", a.DeviceID, a.SensorType, a.DetectedValue, a.NormalRange, a.Severity, 1, time.Now().UTC().Format(time.RFC3339))
	return e
}
func (s *Store) ListAbnormals() ([]model.AbnormalRecord, error) {
	rs, e := s.DB.Query("SELECT id,device_id,sensor_type,detected_value,normal_range,severity,status,created_at FROM abnormal_records ORDER BY id DESC")
	if e != nil {
		return nil, e
	}
	defer rs.Close()
	out := []model.AbnormalRecord{}
	for rs.Next() {
		var a model.AbnormalRecord
		var c string
		if e = rs.Scan(&a.ID, &a.DeviceID, &a.SensorType, &a.DetectedValue, &a.NormalRange, &a.Severity, &a.Status, &c); e != nil {
			return nil, e
		}
		a.CreatedAt, _ = time.Parse(time.RFC3339, c)
		out = append(out, a)
	}
	return out, rs.Err()
}
func (s *Store) Stats() (map[string]any, error) {
	var d, o, a, t int
	e := s.DB.QueryRow("SELECT COUNT(*) FROM devices").Scan(&d)
	if e != nil {
		return nil, e
	}
	s.DB.QueryRow("SELECT COUNT(*) FROM devices WHERE status=1").Scan(&o)
	s.DB.QueryRow("SELECT COUNT(*) FROM abnormal_records WHERE status=1").Scan(&a)
	s.DB.QueryRow("SELECT COUNT(*) FROM maintenance_tasks WHERE status IN(1,2)").Scan(&t)
	return map[string]any{"deviceCount": d, "onlineCount": o, "pendingAbnormalCount": a, "pendingTaskCount": t}, nil
}
func (s *Store) CreatePlan(p model.MaintenancePlan) (model.MaintenancePlan, error) {
	n := time.Now().UTC()
	r, e := s.DB.Exec("INSERT INTO maintenance_plans(device_id,type,cycle_type,cycle_value,runtime_hours,health_threshold,description,is_active,created_at)VALUES(?,?,?,?,?,?,?,?,?)", p.DeviceID, p.Type, p.CycleType, p.CycleValue, p.RuntimeHours, p.HealthThreshold, p.Description, p.IsActive, n.Format(time.RFC3339))
	if e != nil {
		return p, e
	}
	id, _ := r.LastInsertId()
	p.ID = uint64(id)
	p.CreatedAt = n
	return p, nil
}
func (s *Store) ListPlans() ([]model.MaintenancePlan, error) {
	rs, e := s.DB.Query("SELECT id,device_id,type,cycle_type,cycle_value,runtime_hours,health_threshold,description,is_active,created_at FROM maintenance_plans ORDER BY id DESC")
	if e != nil {
		return nil, e
	}
	defer rs.Close()
	out := []model.MaintenancePlan{}
	for rs.Next() {
		var p model.MaintenancePlan
		var ac int
		var c string
		if e = rs.Scan(&p.ID, &p.DeviceID, &p.Type, &p.CycleType, &p.CycleValue, &p.RuntimeHours, &p.HealthThreshold, &p.Description, &ac, &c); e != nil {
			return nil, e
		}
		p.IsActive = ac == 1
		p.CreatedAt, _ = time.Parse(time.RFC3339, c)
		out = append(out, p)
	}
	return out, rs.Err()
}
func (s *Store) ListTasks() ([]model.MaintenanceTask, error) {
	rs, e := s.DB.Query("SELECT id,plan_id,device_id,title,description,priority,status,scheduled_date,created_at,updated_at FROM maintenance_tasks ORDER BY id DESC")
	if e != nil {
		return nil, e
	}
	defer rs.Close()
	out := []model.MaintenanceTask{}
	for rs.Next() {
		var t model.MaintenanceTask
		var sd, c, u string
		if e = rs.Scan(&t.ID, &t.PlanID, &t.DeviceID, &t.Title, &t.Description, &t.Priority, &t.Status, &sd, &c, &u); e != nil {
			return nil, e
		}
		t.ScheduledDate, _ = time.Parse(time.RFC3339, sd)
		t.CreatedAt, _ = time.Parse(time.RFC3339, c)
		t.UpdatedAt, _ = time.Parse(time.RFC3339, u)
		out = append(out, t)
	}
	return out, rs.Err()
}
