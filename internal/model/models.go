package model

import "time"

type Device struct {
	ID           uint64    `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Model        string    `json:"model"`
	Manufacturer string    `json:"manufacturer"`
	InstallDate  time.Time `json:"installDate"`
	Location     string    `json:"location"`
	GroupID      uint64    `json:"groupId"`
	Status       int       `json:"status"`
	HealthScore  float64   `json:"healthScore"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
type SensorData struct {
	ID             uint64    `json:"id"`
	DeviceID       uint64    `json:"deviceId"`
	SensorType     string    `json:"sensorType"`
	Value          float64   `json:"value"`
	IsAbnormal     bool      `json:"isAbnormal"`
	AbnormalReason string    `json:"abnormalReason"`
	CollectedAt    time.Time `json:"collectedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}
type AbnormalRecord struct {
	ID            uint64    `json:"id"`
	DeviceID      uint64    `json:"deviceId"`
	SensorType    string    `json:"sensorType"`
	DetectedValue float64   `json:"detectedValue"`
	NormalRange   string    `json:"normalRange"`
	Severity      int       `json:"severity"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
}
type MaintenancePlan struct {
	ID              uint64    `json:"id"`
	DeviceID        uint64    `json:"deviceId"`
	Type            string    `json:"type"`
	CycleType       string    `json:"cycleType"`
	CycleValue      int       `json:"cycleValue"`
	RuntimeHours    int       `json:"runtimeHours"`
	HealthThreshold float64   `json:"healthThreshold"`
	Description     string    `json:"description"`
	IsActive        bool      `json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
}
type MaintenanceTask struct {
	ID            uint64    `json:"id"`
	PlanID        uint64    `json:"planId"`
	DeviceID      uint64    `json:"deviceId"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Priority      int       `json:"priority"`
	Status        int       `json:"status"`
	ScheduledDate time.Time `json:"scheduledDate"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

const (
	DeviceOnline    = 1
	DeviceOffline   = 2
	DeviceFault     = 3
	DeviceRepairing = 4
	AbnormalPending = 1
	AbnormalHandled = 2
	AbnormalIgnored = 3
	TaskPending     = 1
	TaskRunning     = 2
	TaskCompleted   = 3
	TaskCancelled   = 4
)

type DeviceGroup struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
