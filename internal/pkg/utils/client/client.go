package client

import (
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"os"
)

func NewPineconeClient() *pinecone.Client {
	key := os.Getenv("PINECONE_API_KEY")
	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: key,
	})
	if err != nil {
		exit.Error(err)
	}

	return pc
}
