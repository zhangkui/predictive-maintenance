package verify

import (
	"predictive-maintenance/internal/detector"
	"testing"
)

func TestBug003_BusinessRegression(t *testing.T) {
	c := detector.Config{MinValue: 10, MaxValue: 50, CriticalMin: 0, CriticalMax: 60}
	r := (detector.ThresholdDetector{}).Detect(detector.Sample{SensorType: "temperature", Value: 5}, nil, c)
	if !r.Abnormal {
		t.Fatal("below-normal value was accepted")
	}
	if (detector.ThresholdDetector{}).NormalizedDistance(detector.Sample{Value: 5}, c) <= 0 {
		t.Fatal("lower bound distance was lost")
	}
}
