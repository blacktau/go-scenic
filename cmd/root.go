package cmd

import (
	"log/slog"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var logToConsole bool
var logger *slog.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-scenic",
	Short: "A downloader of desktop backgrounds",
	Long: `go-scenic is a utility that downloads pictures suitable
for desktop background from various online sources.
The sources are configurable.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.AddCommand(fetchCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("initialising go-scenic")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// try XDG_CONFIG_HOME
		viper.AddConfigPath("$XDG_CONFIG_HOME/go-scenic")
		viper.SetConfigType("yaml")
		viper.SetConfigName("go-scenic")
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetDefault("log-level", "info")
	viper.SetDefault("providers", []string{"bing", "unsplash"})
	viper.SetDefault("initial-fetch", 10)
	viper.SetDefault("folder", "$HOME/Pictures/wallpapers")
	viper.SetDefault("provider-folders", true)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if cfgFile != "" {
			logger.Error("config file not found", "config-file", cfgFile)
			os.Exit(-1)
		}
		// if the config file is not found - create the default
		logger.Warn("No config file found. writing default config")
		cfgHome := os.Getenv("XDG_CONFIG_HOME")
		cfgPath := path.Join(cfgHome, "go-scenic")

		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			os.Mkdir(cfgPath, os.ModePerm)
		}

		if err := viper.SafeWriteConfig(); err != nil {
			logger.Error("failed to save default config", "err", err)
			os.Exit(-1)
		}
	}

	logger.Info("read config", "config-file", viper.ConfigFileUsed())
}
