package state

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/spf13/viper"
)

var StateViper *viper.Viper = viper.New()

type TargetOrganization struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type TargetProject struct {
	Name string `json:"name"`
	Id   string `json:"global_id"`
}

var (
	TargetProj = configuration.MarshaledProperty[TargetProject]{
		KeyName:    "target_project",
		ViperStore: StateViper,
		DefaultValue: &TargetProject{
			Name: "",
			Id:   "",
		},
	}
	TargetOrg = configuration.MarshaledProperty[TargetOrganization]{
		KeyName:    "target_org",
		ViperStore: StateViper,
		DefaultValue: &TargetOrganization{
			Name: "",
			Id:   "",
		},
	}
)
var properties = []configuration.Property{
	TargetOrg,
	TargetProj,
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
