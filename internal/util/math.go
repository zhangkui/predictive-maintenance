package util

import "math"

func Clamp(v, low, high float64) float64 {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var total float64
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}
func WeightedMean(values, weights []float64) float64 {
	var total, weight float64
	for i, value := range values {
		if i >= len(weights) {
			break
		}
		total += value * weights[i]
		weight += weights[i]
	}
	if weight == 0 {
		return 0
	}
	return total / weight
}
func StdDev(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	mean := Mean(values)
	var sum float64
	for _, value := range values {
		sum += (value - mean) * (value - mean)
	}
	return math.Sqrt(sum / float64(len(values)))
}
func Percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]float64(nil), values...)
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j] < sorted[j-1]; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}
	p = Clamp(p, 0, 1)
	return sorted[int(float64(len(sorted)-1)*p)]
}
func SafeRate(current, previous, seconds float64) float64 {
	if seconds <= 0 {
		seconds = 1
	}
	return (current - previous) / seconds
}
func LinearSlope(xs, ys []float64) float64 {
	if len(xs) != len(ys) || len(xs) < 2 {
		return 0
	}
	mx, my := Mean(xs), Mean(ys)
	var numerator, denominator float64
	for i := range xs {
		numerator += (xs[i] - mx) * (ys[i] - my)
		denominator += (xs[i] - mx) * (xs[i] - mx)
	}
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}
func Normalize(v, low, high float64) float64 {
	if high <= low {
		return 0
	}
	return Clamp((v-low)/(high-low), 0, 1)
}
