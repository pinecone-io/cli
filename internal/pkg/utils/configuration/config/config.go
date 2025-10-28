package config

import (
	"fmt"
	"strings"

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

	ConfigViper.SetEnvPrefix("pinecone")
	err := ConfigViper.BindEnv(Environment.KeyName)
	if err != nil {
		exit.Error().Err(err).Msg("Error binding environment to environment variable in config file")
	}

	err = validateEnvironment(Environment.Get())
	if err != nil {
		exit.Error().Err(err).Msg("Error validating environment")
	}
}

func validateEnvironment(env string) error {
	validEnvs := []string{"production", "staging"}
	for _, validEnv := range validEnvs {
		if env == validEnv {
			return nil
		}
	}
	quotedEnvs := make([]string, len(validEnvs))
	for i, validEnv := range validEnvs {
		quotedEnvs[i] = fmt.Sprintf("\"%s\"", validEnv)
	}
	return fmt.Errorf("invalid environment: \"%s\", must be one of %s", env, strings.Join(quotedEnvs, ", "))
}
