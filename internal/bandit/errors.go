package bandit

import "errors"

var (
	ErrEmptyBanners       = errors.New("empty banners")
	ErrBannerWithoutViews = errors.New("banner without views")
)
