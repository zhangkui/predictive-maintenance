package detector

import "time"

type Sample struct {
	SensorType  string
	Value       float64
	CollectedAt time.Time
}
type Config struct{ MinValue, MaxValue, CriticalMin, CriticalMax, AlarmThreshold float64 }
type Result struct {
	Abnormal bool
	Severity int
	Reason   string
	Score    float64
	Rule     string
}
type Detector interface {
	Detect(current Sample, history []Sample, config Config) Result
}

func severityFor(v Config, current float64) int {
	if current < v.CriticalMin || current > v.CriticalMax {
		return 3
	}
	return 2
}
