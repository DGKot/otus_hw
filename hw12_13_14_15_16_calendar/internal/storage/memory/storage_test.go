package memorystorage

import (
	"testing"
	"time"

	"github.com/DGKot/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		stor := New()
		ev := storage.Event{
			Title:    "test title",
			Datetime: time.Now(),
			Duration: time.Minute * 30,
		}
		id, err := stor.Create(ev)
		require.NoError(t, err)
		ev.ID = id
		evBuf, err := stor.Get(id)
		require.NoError(t, err)
		require.Equal(t, ev, evBuf)
	})

	t.Run("busy time", func(t *testing.T) {
		stor := New()
		ev := storage.Event{
			Title:    "test title",
			Datetime: time.Now(),
			Duration: time.Minute * 30,
		}
		id, _ := stor.Create(ev)
		ev.ID = id

		evBuf := storage.Event{
			Title:    "test title second",
			Datetime: time.Now(),
			Duration: time.Minute * 10,
		}
		id, err := stor.Create(evBuf)
		require.EqualError(t, err, storage.ErrDateBusy.Error())
		require.Equal(t, "", id)
	})

	t.Run("delete", func(t *testing.T) {
		stor := New()
		ev := storage.Event{
			Title:    "test title",
			Datetime: time.Now(),
			Duration: time.Minute * 30,
		}
		id, _ := stor.Create(ev)
		ev.ID = id

		err := stor.Delete(id)
		require.NoError(t, err)
		event, err := stor.Get(id)
		require.EqualError(t, err, storage.ErrEventNotFound.Error())
		require.Equal(t, event, storage.Event{})
	})

	t.Run("edit", func(t *testing.T) {
		stor := New()
		ev := storage.Event{
			Title:    "test title",
			Datetime: time.Now(),
			Duration: time.Minute * 30,
		}
		id, _ := stor.Create(ev)
		ev.ID = id
		ev.Title = "change title"

		err := stor.Edit(ev)
		require.NoError(t, err)
		event, _ := stor.Get(id)
		require.Equal(t, event.Title, "change title")
	})
}
