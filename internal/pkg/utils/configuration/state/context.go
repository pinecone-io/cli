package state

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
)

type TargetOrganization struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type TargetProject struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type AuthContext string

const (
	AuthNone           AuthContext = ""
	AuthUserToken      AuthContext = "user_token"
	AuthServiceAccount AuthContext = "service_account"
	AuthDefaultAPIKey  AuthContext = "default_api_key"
)

type TargetUser struct {
	AuthContext AuthContext `json:"auth_context"`
	Email       string      `json:"email"`
}

type TargetContext struct {
	Project       TargetProject      `json:"project"`
	Organization  TargetOrganization `json:"organization"`
	User          TargetUser         `json:"user"`
	DefaultAPIKey string             `json:"default_api_key"`
}

func GetTargetContext() *TargetContext {
	return &TargetContext{
		Organization: TargetOrg.Get(),
		Project:      TargetProj.Get(),
		User:         AuthedUser.Get(),
	}
}

func GetTargetOrgId() (string, error) {
	orgId := TargetOrg.Get().Id
	if orgId == "" {
		return "", pcio.Errorf("no target organization set")
	}
	return orgId, nil
}

func GetTargetProjectId() (string, error) {
	projId := TargetProj.Get().Id
	if projId == "" {
		return "", pcio.Errorf("no target project set")
	}
	return projId, nil
}

func GetTargetUserAuthContext() (string, error) {
	context := AuthedUser.Get()
	if context.AuthContext == AuthNone {
		return "", pcio.Errorf("no target authentication context set")
	}
	return string(context.AuthContext), nil
}

func GetTargetUserEmail() string {
	return AuthedUser.Get().Email
}
