package state

type TargetContext struct {
	Api     string
	Project string
	Org     string
}

func GetTargetContext() *TargetContext {
	return &TargetContext{
		Api:     "https://api.pinecone.io",
		Org:     TargetOrg.Get().Name,
		Project: TargetProj.Get().Name,
	}
}
