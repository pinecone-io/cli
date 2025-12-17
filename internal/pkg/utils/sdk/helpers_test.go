package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildMetadataSchema(t *testing.T) {
	t.Run("empty schema returns nil", func(t *testing.T) {
		assert.Nil(t, BuildMetadataSchema([]string{}))
	})

	t.Run("creates filterable fields", func(t *testing.T) {
		fields := []string{"field1", "field2"}

		schema := BuildMetadataSchema(fields)
		if assert.NotNil(t, schema) {
			assert.Len(t, schema.Fields, len(fields))
			for _, field := range fields {
				metadataField, ok := schema.Fields[field]
				assert.True(t, ok)
				assert.True(t, metadataField.Filterable)
			}
		}
	})
}
