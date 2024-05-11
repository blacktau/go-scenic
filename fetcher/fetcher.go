package fetcher

import (
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path"
)

type Wallpaper struct {
	Name string
	Url  string
}

type Fetcher struct {
	destination string
}

func NewFetcher(destination string) Fetcher {
	_, err := os.Stat(destination)
	if os.IsNotExist(err) {
		os.MkdirAll(destination, fs.FileMode(os.O_RDWR))
	}

	return Fetcher{
		destination,
	}
}

func (f *Fetcher) Fetch(w Wallpaper) {
	dest := path.Join(f.destination, w.Name)
	_, err := os.Stat(dest)

	if err == nil || !os.IsNotExist(err) {
		slog.Info("wallpaper already exists", "name", w.Name)
		return
	}

	slog.Info("downloading wall paper", "name", w.Name)
	res, err := http.Get(w.Url)
	if err != nil {
		slog.Error("failed to download wallpaper", "error", err, "url", w.Url)
		return
	}

	defer res.Body.Close()

	file, err := os.Create(dest)
	if err != nil {
		slog.Error("failed to create file", "error", err, "file", dest)
		return
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		slog.Error("failed to write to file", "error", err, "file", dest)
	}

	slog.Info("wallpaper saved", "name", w.Name)
}
