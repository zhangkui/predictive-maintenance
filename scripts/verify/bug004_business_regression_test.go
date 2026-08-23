package verify

import (
	"predictive-maintenance/internal/detector"
	"testing"
	"time"
)

func TestBug004_BusinessRegression(t *testing.T) {
	now := time.Now()
	r := (detector.RateDetector{}).Rate(detector.Sample{Value: 10, CollectedAt: now}, detector.Sample{Value: 1, CollectedAt: now})
	if r != 0 {
		t.Fatalf("invalid interval produced rate %v", r)
	}
}
