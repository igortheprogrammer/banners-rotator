package rotator

import (
	"banners-rotator/internal/rmq"
	"banners-rotator/internal/storage"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

type App interface {
	CreateSlot(description string) (storage.Slot, error)
	CreateBanner(description string) (storage.Banner, error)
	CreateGroup(description string) (storage.Group, error)
	CreateRotation(slotID, bannerID int64) error
	DeleteRotation(slotID, bannerID int64) error
	CreateViewEvent(slotID, bannerID, groupID int64) error
	CreateClickEvent(slotID, bannerID, groupID int64) error
	BannerForSlot(slotID, groupID int64) (storage.Banner, error)
}

type Rotator struct {
	storage Storage
	p       *rmq.Producer
	b       Bandit
}

type Logger interface {
	Debug(msg string, values ...interface{})
	Info(msg string, values ...interface{})
	Warn(msg string, values ...interface{})
	Error(msg string, values ...interface{})
	Lgr() *zap.Logger
}

type Storage interface {
	CreateSlot(description string) (storage.Slot, error)
	CreateBanner(description string) (storage.Banner, error)
	CreateGroup(description string) (storage.Group, error)
	CreateRotation(slotID, bannerID int64) error
	DeleteRotation(slotID, bannerID int64) error
	CreateViewEvent(slotID, bannerID, groupID, date int64) error
	CreateClickEvent(slotID, bannerID, groupID, date int64) error
	NotViewedBanners(slotID int64) ([]storage.Banner, error)
	SlotBanners(slotID int64) ([]storage.Banner, error)
	SlotViews(slotID int64) ([]storage.ViewEvent, error)
	SlotClicks(slotID int64) ([]storage.ClickEvent, error)
}

type Bandit interface {
	RandomBanner(banners []storage.Banner) (storage.Banner, error)
	TopRatedBanner(
		banners []storage.Banner,
		views []storage.ViewEvent,
		clicks []storage.ClickEvent,
	) (storage.Banner, error)
}

func NewApp(s Storage, producer *rmq.Producer, bandit Bandit) App {
	return &Rotator{storage: s, p: producer, b: bandit}
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

	err = r.p.Publish(rmq.QMessage{
		Type:     "view",
		SlotID:   slotID,
		BannerID: bannerID,
		GroupID:  groupID,
		Date:     date,
	})
	if err != nil {
		return fmt.Errorf("rotator -> publish view event -> %w", err)
	}

	return nil
}

func (r *Rotator) CreateClickEvent(slotID, bannerID, groupID int64) error {
	date := time.Now().Unix()
	err := r.storage.CreateClickEvent(slotID, bannerID, groupID, date)
	if err != nil {
		return fmt.Errorf("rotator -> create view event -> %w", err)
	}

	err = r.p.Publish(rmq.QMessage{
		Type:     "click",
		SlotID:   slotID,
		BannerID: bannerID,
		GroupID:  groupID,
		Date:     date,
	})
	if err != nil {
		return fmt.Errorf("rotator -> publish click event -> %w", err)
	}

	return nil
}

func (r *Rotator) BannerForSlot(slotID, groupID int64) (storage.Banner, error) {
	notViewed, err := r.storage.NotViewedBanners(slotID)
	if err != nil {
		return storage.Banner{}, fmt.Errorf("rotator -> banner for slot -> %w", err)
	}
	if banner, err := r.b.RandomBanner(notViewed); err == nil && banner.ID > 0 {
		err = r.CreateViewEvent(slotID, banner.ID, groupID)
		if err != nil {
			return storage.Banner{}, fmt.Errorf("rotator -> banner for slot -> %w", err)
		}

		return banner, nil
	}

	banners, err := r.storage.SlotBanners(slotID)
	if err != nil {
		return storage.Banner{}, fmt.Errorf("rotator -> banner for slot -> %w", err)
	}

	views, err := r.storage.SlotViews(slotID)
	if err != nil {
		return storage.Banner{}, fmt.Errorf("rotator -> banner for slot -> %w", err)
	}

	clicks, err := r.storage.SlotClicks(slotID)
	if err != nil {
		return storage.Banner{}, fmt.Errorf("rotator -> banner for slot -> %w", err)
	}

	banner, err := r.b.TopRatedBanner(banners, views, clicks)
	if err != nil {
		return storage.Banner{}, fmt.Errorf("rotator -> banner for slot -> %w", err)
	}

	err = r.CreateViewEvent(slotID, banner.ID, groupID)
	if err != nil {
		return storage.Banner{}, fmt.Errorf("rotator -> banner for slot -> %w", err)
	}

	return banner, err
}
