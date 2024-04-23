package config

import (
	"os"
	"path/filepath"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
)

func HomeDirPath(subdir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	newPath := filepath.Join(homeDir, subdir)
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
		exit.Error(err)
	}

	if !doesFileExist(configPath) {
		err = os.MkdirAll(configPath, 0755)
		if err != nil {
			exit.Error(err)
		}
	}

	return configPath
}

func ConfigDirPath() string {
	configPath, err := HomeDirPath(".config/pinecone")
	if err != nil {
		exit.Error(err)
	}

	return configPath
}

type ConfigLocations struct {
	ConfigPath string
}

func NewConfigLocations() *ConfigLocations {
	configPath := ConfigDirPath()
	ensureConfigDir()

	return &ConfigLocations{
		ConfigPath: configPath,
	}
}
