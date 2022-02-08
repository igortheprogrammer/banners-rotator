package integrationtests

import (
	"banners-rotator/internal/storage"
	sqlstorage "banners-rotator/internal/storage/sql"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func clearStorage(s *sqlstorage.Storage) error {
	if _, err := s.Exec(`TRUNCATE TABLE rotations RESTART IDENTITY CASCADE;`); err != nil {
		return fmt.Errorf("clear storage for test: %w", err)
	}

	if _, err := s.Exec(`TRUNCATE TABLE views RESTART IDENTITY CASCADE;`); err != nil {
		return fmt.Errorf("clear storage for test: %w", err)
	}

	if _, err := s.Exec(`TRUNCATE TABLE clicks RESTART IDENTITY CASCADE;`); err != nil {
		return fmt.Errorf("clear storage for test: %w", err)
	}

	if _, err := s.Exec(`TRUNCATE TABLE groups RESTART IDENTITY CASCADE;`); err != nil {
		return fmt.Errorf("clear storage for test: %w", err)
	}

	if _, err := s.Exec(`TRUNCATE TABLE banners RESTART IDENTITY CASCADE;`); err != nil {
		return fmt.Errorf("clear storage for test: %w", err)
	}

	if _, err := s.Exec(`TRUNCATE TABLE slots RESTART IDENTITY CASCADE;`); err != nil {
		return fmt.Errorf("clear storage for test: %w", err)
	}

	return nil
}

func getConnectionString() string {
	const connectionString = "host=localhost port=5432 user=postgres password=qwerty dbname=rotator sslmode=disable"
	cs := os.Getenv("TESTS_POSTGRES_DSN")
	if cs == "" {
		return connectionString
	}

	return cs
}

func TestStorage_CreateSlot(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("create slot", func(t *testing.T) {
		desc := uuid.NewString()

		slot, err := s.CreateSlot(desc)
		require.NoError(t, err)
		require.Equal(t, desc, slot.Description)
		require.Greater(t, slot.ID, int64(0))
	})
}

func TestStorage_CreateBanner(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("create banner", func(t *testing.T) {
		desc := uuid.NewString()

		banner, err := s.CreateBanner(desc)
		require.NoError(t, err)
		require.Equal(t, desc, banner.Description)
		require.Greater(t, banner.ID, int64(0))
	})
}

func TestStorage_CreateGroup(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("create group", func(t *testing.T) {
		desc := uuid.NewString()

		group, err := s.CreateGroup(desc)
		require.NoError(t, err)
		require.Equal(t, desc, group.Description)
		require.Greater(t, group.ID, int64(0))
	})
}

func TestStorage_CreateRotation(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("create rotation", func(t *testing.T) {
		desc := uuid.NewString()

		slot, err := s.CreateSlot(desc)
		require.NoError(t, err)
		banner, err := s.CreateBanner(desc)
		require.NoError(t, err)

		err = s.CreateRotation(slot.ID, banner.ID)
		require.NoError(t, err)

		r, err := s.Exec("SELECT * FROM rotations WHERE slot_id=$1 AND banner_id=$2;", slot.ID, banner.ID)
		require.NoError(t, err)
		count, err := r.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), count)

		err = s.CreateRotation(int64(1), int64(-1))
		require.ErrorIs(t, err, storage.ErrRotationNotCreated)
	})
}

func TestStorage_DeleteRotation(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("delete rotation", func(t *testing.T) {
		desc := uuid.NewString()

		slot, err := s.CreateSlot(desc)
		require.NoError(t, err)
		banner, err := s.CreateBanner(desc)
		require.NoError(t, err)

		err = s.CreateRotation(slot.ID, banner.ID)
		require.NoError(t, err)

		err = s.DeleteRotation(slot.ID, banner.ID)
		require.NoError(t, err)

		r, err := s.Exec("SELECT * FROM rotations WHERE slot_id=$1 AND banner_id=$2;", slot.ID, banner.ID)
		require.NoError(t, err)
		count, err := r.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(0), count)
	})
}

func TestStorage_CreateViewEvent(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("create view event", func(t *testing.T) {
		desc := uuid.NewString()

		slot, err := s.CreateSlot(desc)
		require.NoError(t, err)
		banner, err := s.CreateBanner(desc)
		require.NoError(t, err)
		group, err := s.CreateGroup(desc)
		require.NoError(t, err)

		date := time.Now().Unix()
		err = s.CreateViewEvent(slot.ID, banner.ID, group.ID, date)
		require.NoError(t, err)

		r, err := s.Exec(
			"SELECT * FROM views WHERE slot_id=$1 AND banner_id=$2 AND group_id=$3 AND date=$4;",
			slot.ID, banner.ID, group.ID, date,
		)
		require.NoError(t, err)
		count, err := r.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), count)

		err = s.CreateViewEvent(int64(-1), int64(-1), group.ID, date)
		require.ErrorIs(t, err, storage.ErrViewEventNotCreated)
	})
}

func TestStorage_CreateClickEvent(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("create view event", func(t *testing.T) {
		desc := uuid.NewString()

		slot, err := s.CreateSlot(desc)
		require.NoError(t, err)
		banner, err := s.CreateBanner(desc)
		require.NoError(t, err)
		group, err := s.CreateGroup(desc)
		require.NoError(t, err)

		date := time.Now().Unix()
		err = s.CreateClickEvent(slot.ID, banner.ID, group.ID, date)
		require.NoError(t, err)

		r, err := s.Exec(
			"SELECT * FROM clicks WHERE slot_id=$1 AND banner_id=$2 AND group_id=$3 AND date=$4;",
			slot.ID, banner.ID, group.ID, date,
		)
		require.NoError(t, err)
		count, err := r.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), count)

		err = s.CreateClickEvent(int64(-1), int64(-1), group.ID, date)
		require.ErrorIs(t, err, storage.ErrClickEventNotCreated)
	})
}

func TestStorage_NotViewedBanners(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("not viewed banners", func(t *testing.T) {
		desc := uuid.NewString()

		slot, err := s.CreateSlot(desc)
		require.NoError(t, err)
		banner, err := s.CreateBanner(desc)
		require.NoError(t, err)
		group, err := s.CreateGroup(desc)
		require.NoError(t, err)
		err = s.CreateRotation(slot.ID, banner.ID)
		require.NoError(t, err)
		date := time.Now().Unix()
		err = s.CreateViewEvent(slot.ID, banner.ID, group.ID, date)
		require.NoError(t, err)

		banners, err := s.NotViewedBanners(slot.ID)
		require.NoError(t, err)
		require.Len(t, *banners, 0)

		banner2, err := s.CreateBanner(desc)
		require.NoError(t, err)
		err = s.CreateRotation(slot.ID, banner2.ID)
		require.NoError(t, err)

		banners, err = s.NotViewedBanners(slot.ID)
		require.NoError(t, err)
		require.Len(t, *banners, 1)
		require.Equal(t, banner2.ID, (*banners)[0].ID)
	})
}

func TestStorage_SlotBanners(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("slot banners", func(t *testing.T) {
		desc := uuid.NewString()

		slot, err := s.CreateSlot(desc)
		require.NoError(t, err)
		slot2, err := s.CreateSlot(desc)
		require.NoError(t, err)
		banner, err := s.CreateBanner(desc)
		require.NoError(t, err)
		banner2, err := s.CreateBanner(desc)
		require.NoError(t, err)
		err = s.CreateRotation(slot.ID, banner.ID)
		require.NoError(t, err)
		err = s.CreateRotation(slot.ID, banner2.ID)
		require.NoError(t, err)
		err = s.CreateRotation(slot2.ID, banner2.ID)
		require.NoError(t, err)

		banners, err := s.SlotBanners(slot.ID)
		require.NoError(t, err)
		require.Len(t, *banners, 2)

		banners2, err := s.SlotBanners(slot2.ID)
		require.NoError(t, err)
		require.Len(t, *banners2, 1)
	})
}

func TestStorage_SlotViews(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("slot views", func(t *testing.T) {
		desc := uuid.NewString()

		slot, err := s.CreateSlot(desc)
		require.NoError(t, err)
		slot2, err := s.CreateSlot(desc)
		require.NoError(t, err)
		banner, err := s.CreateBanner(desc)
		require.NoError(t, err)
		banner2, err := s.CreateBanner(desc)
		require.NoError(t, err)
		group, err := s.CreateGroup(desc)
		require.NoError(t, err)
		group2, err := s.CreateGroup(desc)
		require.NoError(t, err)

		err = s.CreateViewEvent(slot.ID, banner.ID, group.ID, time.Now().Unix())
		require.NoError(t, err)
		err = s.CreateViewEvent(slot.ID, banner2.ID, group.ID, time.Now().Unix())
		require.NoError(t, err)
		err = s.CreateViewEvent(slot2.ID, banner.ID, group2.ID, time.Now().Unix())
		require.NoError(t, err)

		views, err := s.SlotViews(slot.ID)
		require.NoError(t, err)
		require.Len(t, *views, 2)

		views2, err := s.SlotViews(slot2.ID)
		require.NoError(t, err)
		require.Len(t, *views2, 1)
	})
}

func TestStorage_SlotClicks(t *testing.T) {
	s, err := sqlstorage.NewStorage(context.Background(), getConnectionString())
	require.NoError(t, err)
	err = s.Connect(context.Background())
	require.NoError(t, err)
	err = clearStorage(s)
	require.NoError(t, err)
	defer s.Close()

	t.Run("slot clicks", func(t *testing.T) {
		desc := uuid.NewString()

		slot, err := s.CreateSlot(desc)
		require.NoError(t, err)
		slot2, err := s.CreateSlot(desc)
		require.NoError(t, err)
		banner, err := s.CreateBanner(desc)
		require.NoError(t, err)
		banner2, err := s.CreateBanner(desc)
		require.NoError(t, err)
		group, err := s.CreateGroup(desc)
		require.NoError(t, err)
		group2, err := s.CreateGroup(desc)
		require.NoError(t, err)

		err = s.CreateClickEvent(slot.ID, banner.ID, group.ID, time.Now().Unix())
		require.NoError(t, err)
		err = s.CreateClickEvent(slot.ID, banner2.ID, group.ID, time.Now().Unix())
		require.NoError(t, err)
		err = s.CreateClickEvent(slot2.ID, banner.ID, group2.ID, time.Now().Unix())
		require.NoError(t, err)

		views, err := s.SlotClicks(slot.ID)
		require.NoError(t, err)
		require.Len(t, *views, 2)

		views2, err := s.SlotClicks(slot2.ID)
		require.NoError(t, err)
		require.Len(t, *views2, 1)
	})
}
