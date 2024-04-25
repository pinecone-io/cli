package presenters

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/pinecone"
)

func printOptionalInt(writer io.Writer, label string, value interface{}) {
	if value == nil {
		fmt.Fprintf(writer, "%s\t%s\n", label, "N/A")
		return
	}

	printValue := func(v interface{}) {
		switch v := v.(type) {
		case *int32, *int64:
			if reflect.ValueOf(v).IsNil() {
				fmt.Fprintf(writer, "%s\t%s\n", label, "null")
			} else {
				fmt.Fprintf(writer, "%s\t%d\n", label, reflect.Indirect(reflect.ValueOf(v)).Interface())
			}
		default:
			fmt.Fprintf(writer, "%s\t%s\n", label, "Invalid type")
		}
	}

	printValue(value)
}

func PrintDescribeCollectionTable(coll *pinecone.Collection) {
	writer := NewTabWriter()

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Name\t%s\n", coll.Name)
	fmt.Fprintf(writer, "State\t%s\n", ColorizeCollectionStatus(coll.Status))

	printOptionalInt(writer, "Dimension", coll.Dimension)
	printOptionalInt(writer, "Size", coll.Size)
	printOptionalInt(writer, "VectorCount", coll.VectorCount)

	fmt.Fprintf(writer, "Environment\t%s\n", coll.Environment)

	writer.Flush()
}

func ColorizeCollectionStatus(state pinecone.CollectionStatus) string {
	switch state {
	case pinecone.CollectionStatusReady:
		return style.StatusGreen(string(state))
	case pinecone.CollectionStatusInitializing, pinecone.CollectionStatusTerminating:
		return style.StatusYellow(string(state))
	}

	return string(state)
}
