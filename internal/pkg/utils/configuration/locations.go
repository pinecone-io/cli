package configuration

import (
	"os"
	"path/filepath"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
)

func HomeDirPath(subdir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	newPath := filepath.Join(homeDir, subdir)
	log.Trace().Str("homedir", newPath).Msg("Built home directory")
	return newPath, nil
}

func doesFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func ensureConfigDir() string {
	configPath, err := HomeDirPath(".config/pinecone")
	if err != nil {
		exit.Error(err, "Error getting home directory path")
	}

	if !doesFileExist(configPath) {
		err = os.MkdirAll(configPath, 0755)
		if err != nil {
			exit.Error(err, "Error creating config directory")
		}
	}

	return configPath
}

func ConfigDirPath() string {
	configPath, err := HomeDirPath(".config/pinecone")
	if err != nil {
		exit.Error(err, "Error getting home directory path")
	}

	return configPath
}

type ConfigLocations struct {
	ConfigPath string
}

func NewConfigLocations() *ConfigLocations {
	configPath := ConfigDirPath()

	return &ConfigLocations{
		ConfigPath: configPath,
	}
}

func init() {
	ensureConfigDir()
}
