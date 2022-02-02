package sqlstorage

import (
	"banners-rotator/internal/storage"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	// database/sql implementation for sqlx.
	_ "github.com/lib/pq"
)

type Storage struct {
	store *sqlx.DB
}

func NewStorage(ctx context.Context, connectionString string) (*Storage, error) {
	store, err := sqlx.ConnectContext(ctx, "postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("new storage -> %w", err)
	}

	return &Storage{store}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	if err := s.store.PingContext(ctx); err != nil {
		return fmt.Errorf("connect to storage -> %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.store.Close()
}

func (s *Storage) CreateSlot(description string) (storage.Slot, error) {
	r := s.store.QueryRowx(
		"INSERT INTO slots (description) VALUES ($1) RETURNING id;",
		description,
	)
	var id int64
	if err := r.Scan(&id); err != nil {
		return storage.Slot{}, fmt.Errorf(
			"storage -> create slot -> %w (%s)",
			storage.ErrSlotNotCreated,
			err,
		)
	}

	return storage.Slot{ID: id, Description: description}, nil
}

func (s *Storage) CreateBanner(description string) (storage.Banner, error) {
	r := s.store.QueryRowx(
		"INSERT INTO banners (description) VALUES ($1) RETURNING id;",
		description,
	)
	var id int64
	if err := r.Scan(&id); err != nil {
		return storage.Banner{}, fmt.Errorf(
			"storage -> create banner -> %w (%s)",
			storage.ErrBannerNotCreated,
			err,
		)
	}

	return storage.Banner{ID: id, Description: description}, nil
}

func (s *Storage) CreateGroup(description string) (storage.Group, error) {
	r := s.store.QueryRowx(
		"INSERT INTO groups (description) VALUES ($1) RETURNING id;",
		description,
	)
	var id int64
	if err := r.Scan(&id); err != nil {
		return storage.Group{}, fmt.Errorf(
			"storage -> create group -> %w (%s)",
			storage.ErrGroupNotCreated,
			err,
		)
	}

	return storage.Group{ID: id, Description: description}, nil
}

func (s *Storage) CreateRotation(slotID, bannerID int64) error {
	_, err := s.store.Exec(
		"INSERT INTO rotations (slot_id, banner_id) VALUES ($1, $2);",
		slotID, bannerID,
	)
	if err != nil {
		return fmt.Errorf(
			"storage -> create rotation -> %w (%s)",
			storage.ErrRotationNotCreated,
			err,
		)
	}

	return nil
}

func (s *Storage) DeleteRotation(slotID, bannerID int64) error {
	r, err := s.store.Exec(
		"DELETE FROM rotations WHERE slot_id=$1 AND banner_id=$2;",
		slotID, bannerID,
	)
	if err != nil {
		return fmt.Errorf(
			"storage -> delete rotation -> %w (%s)",
			storage.ErrRotationNotDeleted,
			err,
		)
	}
	if count, err := r.RowsAffected(); err != nil || count == 0 {
		return fmt.Errorf(
			"storage -> delete rotation -> %w (not found)",
			storage.ErrRotationNotDeleted,
		)
	}

	return nil
}

func (s *Storage) CreateViewEvent(slotID, bannerID, groupID, date int64) error {
	_, err := s.store.Exec(
		"INSERT INTO views (slot_id, banner_id, group_id, date) VALUES ($1, $2, $3, $4);",
		slotID, bannerID, groupID, date,
	)
	if err != nil {
		return fmt.Errorf(
			"storage -> create view event -> %w (%s)",
			storage.ErrViewEventNotCreated,
			err,
		)
	}

	return nil
}

func (s *Storage) CreateClickEvent(slotID, bannerID, groupID, date int64) error {
	_, err := s.store.Exec(
		"INSERT INTO clicks (slot_id, banner_id, group_id, date) VALUES ($1, $2, $3, $4);",
		slotID, bannerID, groupID, date,
	)
	if err != nil {
		return fmt.Errorf(
			"storage -> create click event -> %w (%s)",
			storage.ErrClickEventNotCreated,
			err,
		)
	}

	return nil
}

func (s *Storage) NotViewedBanners(slotID int64) ([]storage.Banner, error) {
	var b []storage.Banner
	err := s.store.Select(
		&b,
		`SELECT *
				FROM banners
				WHERE id IN (
					SELECT banner_id
					FROM rotations
					WHERE slot_id = $1
						EXCEPT (SELECT banner_id FROM views WHERE slot_id = $1)
				)`,
		slotID,
	)
	if err != nil {
		return b, fmt.Errorf("storage -> not viewed banners -> %w", err)
	}

	return b, nil
}

func (s *Storage) SlotBanners(slotID int64) ([]storage.Banner, error) {
	var b []storage.Banner
	err := s.store.Select(
		&b,
		`SELECT *
				FROM banners
				WHERE id IN (
					SELECT banner_id
					FROM rotations
					WHERE slot_id = $1
				)`,
		slotID,
	)
	if err != nil {
		return b, fmt.Errorf("storage -> slot banners -> %w", err)
	}

	return b, nil
}

func (s *Storage) SlotViews(slotID int64) ([]storage.ViewEvent, error) {
	var v []storage.ViewEvent
	err := s.store.Select(
		&v,
		`SELECT * FROM views WHERE slot_id= $1`,
		slotID,
	)
	if err != nil {
		return v, fmt.Errorf("storage -> slot views -> %w", err)
	}

	return v, nil
}

func (s *Storage) SlotClicks(slotID int64) ([]storage.ClickEvent, error) {
	var v []storage.ClickEvent
	err := s.store.Select(
		&v,
		`SELECT * FROM clicks WHERE slot_id= $1`,
		slotID,
	)
	if err != nil {
		return v, fmt.Errorf("storage -> slot clicks -> %w", err)
	}

	return v, nil
}

func (s *Storage) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.store.Exec(query, args...)
}
