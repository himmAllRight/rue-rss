package main

import (
	"fmt"

	"github.com/spf13/viper"
)

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
	// Single Items (db src in this case)
	dbSrc := viper.GetString("db.src")
	debugPrint("Grabbed db src from config?")
	debugPrint(dbSrc)

	//	// Now how about listed items? (like feed urls)
	//	for _, feed := range feeds {
	//		debugPrint("In Cat Loop")
	//		for _, url := range feed {
	//			debugPrint(url)
	//		}
	//	}

}

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

func configFeedSources() []FeedSourceConfig {
	configFeedSourceSlice := []FeedSourceConfig{}
	categories := getConfigKeys("feeds")
	debugPrint("Categories:")
	for _, cat := range categories {
		//catConfig := viper.GetStringSlice(fmt.Sprintf("feeds.%s", categories[i]))
		catConfig := viper.GetStringSlice(fmt.Sprintf("feeds.%s", cat))
		debugPrint(fmt.Sprintf("Loaded Cat Config: %s", cat))
		for _, url := range catConfig {
			//feedSourceConfig := FeedSourceConfig(categories[i], url)
			debugPrint(url)
		}
	}
	return configFeedSourceSlice
}
