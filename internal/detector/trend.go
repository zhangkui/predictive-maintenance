package detector

import (
	"fmt"
	"predictive-maintenance/internal/util"
)

type TrendDetector struct {
	Window     int
	SlopeLimit float64
}

func (t TrendDetector) Detect(current Sample, history []Sample, c Config) Result {
	out := Result{Score: 100, Rule: "trend"}
	points := Recent(append(history, current), t.Window)
	if len(points) < 2 {
		return out
	}
	xs := make([]float64, len(points))
	ys := make([]float64, len(points))
	base := points[0].CollectedAt
	for i, p := range points {
		xs[i] = p.CollectedAt.Sub(base).Minutes()
		ys[i] = p.Value
	}
	slope := util.LinearSlope(xs, ys)
	limit := t.SlopeLimit
	if limit == 0 {
		limit = (c.MaxValue - c.MinValue) / mathMax(xs[len(xs)-1], 1) * 0.2
	}
	if (slope > limit && current.Value > c.MinValue) || (slope < -limit && current.Value < c.MaxValue) {
		out.Abnormal = true
		out.Severity = 1
		out.Score = 65
		out.Reason = fmt.Sprintf("persistent trend slope %.4f", slope)
	}
	return out
}
func TrendDirection(history []Sample) string {
	if len(history) < 2 {
		return "stable"
	}
	if history[len(history)-1].Value > history[0].Value {
		return "rising"
	}
	if history[len(history)-1].Value < history[0].Value {
		return "falling"
	}
	return "stable"
}
func Forecast(history []Sample, hours float64) float64 {
	if len(history) < 2 {
		return 0
	}
	xs := make([]float64, len(history))
	ys := make([]float64, len(history))
	base := history[0].CollectedAt
	for i, p := range history {
		xs[i] = p.CollectedAt.Sub(base).Hours()
		ys[i] = p.Value
	}
	return history[len(history)-1].Value + util.LinearSlope(xs, ys)*hours
}
func mathMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
