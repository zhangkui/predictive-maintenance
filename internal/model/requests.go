package model

import "time"

type SensorConfig struct {
	ID             uint64  `json:"id"`
	DeviceID       uint64  `json:"deviceId"`
	SensorType     string  `json:"sensorType"`
	Unit           string  `json:"unit"`
	MinValue       float64 `json:"minValue"`
	MaxValue       float64 `json:"maxValue"`
	CriticalMin    float64 `json:"criticalMin"`
	CriticalMax    float64 `json:"criticalMax"`
	AlarmThreshold float64 `json:"alarmThreshold"`
	SampleInterval int     `json:"sampleInterval"`
}
type SensorDataBatch struct {
	DeviceID uint64          `json:"deviceId"`
	Items    []SensorReading `json:"items"`
}
type SensorReading struct {
	SensorType  string    `json:"sensorType"`
	Value       float64   `json:"value"`
	CollectedAt time.Time `json:"collectedAt"`
}
type DeviceFilter struct {
	Status  *int
	Type    string
	GroupID *uint64
	Keyword string
}
type HealthPoint struct {
	DeviceID     uint64             `json:"deviceId"`
	Score        float64            `json:"score"`
	Components   map[string]float64 `json:"components"`
	CalculatedAt time.Time          `json:"calculatedAt"`
}
type Prediction struct {
	DeviceID           uint64    `json:"deviceId"`
	RemainingHours     float64   `json:"remainingHours"`
	FailureProbability float64   `json:"failureProbability"`
	Confidence         float64   `json:"confidence"`
	Trend              string    `json:"trend"`
	GeneratedAt        time.Time `json:"generatedAt"`
}
type Ticket struct {
	ID          uint64    `json:"id"`
	DeviceID    uint64    `json:"deviceId"`
	AbnormalID  uint64    `json:"abnormalId"`
	Title       string    `json:"title"`
	Priority    int       `json:"priority"`
	Status      int       `json:"status"`
	AssigneeID  uint64    `json:"assigneeId"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
type Dashboard struct {
	DeviceCount     int     `json:"deviceCount"`
	OnlineCount     int     `json:"onlineCount"`
	FaultCount      int     `json:"faultCount"`
	AverageHealth   float64 `json:"averageHealth"`
	PendingAbnormal int     `json:"pendingAbnormal"`
	PendingTasks    int     `json:"pendingTasks"`
	CriticalToday   int     `json:"criticalToday"`
}
