package state

type TargetContext struct {
	Api            string
	ApiStaging     string
	Project        string
	Org            string
	KnowledgeModel string
}

func GetTargetContext() *TargetContext {
	return &TargetContext{
		Api:            "https://api.pinecone.io",
		ApiStaging:     "https://api-staging.pinecone.io",
		Org:            TargetOrg.Get().Name,
		Project:        TargetProj.Get().Name,
		KnowledgeModel: TargetKm.Get().Name,
	}
}
