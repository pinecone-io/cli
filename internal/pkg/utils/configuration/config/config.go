package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/spf13/viper"
)

var ConfigViper *viper.Viper

const colorKey = "color"

func init() {
	ConfigViper = viper.New()
	ConfigViper.SetDefault(colorKey, true)

	locations := configuration.NewConfigLocations()

	ConfigViper.SetConfigName("config") // name of config file (without extension)
	ConfigViper.SetConfigType("yaml")
	ConfigViper.AddConfigPath(locations.ConfigPath)

	ConfigViper.SetDefault(colorKey, true)
	ConfigViper.SafeWriteConfig()
}

type ConfigPropertyBool struct {
	KeyName string
}

func (c ConfigPropertyBool) Set(value bool) {
	ConfigViper.Set(c.KeyName, value)
	SaveConfig()
}

func (c ConfigPropertyBool) Get() bool {
	return ConfigViper.GetBool(c.KeyName)
}

var (
	Color = ConfigPropertyBool{KeyName: colorKey}
)

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
