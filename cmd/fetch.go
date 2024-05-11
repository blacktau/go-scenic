package cmd

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/blacktau/go-scenic/fetcher"
	"github.com/blacktau/go-scenic/providers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches images based on the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var pvdrs []providers.Provider
		cfgProviders := viper.GetStringSlice("providers")
		dest := viper.GetString("folder")

		slog.Info("configured providers", "providers", cfgProviders)

		if slices.Contains(cfgProviders, "bing") {
			bing := providers.NewBingProvider(10)
			pvdrs = append(pvdrs, &bing)
		}

		var papers []fetcher.Wallpaper

		for _, pvd := range pvdrs {
			slog.Info(fmt.Sprintf("indexing wallpapers from %s", pvd.GetName()))
			papers = append(papers, pvd.Fetch()...)
		}

		fetcher := fetcher.NewFetcher(dest)
		for _, prp := range papers {
			fetcher.Fetch(prp)
		}
	},
}
