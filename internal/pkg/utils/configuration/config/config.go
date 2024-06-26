package config

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/spf13/viper"
)

var ConfigViper *viper.Viper = viper.New()

var (
	Color = configuration.ConfigProperty[bool]{
		KeyName:      "color",
		ViperStore:   ConfigViper,
		DefaultValue: true,
	}
	Environment = configuration.ConfigProperty[string]{
		KeyName:      "environment",
		ViperStore:   ConfigViper,
		DefaultValue: "production",
	}
)
var properties = []configuration.Property{
	Color,
	Environment,
}

var configFile = configuration.ConfigFile{
	FileName:   "config",
	FileFormat: "yaml",
	Properties: properties,
	ViperStore: ConfigViper,
}

func init() {
	configFile.Init()
}
