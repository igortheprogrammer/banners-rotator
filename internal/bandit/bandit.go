package bandit

import (
	"banners-rotator/internal/rotator"
	"banners-rotator/internal/storage"
	"math"
	"math/rand"
)

type scoresItem struct {
	banner storage.Banner
	score  float64
}

type Bandit struct{}

func NewBandit() rotator.Bandit {
	return &Bandit{}
}

func (b *Bandit) RandomBanner(banners []storage.Banner) (storage.Banner, error) {
	if len(banners) > 0 {
		i := rand.Intn(len(banners))
		return banners[i], nil
	}

	return storage.Banner{}, ErrEmptyBanners
}

func (b *Bandit) TopRatedBanner(
	banners []storage.Banner,
	views []storage.ViewEvent,
	clicks []storage.ClickEvent,
) (storage.Banner, error) {
	cViews, cClicks, err := b.prepare(banners, views, clicks)
	if err != nil {
		return storage.Banner{}, err
	}

	scores := make(map[int64]scoresItem)
	for _, banner := range banners {
		scores[banner.ID] = scoresItem{
			banner,
			b.bannerScore(
				float64(cViews[banner.ID]),
				float64(cClicks[banner.ID]),
				float64(len(cClicks)),
			),
		}
	}

	top := b.topBanners(scores)
	banner, err := b.RandomBanner(top)

	return banner, err
}

func (b Bandit) prepare(
	banners []storage.Banner,
	views []storage.ViewEvent,
	clicks []storage.ClickEvent,
) (map[int64]int, map[int64]int, error) {
	cachedViews := make(map[int64]int)
	cachedClicks := make(map[int64]int)

	for _, view := range views {
		cachedViews[view.BannerID]++
	}

	for _, click := range clicks {
		cachedClicks[click.BannerID]++
	}

	for _, banner := range banners {
		if cachedViews[banner.ID] == 0 {
			return nil, nil, ErrBannerWithoutViews
		}
	}

	return cachedViews, cachedClicks, nil
}

func (b *Bandit) bannerScore(views float64, clicks float64, totalViews float64) float64 {
	clickViewRate := clicks / views
	banditRate := math.Sqrt(2 * math.Log(totalViews) / views)

	return clickViewRate + banditRate
}

func (b *Bandit) topBanners(scores map[int64]scoresItem) []storage.Banner {
	max := 0.0
	for _, item := range scores {
		if item.score > max {
			max = item.score
		}
	}

	var banners []storage.Banner
	for _, item := range scores {
		if item.score == max {
			banners = append(banners, item.banner)
		}
	}

	return banners
}
