package config

import (
	"github.com/spf13/viper"
)

const colorKey = "color"

func init() {
	viper.SetDefault(colorKey, true)
}

type ColorConfig struct {
	keyName string
}

var Color = ColorConfig{
	keyName: colorKey,
}

func (a ColorConfig) Set(newColorSetting bool) {
	viper.Set(colorKey, newColorSetting)
	SaveConfig()
}

func (a ColorConfig) Get() bool {
	colorSetting := viper.GetBool(colorKey)
	return colorSetting
}
