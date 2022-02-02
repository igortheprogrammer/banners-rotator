package bandit

import (
	"banners-rotator/internal/storage"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBandit_RandomBanner(t *testing.T) {
	bnd := NewBandit()

	t.Run("random banner", func(t *testing.T) {
		banners := getBanners(10)
		result, err := bnd.RandomBanner(banners)
		require.NoError(t, err)
		require.Greater(t, result.ID, int64(0))
	})

	t.Run("random banner without banners", func(t *testing.T) {
		banners := getBanners(0)
		result, err := bnd.RandomBanner(banners)
		require.ErrorIs(t, err, ErrEmptyBanners, "actual error is %s", err)
		require.Equal(t, result.ID, int64(0))
	})
}

func TestBandit_TopRatedBanner(t *testing.T) {
	bnd := NewBandit()

	t.Run("top rated banner", func(t *testing.T) {
		banners := getBanners(3)
		views := getViews(banners)
		clicks := getClicks(banners)
		clicks = append(clicks, clicks[len(clicks)-1])

		banner, err := bnd.TopRatedBanner(banners, views, clicks)
		require.NoError(t, err)
		require.Equal(t, clicks[len(clicks)-1].BannerID, banner.ID)
	})
}

func TestBandit_prepare(t *testing.T) {
	bnd := &Bandit{}

	t.Run("prepare", func(t *testing.T) {
		banners := getBanners(3)
		views := getViews(banners)
		clicks := getClicks(banners)
		cViews, cClicks, err := bnd.prepare(banners, views, clicks)
		require.NoError(t, err)
		require.Len(t, cViews, len(banners))
		require.Len(t, cClicks, len(banners))
	})

	t.Run("prepare with banner without views", func(t *testing.T) {
		banners := getBanners(3)
		views := getViews(banners)
		views = views[1:]
		clicks := getClicks(banners)
		_, _, err := bnd.prepare(banners, views, clicks)
		require.ErrorIs(t, err, ErrBannerWithoutViews, "actual error is %s", err)
	})
}

func TestBandit_bannerScore(t *testing.T) {
	bnd := &Bandit{}

	t.Run("banner score", func(t *testing.T) {
		views := 1000.0
		clicks := 100.0
		totalViews := 100000.0

		// Пока продолжается работа:
		// дёрнуть за ручку j, для которой максимальна величина
		// https://hsto.org/r/w1560/getpro/habr/post_images/979/478/120/979478120f7e7588c5dad40d61e97d04.png
		// где xj – средний доход от ручки j, nj – то, сколько раз мы дёрнули за ручку j, а n – то,
		// сколько раз мы дёргали за все ручки.
		expected := (clicks / views) + math.Sqrt(2*math.Log(totalViews)/views)
		result := bnd.bannerScore(views, clicks, totalViews)
		require.Equal(t, expected, result)
	})
}

func TestBandit_topBanners(t *testing.T) {
	bnd := &Bandit{}

	t.Run("top banners", func(t *testing.T) {
		banners := getBanners(3)
		views := getViews(banners)
		clicks := getClicks(banners)
		clicks = append(clicks, clicks[len(clicks)-1])
		cViews, cClicks, err := bnd.prepare(banners, views, clicks)
		require.NoError(t, err)

		scores := make(map[int64]scoresItem)
		for _, banner := range banners {
			scores[banner.ID] = scoresItem{
				banner,
				bnd.bannerScore(
					float64(cViews[banner.ID]),
					float64(cClicks[banner.ID]),
					float64(len(cClicks)),
				),
			}
		}

		top := bnd.topBanners(scores)
		require.Len(t, top, 1)
		require.Equal(t, clicks[len(clicks)-1].BannerID, top[0].ID)
	})
}

func getBanners(count int) []storage.Banner {
	var banners []storage.Banner

	for i := 0; i < count; i++ {
		banners = append(banners, storage.Banner{
			ID:          int64(i + 1),
			Description: fmt.Sprintf("Banner %d", i),
		})
	}

	return banners
}

func getViews(banners []storage.Banner) []storage.ViewEvent {
	var views []storage.ViewEvent

	for i, banner := range banners {
		views = append(views, storage.ViewEvent{
			SlotID:   int64(i + 1),
			BannerID: banner.ID,
			GroupID:  int64(i + 1),
			Date:     time.Now().Unix(),
		})
	}

	return views
}

func getClicks(banners []storage.Banner) []storage.ClickEvent {
	var clicks []storage.ClickEvent

	for i, banner := range banners {
		clicks = append(clicks, storage.ClickEvent{
			SlotID:   int64(i + 1),
			BannerID: banner.ID,
			GroupID:  int64(i + 1),
			Date:     time.Now().Unix(),
		})
	}

	return clicks
}
