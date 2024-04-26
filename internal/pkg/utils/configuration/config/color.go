package config

type ColorConfig struct {
	keyName string
}

var Color = ColorConfig{
	keyName: colorKey,
}

func (a ColorConfig) Set(newColorSetting bool) {
	ConfigViper.Set(colorKey, newColorSetting)
	SaveConfig()
}

func (a ColorConfig) Get() bool {
	colorSetting := ConfigViper.GetBool(colorKey)
	return colorSetting
}
