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

type TargetKnowledgeModel struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ChatHistory struct {
	History *models.KnowledgeModelChatHistory `json:"history"`
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
	TargetKm = configuration.MarshaledProperty[TargetKnowledgeModel]{
		KeyName:    "target_knowledge_model",
		ViperStore: StateViper,
		DefaultValue: &TargetKnowledgeModel{
			Name: "",
		},
	}
	ChatHist = configuration.MarshaledProperty[ChatHistory]{
		KeyName:    "chat_history",
		ViperStore: StateViper,
		DefaultValue: &ChatHistory{
			History: &models.KnowledgeModelChatHistory{},
		},
	}
)
var properties = []configuration.Property{
	TargetOrg,
	TargetProj,
	TargetKm,
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
