package presenters

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

func ColorizeBool(b bool) string {
	if b {
		return style.SuccessStyle().Render("true")
	}
	return style.ErrorStyle().Render("false")
}

func DisplayOrNone(val any) any {
	if val == nil {
		return "<none>"
	}

	v := reflect.ValueOf(val)
	for v.IsValid() {
		switch v.Kind() {
		case reflect.Ptr, reflect.Interface:
			if v.IsNil() {
				return "<none>"
			}
			v = v.Elem()
		default:
			return v.Interface()
		}
	}

	return "<none>"
}

// FormatSize formats byte size into human-readable format
func FormatSize(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatTags formats tags for display
func FormatTags(tags *pinecone.IndexTags) string {
	if tags == nil {
		return ""
	}

	var tagStrings []string
	for key, value := range *tags {
		if value == "" {
			tagStrings = append(tagStrings, key)
		} else {
			tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return strings.Join(tagStrings, ", ")
}

// ColorizeStatus applies appropriate styling to any status string
// Works for both backup and index statuses
func ColorizeStatus(status string) string {
	// Normalize the status to lowercase for comparison
	normalizedStatus := strings.ToLower(strings.TrimSpace(status))

	switch normalizedStatus {
	// Success states (green)
	case "ready", "completed":
		return style.SuccessStyle().Render(status)

	// Warning/In-progress states (yellow)
	case "initializing", "inprogress", "in_progress", "pending", "terminating",
		"scalingdown", "scalingdownpodsize", "scalingup", "scalinguppodsize":
		return style.WarningStyle().Render(status)

	// Error states (red)
	case "failed", "initializationfailed":
		return style.ErrorStyle().Render(status)

	default:
		// If status is empty or unknown, show it as-is
		if status == "" {
			return "-"
		}
		return status
	}
}
