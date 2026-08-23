package util

import "time"

func NowUTC() time.Time { return time.Now().UTC() }
func StartOfDay(t time.Time) time.Time {
	u := t.UTC()
	return time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
}
func EndOfDay(t time.Time) time.Time { return StartOfDay(t).Add(24*time.Hour - time.Nanosecond) }
func IsWeekend(t time.Time) bool     { return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday }
func AddCycle(t time.Time, cycle string, value int) time.Time {
	if value < 1 {
		value = 1
	}
	switch cycle {
	case "daily":
		return t.AddDate(0, 0, value)
	case "weekly":
		return t.AddDate(0, 0, 7*value)
	case "monthly":
		return t.AddDate(0, value, 0)
	case "yearly":
		return t.AddDate(value, 0, 0)
	default:
		return t.AddDate(0, 0, value)
	}
}
func ParseOrNow(v string) time.Time {
	if t, e := time.Parse(time.RFC3339, v); e == nil {
		return t
	}
	return NowUTC()
}
func HoursBetween(a, b time.Time) float64 { return b.Sub(a).Hours() }
