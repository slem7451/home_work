package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server/http/structures" //nolint:depguard
	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestEventLogic(t *testing.T) {
	time.Sleep(7*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	event := structures.Event{
		Title: "test title",
		Descr: "test descr",
		EventDate: time.Now().Format(time.DateTime),
		DateSince: time.Now().Format(time.DateTime),
		UserID: 1,
	}

	eventB, err := json.Marshal(event)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, "POST", "http://calendar:8080/event", bytes.NewBuffer(eventB))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	idB := make([]byte, 255)
	n, _ := resp.Body.Read(idB)
	idB = idB[:n]

	var eventId map[string]int
	err = json.Unmarshal(idB, &eventId)
	require.NoError(t, err)

	event = structures.Event{
		Title: "test title update",
		Descr: "test descr",
		EventDate: time.Now().Format(time.DateTime),
		DateSince: time.Now().Format(time.DateTime),
		UserID: 1,
	}

	eventB, err = json.Marshal(event)
	require.NoError(t, err)

	req, err = http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("http://calendar:8080/event/%d", eventId["id"]), bytes.NewBuffer(eventB))
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	date := map[string]string {"date": time.Now().Format(time.DateOnly)}
	dateB, err := json.Marshal(date)
	require.NoError(t, err)

	req, err = http.NewRequestWithContext(ctx, "GET", "http://calendar:8080/event/day", bytes.NewBuffer(dateB))
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequestWithContext(ctx, "GET", "http://calendar:8080/event/week", bytes.NewBuffer(dateB))
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequestWithContext(ctx, "GET", "http://calendar:8080/event/month", bytes.NewBuffer(dateB))
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("http://calendar:8080/event/%d", eventId["id"]), nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}