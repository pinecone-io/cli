package presenters

import (
	"testing"

	"github.com/pinecone-io/go-pinecone/pinecone"
)

// TODO : clean up the old pointer print stuff
func TestPrintDescribeCollectionTable(t *testing.T) {
	var int64Zero int64 = 0
	// var int64NonZero int64 = 1000
	var int32Zero int32 = 0
	var int32NonZero int32 = 1000
	// var int64Nil = (*int64)(nil)
	// var int32Nil = (*int32)(nil)

	var tests = []struct {
		name        string
		size        int64
		dimension   int32
		vectorCount int32
		status      pinecone.CollectionStatus
	}{
		{"testcoll", int64Zero, int32Zero, int32Zero, pinecone.CollectionStatusReady},
		{"testcoll", int64Zero, int32NonZero, int32NonZero, pinecone.CollectionStatusInitializing},
		// {"testcoll", int64Nil, int32NonZero, int32NonZero, pinecone.CollectionStatusReady},
		// {"testcoll", int64NonZero, int32Nil, int32NonZero, pinecone.CollectionStatusReady},
		// {"testcoll", int64NonZero, int32NonZero, int32Nil, pinecone.CollectionStatusReady},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// type Collection struct {
			// 	Name        string
			// 	Size        *int64
			// 	Status      CollectionStatus
			// 	Dimension   *int32
			// 	VectorCount *int32
			// 	Environment string
			// }
			coll := &pinecone.Collection{
				Name:        tt.name,
				Dimension:   tt.dimension,
				VectorCount: tt.vectorCount,
				Size:        tt.size,
			}

			PrintDescribeCollectionTable(coll)
		})
	}
}
