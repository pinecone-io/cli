package state

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
