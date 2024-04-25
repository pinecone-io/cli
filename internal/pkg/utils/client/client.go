package client

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

func NewPineconeClient() *pinecone.Client {
	key := config.ApiKey.Get()
	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey:    key,
		SourceTag: "pinecone-cli",
	})
	if err != nil {
		exit.Error(err)
	}

	return pc
}
