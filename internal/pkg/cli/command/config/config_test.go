package config

import "context"

// mockConfigService implements ConfigService for unit tests.
// Each field controls what the corresponding method returns.
// The last* fields record the arguments of the most recent call.
type mockConfigService struct {
	// Get
	getValue          string
	getSensitive      bool
	getEnvVarName     string
	getEnvVarOverride bool
	getErr            error
	lastGetKey        string

	// GetStored — defaults to getValue/getSensitive/getErr when not explicitly set
	getStoredValue     string
	getStoredSensitive bool
	getStoredErr       error
	getStoredOverride  bool // set to true to use getStored* fields instead of get* fields
	lastGetStoredKey   string

	// Set
	setLines     []string
	setErr       error
	lastSetKey   string
	lastSetValue string

	// Unset
	unsetLines   []string
	unsetErr     error
	lastUnsetKey string

	// List
	listResult []ConfigEntry

	// Describe
	describeResult  ConfigDescription
	describeErr     error
	lastDescribeKey string
}

func (m *mockConfigService) Get(key string) (value string, sensitive bool, envVarName string, envVarOverride bool, err error) {
	m.lastGetKey = key
	return m.getValue, m.getSensitive, m.getEnvVarName, m.getEnvVarOverride, m.getErr
}

func (m *mockConfigService) GetStored(key string) (value string, sensitive bool, err error) {
	m.lastGetStoredKey = key
	if m.getStoredOverride {
		return m.getStoredValue, m.getStoredSensitive, m.getStoredErr
	}
	return m.getValue, m.getSensitive, m.getErr
}

func (m *mockConfigService) Set(ctx context.Context, key, value string) ([]string, error) {
	m.lastSetKey = key
	m.lastSetValue = value
	return m.setLines, m.setErr
}

func (m *mockConfigService) Unset(ctx context.Context, key string) ([]string, error) {
	m.lastUnsetKey = key
	return m.unsetLines, m.unsetErr
}

func (m *mockConfigService) List(includeHidden bool) []ConfigEntry {
	return m.listResult
}

func (m *mockConfigService) Describe(key string) (ConfigDescription, error) {
	m.lastDescribeKey = key
	return m.describeResult, m.describeErr
}
