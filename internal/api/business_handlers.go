package api

import (
	"errors"
	"net/http"
	"predictive-maintenance/internal/model"
	"strconv"
	"time"
)

func (a *API) sensorConfigsHandler(w http.ResponseWriter, r *http.Request) {
	device, e := strconv.ParseUint(r.URL.Query().Get("deviceId"), 10, 64)
	if e != nil || device == 0 {
		fail(w, 400, errors.New("invalid deviceId"))
		return
	}
	if r.Method == http.MethodGet {
		items, e := a.configs.List(r.Context(), device)
		if e != nil {
			fail(w, 500, e)
			return
		}
		write(w, 200, items)
		return
	}
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		fail(w, 405, errors.New("method not allowed"))
		return
	}
	var c model.SensorConfig
	if e = decode(r, &c); e != nil {
		fail(w, 400, e)
		return
	}
	c.DeviceID = device
	v, e := a.configs.Save(r.Context(), c)
	if e != nil {
		fail(w, 400, e)
		return
	}
	write(w, 200, v)
}
func (a *API) healthHistoryHandler(w http.ResponseWriter, r *http.Request) {
	device, e := strconv.ParseUint(r.URL.Query().Get("deviceId"), 10, 64)
	if e != nil || device == 0 {
		fail(w, 400, errors.New("invalid deviceId"))
		return
	}
	to := time.Now().UTC()
	from := to.Add(-30 * 24 * time.Hour)
	if value := r.URL.Query().Get("from"); value != "" {
		from, e = time.Parse(time.RFC3339, value)
		if e != nil {
			fail(w, 400, e)
			return
		}
	}
	if value := r.URL.Query().Get("to"); value != "" {
		to, e = time.Parse(time.RFC3339, value)
		if e != nil {
			fail(w, 400, e)
			return
		}
	}
	items, e := a.healthHistory.List(r.Context(), device, from, to, 200)
	if e != nil {
		fail(w, 500, e)
		return
	}
	write(w, 200, items)
}
func (a *API) groupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		items, e := a.groups.List(r.Context())
		if e != nil {
			fail(w, 500, e)
			return
		}
		write(w, 200, items)
		return
	}
	if r.Method != http.MethodPost {
		fail(w, 405, errors.New("method not allowed"))
		return
	}
	var group model.DeviceGroup
	if e := decode(r, &group); e != nil {
		fail(w, 400, e)
		return
	}
	created, e := a.groups.Create(r.Context(), group)
	if e != nil {
		fail(w, 400, e)
		return
	}
	write(w, 201, created)
}
