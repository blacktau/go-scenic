package providers

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/blacktau/go-scenic/fetcher"
	"github.com/mmcdole/gofeed"
)

const userAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64, x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
const bingRssUrl string = "https://www.bing.com/HPImageArchive.aspx?format=rss&idx=0&n=%d"
const bingBaseUrl string = "https://www.bing.com%s"

var nameRegex regexp.Regexp = *regexp.MustCompile("id=(.+?)&")

type BingProvider struct {
	cnt  int
	Name string
}

func NewBingProvider(cnt int) BingProvider {
	return BingProvider{
		cnt: cnt,
	}
}

func (bp *BingProvider) GetName() string {
	return "Bing"
}

func (bp *BingProvider) Fetch() []fetcher.Wallpaper {

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(fmt.Sprintf(bingRssUrl, bp.cnt))

	if err != nil {
		slog.Error("error fetching bing feed", "error", err)
		return nil
	}

	var papers []fetcher.Wallpaper

	slog.Info("fetched feed items", "count", len(feed.Items))

	for _, itm := range feed.Items {
		matches := nameRegex.FindStringSubmatch(itm.Link)

		paper := fetcher.Wallpaper{
			Name: matches[1],
			Url:  fmt.Sprintf(bingBaseUrl, itm.Link),
		}

		papers = append(papers, paper)
	}

	return papers
}
