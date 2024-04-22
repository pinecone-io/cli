package client

import (
	"fmt"
	"os"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

func NewPineconeClient() *pinecone.Client {
	key := os.Getenv("PINECONE_API_KEY")
	fmt.Println("list called with key:", key)

	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: key,
	})
	if err != nil {
		exit.Error(err)
	}

	return pc
}