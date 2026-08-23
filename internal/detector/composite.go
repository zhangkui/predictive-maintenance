package detector

import "sort"

type CompositeDetector struct{ Detectors []Detector }

func (c CompositeDetector) Detect(current Sample, history []Sample, config Config) Result {
	if len(c.Detectors) == 0 {
		return Result{Score: 100, Rule: "none"}
	}
	results := make([]Result, 0, len(c.Detectors))
	for _, d := range c.Detectors {
		results = append(results, d.Detect(current, history, config))
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Severity > results[j].Severity })
	best := results[0]
	var abnormal int
	var score float64
	for _, r := range results {
		if r.Abnormal {
			abnormal++
		}
		score += r.Score
	}
	if abnormal >= 2 {
		best.Abnormal = true
		if best.Severity < 2 {
			best.Severity = 2
		}
		best.Rule = "composite"
	}
	best.Score = score / float64(len(results))
	return best
}
func (c CompositeDetector) Healthy(current Sample, history []Sample, config Config) bool {
	return !c.Detect(current, history, config).Abnormal
}
