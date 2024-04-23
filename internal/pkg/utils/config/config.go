package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/spf13/viper"
)

func InitConfigFile() {
	SetupDefaults()

	// We use SafeWriteConfig() instead of WriteConfig() to avoid overwriting
	// the config file if it already exists
	viper.SafeWriteConfig()
}

func SetupDefaults() {
	locations := NewConfigLocations()

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(locations.ConfigPath)

	viper.SetDefault("PINECONE_API_KEY", "")
}

func LoadConfig() {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		exit.Error(err)
	}
}

func SaveConfig() {
	err := viper.WriteConfig()
	if err != nil {
		exit.Error(err)
	}
}
