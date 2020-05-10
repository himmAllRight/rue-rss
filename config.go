package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// FeedSourceConfig is an object that contains the information from the config
// file, required to create a new FeedSource item.
type FeedSourceConfig struct {
	Feedurl  string
	Category string
}

// Load and return a config from file.
func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/rue")
	viper.AddConfigPath("/etc/rue")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

// Grabs the sub configs keys of the config provided and returns them as a
// []string
func getConfigKeys(configName string) []string {
	configItems := viper.GetStringMapStringSlice(configName)
	keys := make([]string, len(configItems))

	i := 0
	for k := range configItems {
		keys[i] = k
		i++
	}
	return keys
}

// Grabs all the Feed Sources defined in the config and returns objects for
// them.
func configFeedSources() []FeedSourceConfig {
	configFeedSourceSlice := []FeedSourceConfig{}
	categories := getConfigKeys("feed_sources")
	for _, category := range categories {
		catConfig := viper.GetStringSlice(fmt.Sprintf("feed_sources.%s", category))
		for _, url := range catConfig {
			feedSourceConfig := FeedSourceConfig{url, category}
			configFeedSourceSlice = append(configFeedSourceSlice, feedSourceConfig)
		}
	}
	return configFeedSourceSlice
}
