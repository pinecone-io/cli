package state

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/spf13/viper"
)

var StateViper *viper.Viper

const targetProjectName string = "target_project_name"
const targetProjectGlobalId string = "target_project_global_id"
const targetOrgName string = "target_org_name"
const targetOrgId string = "target_org_id"
const humanMode string = "human_mode"

type ConfigProperty struct {
	KeyName string
}

func (c ConfigProperty) Set(value string) {
	StateViper.Set(c.KeyName, value)
	SaveState()
}

func (c ConfigProperty) Get() string {
	return StateViper.GetString(c.KeyName)
}

type ConfigPropertyBool struct {
	KeyName string
}

func (c ConfigPropertyBool) Set(value bool) {
	StateViper.Set(c.KeyName, value)
	SaveState()
}

func (c ConfigPropertyBool) Get() bool {
	return StateViper.GetBool(c.KeyName)
}

var (
	TargetProjectName     = ConfigProperty{KeyName: targetProjectName}
	TargetProjectGlobalId = ConfigProperty{KeyName: targetProjectGlobalId}

	TargetOrgName = ConfigProperty{KeyName: targetOrgName}
	TargetOrgId   = ConfigProperty{KeyName: targetOrgId}

	HumanMode = ConfigPropertyBool{KeyName: humanMode}
)

func init() {
	StateViper = viper.New()
	locations := configuration.NewConfigLocations()

	StateViper.SetConfigName("state")
	StateViper.SetConfigType("yaml")
	StateViper.AddConfigPath(locations.ConfigPath)

	StateViper.SetDefault(targetProjectName, "")
	StateViper.SetDefault(targetProjectGlobalId, "")
	StateViper.SetDefault(targetOrgName, "")
	StateViper.SetDefault(targetOrgId, "")
	StateViper.SetDefault(humanMode, false)
	StateViper.SafeWriteConfig()
}

func Clear() {
	TargetProjectName.Set("")
	TargetProjectGlobalId.Set("")
	TargetOrgName.Set("")
	TargetOrgId.Set("")
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
