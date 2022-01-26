package sqlstorage

import (
	"banners-rotator/internal/storage"
	"context"
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
		return nil, fmt.Errorf("connect to storage -> %w", err)
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
	return storage.Slot{}, nil
}

func (s *Storage) CreateBanner(description string) (storage.Banner, error) {
	return storage.Banner{}, nil
}

func (s *Storage) CreateGroup(description string) (storage.Group, error) {
	return storage.Group{}, nil
}

func (s *Storage) CreateRotation(slotID, bannerID int64) error {
	return nil
}

func (s *Storage) DeleteRotation(slotID, bannerID int64) error {
	return nil
}

func (s *Storage) CreateViewEvent(slotID, bannerID, groupID, date int64) error {
	return nil
}

func (s *Storage) CreateClickEvent(slotID, bannerID, groupID, date int64) error {
	return nil
}
