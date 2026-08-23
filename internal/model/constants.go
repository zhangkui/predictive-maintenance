package model

const (
	SensorTemperature = "temperature"
	SensorVibration   = "vibration"
	SensorPressure    = "pressure"
	SensorCurrent     = "current"
	SensorSpeed       = "speed"
	SeverityWarning   = 1
	SeveritySerious   = 2
	SeverityCritical  = 3
	PriorityLow       = 1
	PriorityMedium    = 2
	PriorityHigh      = 3
	PriorityEmergency = 4
)

var SensorUnits = map[string]string{SensorTemperature: "degC", SensorVibration: "mm/s", SensorPressure: "MPa", SensorCurrent: "A", SensorSpeed: "RPM"}
var DeviceStatusNames = map[int]string{DeviceOnline: "online", DeviceOffline: "offline", DeviceFault: "fault", DeviceRepairing: "repairing"}
var TaskStatusNames = map[int]string{TaskPending: "pending", TaskRunning: "running", 3: "completed", 4: "cancelled"}

func ValidSensorType(v string) bool { _, ok := SensorUnits[v]; return ok }
func ValidDeviceStatus(v int) bool  { _, ok := DeviceStatusNames[v]; return ok }
func ValidTaskStatus(v int) bool    { _, ok := TaskStatusNames[v]; return ok }
