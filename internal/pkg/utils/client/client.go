package client

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

func NewPineconeClient() *pinecone.Client {
	key := config.ApiKey.Get()
	if key == "" {
		exit.Error(fmt.Errorf("API key not set. Please run %s or %s", style.Code("pinecone auth login"), style.Code("pinecone auth set-api-key")))
	}

	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey:    key,
		SourceTag: "pinecone-cli",
	})
	if err != nil {
		exit.Error(err)
	}

	return pc
}
