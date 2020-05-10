package main

import (
	"fmt"

	"github.com/spf13/viper"
)

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
		debugPrint("In config keys loop")
		debugPrint(k)
		keys[i] = k
		i++
	}
	return keys
}

func configFeedSources() {
	categories := getConfigKeys("feeds")
	debugPrint("Categories:")
	for i := range categories {
		debugPrint(fmt.Sprintf("%d", cat))
		debugPrint(categories[cat])
	}
}
