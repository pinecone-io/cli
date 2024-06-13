package state

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/models"
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

type TargetAssistant struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ChatHistory struct {
	History *models.AssistantChatHistory `json:"history"`
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
	TargetAsst = configuration.MarshaledProperty[TargetAssistant]{
		KeyName:    "target_assistant",
		ViperStore: StateViper,
		DefaultValue: &TargetAssistant{
			Name: "",
		},
	}
	ChatHist = configuration.MarshaledProperty[ChatHistory]{
		KeyName:    "chat_history",
		ViperStore: StateViper,
		DefaultValue: &ChatHistory{
			History: &models.AssistantChatHistory{},
		},
	}
)
var properties = []configuration.Property{
	TargetOrg,
	TargetProj,
	TargetAsst,
	ChatHist,
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
