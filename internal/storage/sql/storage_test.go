package sqlstorage

import (
	"banners-rotator/internal/storage"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestStorage_CreateSlot(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	s := Storage{store: sqlxDB}

	t.Run("create slot", func(t *testing.T) {
		desc := uuid.NewString()
		mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO slots (description) VALUES ($1) RETURNING id;`)).
			WithArgs(desc).
			WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))
		_, err = s.CreateSlot(desc)
		require.NoError(t, err)

		mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO slots (description) VALUES ($1) RETURNING id;`)).
			WithArgs(desc).
			WillReturnError(fmt.Errorf("test error"))
		_, err = s.CreateSlot(desc)
		require.ErrorIs(t, err, storage.ErrSlotNotCreated)
	})

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStorage_CreateBanner(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	s := Storage{store: sqlxDB}

	t.Run("create banner", func(t *testing.T) {
		desc := uuid.NewString()
		mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO banners (description) VALUES ($1) RETURNING id;`)).
			WithArgs(desc).
			WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))
		_, err = s.CreateBanner(desc)
		require.NoError(t, err)

		mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO banners (description) VALUES ($1) RETURNING id;`)).
			WithArgs(desc).
			WillReturnError(fmt.Errorf("test error"))
		_, err = s.CreateBanner(desc)
		require.ErrorIs(t, err, storage.ErrBannerNotCreated)
	})

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStorage_CreateGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	s := Storage{store: sqlxDB}

	t.Run("create group", func(t *testing.T) {
		desc := uuid.NewString()
		mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO groups (description) VALUES ($1) RETURNING id;`)).
			WithArgs(desc).
			WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))
		_, err = s.CreateGroup(desc)
		require.NoError(t, err)

		mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO groups (description) VALUES ($1) RETURNING id;`)).
			WithArgs(desc).
			WillReturnError(fmt.Errorf("test error"))
		_, err = s.CreateGroup(desc)
		require.ErrorIs(t, err, storage.ErrGroupNotCreated)
	})

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStorage_CreateRotation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	s := Storage{store: sqlxDB}

	t.Run("create rotation", func(t *testing.T) {
		mock.
			ExpectExec(regexp.QuoteMeta(`INSERT INTO rotations (slot_id, banner_id) VALUES ($1, $2);`)).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err = s.CreateRotation(1, 1)
		require.NoError(t, err)

		mock.
			ExpectExec(regexp.QuoteMeta(`INSERT INTO rotations (slot_id, banner_id) VALUES ($1, $2);`)).
			WithArgs(1, 1).
			WillReturnError(fmt.Errorf("test error"))
		err = s.CreateRotation(1, 1)
		require.ErrorIs(t, err, storage.ErrRotationNotCreated)
	})

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStorage_DeleteRotation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	s := Storage{store: sqlxDB}

	t.Run("delete rotation", func(t *testing.T) {
		mock.
			ExpectExec(regexp.QuoteMeta(`DELETE FROM rotations WHERE slot_id=$1 AND banner_id=$2;`)).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err = s.DeleteRotation(1, 1)
		require.NoError(t, err)

		mock.
			ExpectExec(regexp.QuoteMeta(`DELETE FROM rotations WHERE slot_id=$1 AND banner_id=$2;`)).
			WithArgs(1, 1).
			WillReturnError(fmt.Errorf("test error"))
		err = s.DeleteRotation(1, 1)
		require.ErrorIs(t, err, storage.ErrRotationNotDeleted)
	})

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStorage_CreateViewEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	s := Storage{store: sqlxDB}

	t.Run("create view event", func(t *testing.T) {
		mock.
			ExpectExec(regexp.QuoteMeta(
				`INSERT INTO views (slot_id, banner_id, group_id, date) VALUES ($1, $2, $3, $4);`,
			)).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err = s.CreateViewEvent(1, 1, 1, 1)
		require.NoError(t, err)

		mock.
			ExpectExec(regexp.QuoteMeta(
				`INSERT INTO views (slot_id, banner_id, group_id, date) VALUES ($1, $2, $3, $4);`,
			)).
			WithArgs(1, 1, 1, 1).
			WillReturnError(fmt.Errorf("test error"))
		err = s.CreateViewEvent(1, 1, 1, 1)
		require.ErrorIs(t, err, storage.ErrViewEventNotCreated)
	})

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStorage_CreateClickEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	s := Storage{store: sqlxDB}

	t.Run("create click event", func(t *testing.T) {
		mock.
			ExpectExec(regexp.QuoteMeta(
				`INSERT INTO clicks (slot_id, banner_id, group_id, date) VALUES ($1, $2, $3, $4);`,
			)).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err = s.CreateClickEvent(1, 1, 1, 1)
		require.NoError(t, err)

		mock.
			ExpectExec(regexp.QuoteMeta(
				`INSERT INTO clicks (slot_id, banner_id, group_id, date) VALUES ($1, $2, $3, $4);`,
			)).
			WithArgs(1, 1, 1, 1).
			WillReturnError(fmt.Errorf("test error"))
		err = s.CreateClickEvent(1, 1, 1, 1)
		require.ErrorIs(t, err, storage.ErrClickEventNotCreated)
	})

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
