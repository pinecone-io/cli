package state

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/spf13/viper"
)

var StateViper *viper.Viper = viper.New()

var (
	TargetProjectName = configuration.ConfigProperty[string]{
		KeyName:      "target_project_name",
		ViperStore:   StateViper,
		DefaultValue: "",
	}
	TargetProjectGlobalId = configuration.ConfigProperty[string]{
		KeyName:      "target_project_global_id",
		ViperStore:   StateViper,
		DefaultValue: "",
	}
	TargetOrgName = configuration.ConfigProperty[string]{
		KeyName:      "target_org_name",
		ViperStore:   StateViper,
		DefaultValue: "",
	}
	TargetOrgId = configuration.ConfigProperty[string]{
		KeyName:      "target_org_id",
		ViperStore:   StateViper,
		DefaultValue: "",
	}
	HumanMode = configuration.ConfigProperty[bool]{
		KeyName:      "human_mode",
		ViperStore:   StateViper,
		DefaultValue: true,
	}
)
var properties = []configuration.Property{
	TargetProjectName,
	TargetProjectGlobalId,
	TargetOrgName,
	TargetOrgId,
	HumanMode,
}

func init() {
	locations := configuration.NewConfigLocations()

	StateViper.SetConfigName("state")
	StateViper.SetConfigType("yaml")
	StateViper.AddConfigPath(locations.ConfigPath)

	for _, property := range properties {
		property.Init()
	}

	StateViper.SafeWriteConfig()
}

func Clear() {
	for _, property := range properties {
		property.Clear()
	}
	SaveState()
}

func LoadState() {
	err := StateViper.ReadInConfig() // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		exit.Error(err)
	}
}

func SaveState() {
	err := StateViper.WriteConfig()
	if err != nil {
		exit.Error(err)
	}
}

type TargetContext struct {
	Api     string
	Project string
	Org     string
}

func GetTargetContext() *TargetContext {
	LoadState()
	return &TargetContext{
		Api:     "https://api.pinecone.io",
		Org:     TargetOrgName.Get(),
		Project: TargetProjectName.Get(),
	}
}
