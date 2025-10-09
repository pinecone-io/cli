package state

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/spf13/viper"
)

var StateViper *viper.Viper = viper.New()

var (
	TargetProj = configuration.MarshaledProperty[TargetProject]{
		KeyName:    "target_project",
		ViperStore: StateViper,
		DefaultValue: TargetProject{
			Name: "",
			Id:   "",
		},
	}
	TargetOrg = configuration.MarshaledProperty[TargetOrganization]{
		KeyName:    "target_org",
		ViperStore: StateViper,
		DefaultValue: TargetOrganization{
			Name: "",
			Id:   "",
		},
	}
	AuthedUser = configuration.MarshaledProperty[TargetUser]{
		KeyName:    "user_context",
		ViperStore: StateViper,
		DefaultValue: TargetUser{
			AuthContext: AuthNone,
			Email:       "",
		},
	}
)
var properties = []configuration.Property{
	TargetOrg,
	TargetProj,
	AuthedUser,
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
