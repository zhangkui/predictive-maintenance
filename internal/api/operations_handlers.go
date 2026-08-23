package api

import (
	"errors"
	"net/http"
	"predictive-maintenance/internal/service"
	"strconv"
	"time"
)

func (a *API) runtimeHandler(w http.ResponseWriter, r *http.Request) {
	device, e := strconv.ParseUint(r.URL.Query().Get("deviceId"), 10, 64)
	if e != nil || device == 0 {
		fail(w, 400, errors.New("valid deviceId is required"))
		return
	}
	if r.Method == http.MethodPost {
		var q struct {
			Date   string  `json:"date"`
			Hours  float64 `json:"hours"`
			Source string  `json:"source"`
		}
		if e = decode(r, &q); e != nil {
			fail(w, 400, e)
			return
		}
		date := time.Now().UTC()
		if q.Date != "" {
			date, e = time.Parse("2006-01-02", q.Date)
			if e != nil {
				fail(w, 400, e)
				return
			}
		}
		if e = a.runtime.Record(r.Context(), service.RuntimeRecord{DeviceID: device, RunDate: date, Hours: q.Hours, Source: q.Source}); e != nil {
			fail(w, 400, e)
			return
		}
		write(w, 201, map[string]any{"recorded": true})
		return
	}
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -30)
	items, e := a.runtime.Daily(r.Context(), device, from, to)
	if e != nil {
		fail(w, 500, e)
		return
	}
	total, _ := a.runtime.Total(r.Context(), device, from, to)
	write(w, 200, map[string]any{"items": items, "totalHours": total})
}
func (a *API) retentionHandler(w http.ResponseWriter, r *http.Request) {
	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if r.Method == http.MethodGet {
		v, e := a.retention.Preview(r.Context(), days)
		if e != nil {
			fail(w, 500, e)
			return
		}
		write(w, 200, v)
		return
	}
	if r.Method != http.MethodPost {
		fail(w, 405, errors.New("method not allowed"))
		return
	}
	v, e := a.retention.Cleanup(r.Context(), days)
	if e != nil {
		fail(w, 500, e)
		return
	}
	write(w, 200, v)
}
