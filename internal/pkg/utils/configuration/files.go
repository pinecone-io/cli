package configuration

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	FileName   string
	FileFormat string
	Properties []Property
	ViperStore *viper.Viper
}

func (c ConfigFile) Init() {
	fmt.Printf("Setting up config file: %s.%s\n", c.FileName, c.FileFormat)
	locations := NewConfigLocations()

	c.ViperStore.SetConfigName(c.FileName) // name of config file (without extension)
	c.ViperStore.SetConfigType(c.FileFormat)
	c.ViperStore.AddConfigPath(locations.ConfigPath)

	for _, property := range c.Properties {
		fmt.Printf("Setting default value for property: %v\n", property)
		property.Init()
	}
	c.ViperStore.SafeWriteConfig()
	c.LoadConfig()
}

func (c ConfigFile) Clear() {
	for _, property := range c.Properties {
		property.Clear()
	}
	c.Save()
}

func (c ConfigFile) LoadConfig() {
	fmt.Printf("Loading config file %s.%s\n", c.FileName, c.FileFormat)
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
