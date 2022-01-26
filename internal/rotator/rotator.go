package rotator

import (
	"banners-rotator/internal/storage"
	"fmt"
	"strings"
	"time"
)

type App interface {
	CreateSlot(description string) (storage.Slot, error)
	CreateBanner(description string) (storage.Banner, error)
	CreateGroup(description string) (storage.Group, error)
	CreateRotation(slotID, bannerID int64) error
	DeleteRotation(slotID, bannerID int64) error
	CreateViewEvent(slotID, bannerID, groupID int64) error
	CreateClickEvent(slotID, bannerID, groupID int64) error
	BannerForSlot(slotID, groupID int64) (int64, error)
}

type Rotator struct {
	storage Storage
}

type Logger interface {
	Debug(msg string, values ...interface{})
	Info(msg string, values ...interface{})
	Warn(msg string, values ...interface{})
	Error(msg string, values ...interface{})
}

type Storage interface {
	CreateSlot(description string) (storage.Slot, error)
	CreateBanner(description string) (storage.Banner, error)
	CreateGroup(description string) (storage.Group, error)
	CreateRotation(slotID, bannerID int64) error
	DeleteRotation(slotID, bannerID int64) error
	CreateViewEvent(slotID, bannerID, groupID, date int64) error
	CreateClickEvent(slotID, bannerID, groupID, date int64) error
}

func NewApp(s Storage) App {
	return &Rotator{storage: s}
}

func (r *Rotator) CreateSlot(description string) (storage.Slot, error) {
	return r.storage.CreateSlot(strings.TrimSpace(description))
}

func (r *Rotator) CreateBanner(description string) (storage.Banner, error) {
	return r.storage.CreateBanner(strings.TrimSpace(description))
}

func (r *Rotator) CreateGroup(description string) (storage.Group, error) {
	return r.storage.CreateGroup(strings.TrimSpace(description))
}

func (r *Rotator) CreateRotation(slotID, bannerID int64) error {
	return r.storage.CreateRotation(slotID, bannerID)
}

func (r *Rotator) DeleteRotation(slotID, bannerID int64) error {
	return r.storage.DeleteRotation(slotID, bannerID)
}

func (r *Rotator) CreateViewEvent(slotID, bannerID, groupID int64) error {
	date := time.Now().Unix()
	err := r.storage.CreateViewEvent(slotID, bannerID, groupID, date)
	if err != nil {
		return fmt.Errorf("rotator -> create view event -> %w", err)
	}

	// TODO: send to RMQ
	return nil
}

func (r *Rotator) CreateClickEvent(slotID, bannerID, groupID int64) error {
	date := time.Now().Unix()
	err := r.storage.CreateClickEvent(slotID, bannerID, groupID, date)
	if err != nil {
		return fmt.Errorf("rotator -> create view event -> %w", err)
	}

	// TODO: send to RMQ
	return nil
}

func (r *Rotator) BannerForSlot(slotID, groupID int64) (int64, error) {
	bannerID := int64(0)
	err := r.CreateViewEvent(slotID, bannerID, groupID)
	if err != nil {
		return 0, fmt.Errorf("rotator -> banner for slot -> %w", err)
	}

	return 0, err
}
