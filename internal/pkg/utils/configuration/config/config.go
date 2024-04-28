package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/spf13/viper"
)

var ConfigViper *viper.Viper = viper.New()

var (
	Color = configuration.ConfigProperty[bool]{
		KeyName:      "color",
		ViperStore:   ConfigViper,
		DefaultValue: true,
	}
)
var properties = []configuration.Property{
	Color,
}

func init() {
	locations := configuration.NewConfigLocations()

	ConfigViper.SetConfigName("config") // name of config file (without extension)
	ConfigViper.SetConfigType("yaml")
	ConfigViper.AddConfigPath(locations.ConfigPath)

	for _, property := range properties {
		property.Init()
	}

	ConfigViper.SafeWriteConfig()
}

func Clear() {
	for _, property := range properties {
		property.Clear()
	}
	SaveConfig()
}

func LoadConfig() {
	err := ConfigViper.ReadInConfig() // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		exit.Error(err)
	}
}

func SaveConfig() {
	err := ConfigViper.WriteConfig()
	if err != nil {
		exit.Error(err)
	}
}
