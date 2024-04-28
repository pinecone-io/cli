package state

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
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

var ConfigFile = configuration.ConfigFile{
	FileName:   "state",
	FileFormat: "yaml",
	Properties: properties,
	ViperStore: StateViper,
}

func init() {
	ConfigFile.Init()
}
