package configuration

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	FileName   string
	FileFormat string
	Properties []Property
	ViperStore *viper.Viper
}

func (c ConfigFile) Init() {
	log.Trace().Str("file_name", c.FileName).Str("file_format", c.FileFormat).Msg("Initializing config file")
	locations := NewConfigLocations()

	c.ViperStore.SetConfigName(c.FileName) // name of config file (without extension)
	c.ViperStore.SetConfigType(c.FileFormat)
	c.ViperStore.AddConfigPath(locations.ConfigPath)

	for _, property := range c.Properties {
		property.Init()
	}
	c.ViperStore.SafeWriteConfig()
	c.LoadConfig()
}

func (c ConfigFile) Clear() {
	log.Debug().Str("file_name", c.FileName).Msg("Clearing config file")
	for _, property := range c.Properties {
		property.Clear()
	}
	c.Save()
}

func (c ConfigFile) LoadConfig() {
	log.Debug().Str("file_name", c.FileName).Str("file_format", c.FileFormat).Msg("Loading config file")
	err := c.ViperStore.ReadInConfig() // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		exit.Error(err)
	}
}

func (c ConfigFile) Save() {
	err := c.ViperStore.WriteConfig()
	if err != nil {
		exit.Error(err)
	}
}