package verify

import (
	"predictive-maintenance/internal/detector"
	"testing"
	"time"
)

func TestBug005_BusinessRegression(t *testing.T) {
	now := time.Now()
	h := []detector.Sample{{Value: 0, CollectedAt: now.Add(-time.Minute)}}
	r := detector.TrendDetector{Window: 8}.Detect(detector.Sample{Value: 100, CollectedAt: now}, h, detector.Config{MinValue: 0, MaxValue: 00})
	if r.Abnormal {
		t.Fatal("two-point noise was classified as persistent trend")
	}
}
