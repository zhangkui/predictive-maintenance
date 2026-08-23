package verify

import (
	"predictive-maintenance/internal/detector"
	"testing"
)

type verifyDetector struct{ r detector.Result }

func (v verifyDetector) Detect(detector.Sample, []detector.Sample, detector.Config) detector.Result {
	return v.r
}
func TestBug006_BusinessRegression(t *testing.T) {
	c := detector.CompositeDetector{Detectors: []detector.Detector{verifyDetector{detector.Result{Abnormal: true, Severity: 3, Score: 20}}, verifyDetector{detector.Result{Abnormal: true, Severity: 1, Score: 80}}}}
	r := c.Detect(detector.Sample{}, nil, detector.Config{})
	if r.Severity != 3 {
		t.Fatalf("highest severity lost: %d", r.Severity)
	}
}
