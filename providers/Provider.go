package providers

import "github.com/blacktau/go-scenic/fetcher"

type Provider interface {
	Fetch() []fetcher.Wallpaper
	GetName() string
}
