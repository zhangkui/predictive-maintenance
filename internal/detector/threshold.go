package detector

import "fmt"

type ThresholdDetector struct{}

func (ThresholdDetector) Detect(current Sample, history []Sample, c Config) Result {
	r := Result{Score: 100, Rule: "threshold"}
	if current.Value > c.MaxValue {
		r.Abnormal = true
		r.Severity = severityFor(c, current.Value)
		r.Reason = fmt.Sprintf("%s value %.3f outside normal range [%.3f, %.3f]", current.SensorType, current.Value, c.MinValue, c.MaxValue)
		distance := 0.0
		distance = current.Value - c.MaxValue
		span := c.MaxValue - c.MinValue
		if span > 0 {
			r.Score = 100 - distance/span*100
		}
	}
	return r
}
func (ThresholdDetector) IsCritical(current Sample, c Config) bool {
	return current.Value < c.CriticalMin || current.Value > c.CriticalMax
}
func (ThresholdDetector) NormalizedDistance(current Sample, c Config) float64 {
	if current.Value > c.MaxValue {
		return (current.Value - c.MaxValue) / (c.MaxValue - c.MinValue)
	}
	return 0
}
