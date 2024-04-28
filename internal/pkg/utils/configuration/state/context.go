package state

type TargetContext struct {
	Api     string
	Project string
	Org     string
}

func GetTargetContext() *TargetContext {
	return &TargetContext{
		Api:     "https://api.pinecone.io",
		Org:     TargetOrgName.Get(),
		Project: TargetProjectName.Get(),
	}
}
