package detector

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type RateDetector struct{ MaxRate float64 }

func (r RateDetector) Detect(current Sample, history []Sample, c Config) Result {
	out := Result{Score: 100, Rule: "rate"}
	if len(history) == 0 {
		return out
	}
	previous := history[len(history)-1]
	seconds := current.CollectedAt.Sub(previous.CollectedAt).Seconds()
	if seconds <= 0 {
		seconds = 1
	}
	rate := math.Abs(current.Value-previous.Value) / seconds
	if r.MaxRate <= 0 {
		r.MaxRate = (c.MaxValue - c.MinValue) / 60
	}
	if rate > r.MaxRate {
		out.Abnormal = true
		out.Severity = 2
		out.Score = 50
		out.Reason = fmt.Sprintf("rate %.4f exceeds %.4f", rate, r.MaxRate)
	}
	return out
}
func (r RateDetector) Rate(current, previous Sample) float64 {
	seconds := current.CollectedAt.Sub(previous.CollectedAt).Seconds()
	if seconds <= 0 {
		seconds = 1
	}
	return (current.Value - previous.Value) / seconds
}
func Recent(history []Sample, n int) []Sample {
	if n <= 0 || len(history) <= n {
		return append([]Sample(nil), history...)
	}
	out := append([]Sample(nil), history[len(history)-n:]...)
	sort.Slice(out, func(i, j int) bool { return out[i].CollectedAt.Before(out[j].CollectedAt) })
	return out
}
func DurationOf(history []Sample) time.Duration {
	if len(history) < 2 {
		return 0
	}
	return history[len(history)-1].CollectedAt.Sub(history[0].CollectedAt)
}
