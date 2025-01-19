package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server/http/structures"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage"
)

const (
	dayMode = "day"
	weekMode = "week"
	monthMode = "month"
	betweenMode = "between"
)

type calendarHandler struct{
	app server.Application
	ctx context.Context
}

func validateEvent(event structures.Event) error {
	err := ""

	if event.Title == "" {
		err += "title is required field\n"
	}

	if event.EventDate == "" {
		err += "event_date is required field\n"
	}

	if event.DateSince == "" {
		err += "date_since is required field\n"
	}

	if event.UserID == 0 {
		err += "user_id is required field\n"
	}

	if err != "" {
		return errors.New(strings.Trim(err, "\n"))
	}

	return nil
}

func validateAndParseFindRequest(mode string, reqMap map[string]string) (map[string]time.Time, error) {
	switch mode {
	case dayMode, weekMode, monthMode:
		dateStr, ok := reqMap["date"]
		if !ok {
			return make(map[string]time.Time), errors.New("date is required")
		}
		
		date, err := time.Parse(time.DateOnly, dateStr)
		if err != nil {
			return make(map[string]time.Time), err
		}

		return map[string]time.Time{"date": date}, nil
	case betweenMode:
		startStr, ok := reqMap["start"]
		if !ok {
			return make(map[string]time.Time), errors.New("start is required")
		}

		endStr, ok := reqMap["end"]
		if !ok {
			return make(map[string]time.Time), errors.New("end is required")
		}
		
		start, err := time.Parse(time.DateOnly, startStr)
		if err != nil {
			return make(map[string]time.Time), err
		}

		end, err := time.Parse(time.DateOnly, endStr)
		if err != nil {
			return make(map[string]time.Time), err
		}

		return map[string]time.Time{"start": start, "end": end}, nil
	default:
		return make(map[string]time.Time), errors.New("unknown mode")
	}
}

func parseEventStruct (reqEvent structures.Event) (storage.Event, error) {
	if err := validateEvent(reqEvent); err != nil {
		return storage.Event{}, err
	}

	eventDate, err := time.Parse(time.DateTime, reqEvent.EventDate)
	if err != nil {
		return storage.Event{}, err
	}

	dateSince, err := time.Parse(time.DateTime, reqEvent.DateSince)
	if err != nil {
		return storage.Event{}, err
	}

	var notifyDate time.Time
	if reqEvent.NotifyDate != "" {
		notifyDate, err = time.Parse(time.DateTime, reqEvent.NotifyDate)
		if err != nil {
			return storage.Event{}, err
		}
	}

	event := storage.Event{
		ID: reqEvent.ID,
		Title: reqEvent.Title,
		EventDate: eventDate,
		DateSince: dateSince,
		Descr: reqEvent.Descr,
		UserID: reqEvent.UserID,
		NotifyDate: notifyDate,
	}

	return event, nil
}

func (h *calendarHandler) hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello, world"))
}

func (h *calendarHandler) create(w http.ResponseWriter, r *http.Request) {
	reqEvent := structures.Event{}

	if err := json.NewDecoder(r.Body).Decode(&reqEvent); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	event, err := parseEventStruct(reqEvent)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := h.app.CreateEvent(h.ctx, event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *calendarHandler) update(w http.ResponseWriter, r *http.Request) {
	reqEvent := structures.Event{}

	if err := json.NewDecoder(r.Body).Decode(&reqEvent); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	event, err := parseEventStruct(reqEvent)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.app.UpdateEvent(h.ctx, id, event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *calendarHandler) delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.app.DeleteEvent(h.ctx, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *calendarHandler) find(w http.ResponseWriter, r *http.Request) {
	reqMap := make(map[string]string)

	if err := json.NewDecoder(r.Body).Decode(&reqMap); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	var events []storage.Event
	var err error
	mode := r.PathValue("mode")

	parsedDates, err := validateAndParseFindRequest(mode, reqMap)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	switch mode {
	case dayMode:
		events, err = h.app.FindEventsForDay(h.ctx, parsedDates["date"])
	case weekMode:
		events, err = h.app.FindEventsForWeek(h.ctx, parsedDates["date"])
	case monthMode:
		events, err = h.app.FindEventsForMonth(h.ctx, parsedDates["date"])
	case betweenMode:
		events, err = h.app.FindEventsBetweenDates(h.ctx, parsedDates["start"], parsedDates["end"])
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}