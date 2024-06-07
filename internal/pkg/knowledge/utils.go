package knowledge

import (
	"path/filepath"
	"strings"
)

func IsSupportedFile(filePath string) bool {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(filePath)))

	return ext == ".pdf" || ext == ".txt"
}
