package memorystorage

import (
	"context"
	"testing"
	"time"

	storagelib "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("insert right w/o ID", func(t *testing.T) {
		var storage = New()
		events := make(map[int]storagelib.Event)

		event := storagelib.Event{Title: "test1", EventDate: time.Now(), DateSince: time.Now(), Descr: "test descr", UserID: 1, NotifyDate: time.Now()}
		id, err := storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 1, id)
		event.ID = id
		events[id] = event

		event = storagelib.Event{Title: "test2", EventDate: time.Now(), DateSince: time.Now(), UserID: 1}
		id, err = storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 2, id)
		event.ID = id
		events[id] = event

		eventsInMem, err := storage.FindForDay(context.Background(), time.Now())
		require.Nil(t, err)
		require.Equal(t, len(events), len(eventsInMem))

		for _, v := range eventsInMem {
			require.Equal(t, v, events[v.ID])
		}
	})

	t.Run("insert right with ID", func(t *testing.T) {
		var storage = New()
		events := make(map[int]storagelib.Event)

		event := storagelib.Event{ID: 1, Title: "test1", EventDate: time.Now(), DateSince: time.Now(), Descr: "test descr", UserID: 1, NotifyDate: time.Now()}
		id, err := storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 1, id)
		events[id] = event

		event = storagelib.Event{ID: 2, Title: "test2", EventDate: time.Now(), DateSince: time.Now(), UserID: 1}
		id, err = storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 2, id)
		events[id] = event

		eventsInMem, err := storage.FindForDay(context.Background(), time.Now())
		require.Nil(t, err)
		require.Equal(t, len(events), len(eventsInMem))

		for _, v := range eventsInMem {
			require.Equal(t, v, events[v.ID])
		}
	})

	t.Run("insert wrong", func(t *testing.T) {
		var storage = New()

		event := storagelib.Event{ID: 1, Title: "test1", EventDate: time.Now(), DateSince: time.Now(), Descr: "test descr", UserID: 1, NotifyDate: time.Now()}
		id, err := storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 1, id)

		event = storagelib.Event{ID: 1, Title: "test2", EventDate: time.Now(), DateSince: time.Now(), UserID: 1}
		id, err = storage.Create(context.Background(), event)
		require.ErrorIs(t, err, ErrNotUniqueID)
		require.Equal(t, 0, id)
	})

	t.Run("update right", func(t *testing.T) {
		var storage = New()

		event := storagelib.Event{ID: 1, Title: "test1", EventDate: time.Now(), DateSince: time.Now(), Descr: "test descr", UserID: 1, NotifyDate: time.Now()}
		id, err := storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 1, id)

		newEvent := storagelib.Event{ID: 2, Title: "test2", EventDate: time.Now(), DateSince: time.Now(), UserID: 1}
		err = storage.Update(context.Background(), id, newEvent)
		require.Nil(t, err)

		eventsInMem, err := storage.FindForDay(context.Background(), time.Now())
		require.Nil(t, err)
		require.Equal(t, 1, len(eventsInMem))

		require.Equal(t, newEvent, eventsInMem[0])
	})

	t.Run("update wrong", func(t *testing.T) {
		var storage = New()

		event := storagelib.Event{ID: 1, Title: "test1", EventDate: time.Now(), DateSince: time.Now(), Descr: "test descr", UserID: 1, NotifyDate: time.Now()}
		id, err := storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 1, id)

		newEvent := storagelib.Event{ID: 2, Title: "test2", EventDate: time.Now(), DateSince: time.Now(), UserID: 1}
		err = storage.Update(context.Background(), 2, newEvent)
		require.ErrorIs(t, err, ErrEventNotFound)
	})

	t.Run("delete", func(t *testing.T) {
		var storage = New()

		event := storagelib.Event{ID: 1, Title: "test1", EventDate: time.Now(), DateSince: time.Now(), Descr: "test descr", UserID: 1, NotifyDate: time.Now()}
		id, err := storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 1, id)

		require.Nil(t, storage.Delete(context.Background(), id))
	})

	t.Run("selects", func(t *testing.T) {
		var storage = New()
		events := make(map[int]storagelib.Event)
		testTime := time.Date(2024, 12, 2, 0, 0, 0, 0, time.Now().Location())

		event := storagelib.Event{Title: "test1", EventDate: testTime, DateSince: time.Now(), Descr: "test descr", UserID: 1, NotifyDate: time.Now()}
		id, err := storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 1, id)
		event.ID = id
		events[id] = event

		event = storagelib.Event{Title: "test2", EventDate: testTime.Add(time.Hour * 24 * 3), DateSince: time.Now(), UserID: 1}
		id, err = storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 2, id)
		event.ID = id
		events[id] = event

		event = storagelib.Event{Title: "test3", EventDate: testTime.Add(time.Hour * 24 * 13), DateSince: time.Now(), UserID: 1}
		id, err = storage.Create(context.Background(), event)
		require.Nil(t, err)
		require.Equal(t, 3, id)
		event.ID = id
		events[id] = event

		eventsInMem, err := storage.FindForDay(context.Background(), testTime)
		require.Nil(t, err)
		require.Equal(t, 1, len(eventsInMem))

		for _, v := range eventsInMem {
			require.Equal(t, v, events[v.ID])
		}

		eventsInMem, err = storage.FindForWeek(context.Background(), testTime)
		require.Nil(t, err)
		require.Equal(t, 2, len(eventsInMem))

		for _, v := range eventsInMem {
			require.Equal(t, v, events[v.ID])
		}

		eventsInMem, err = storage.FindForMonth(context.Background(), testTime)
		require.Nil(t, err)
		require.Equal(t, len(events), len(eventsInMem))

		for _, v := range eventsInMem {
			require.Equal(t, v, events[v.ID])
		}
	})
}