package state

import "github.com/pinecone-io/cli/internal/pkg/utils/pcio"

type TargetContext struct {
	Project string
	Org     string
}

func GetTargetContext() *TargetContext {
	return &TargetContext{
		Org:     TargetOrg.Get().Name,
		Project: TargetProj.Get().Name,
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
