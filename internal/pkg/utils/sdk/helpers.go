package sdk

import "github.com/pinecone-io/go-pinecone/v5/pinecone"

// Currently, passing a MetadataSchema field with "filterable: false" is not supported.
// We allow users to pass a slice of metadata fields, and then construct the MetadataSchema object from that.
func BuildMetadataSchema(schema []string) *pinecone.MetadataSchema {
	if len(schema) == 0 {
		return nil
	}

	metadataSchema := &pinecone.MetadataSchema{
		Fields: make(map[string]pinecone.MetadataSchemaField, len(schema)),
	}

	for _, field := range schema {
		metadataSchema.Fields[field] = pinecone.MetadataSchemaField{
			Filterable: true,
		}
	}

	return metadataSchema
}
