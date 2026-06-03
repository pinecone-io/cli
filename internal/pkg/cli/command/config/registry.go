package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	conf "github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

// ErrNoChange is returned by validateStr when the incoming value is equivalent
// to the current stored value and no write is needed.
var ErrNoChange = errors.New("no change")

// keyDescriptor describes a single user-configurable setting.
type keyDescriptor struct {
	Description     string
	LongDescription string // optional multi-paragraph detail shown by `pc config describe`
	Sensitive       bool
	Hidden          bool
	ValidValues     []string // non-nil: values shown in help; nil: any non-empty string accepted
	defaultVal      string   // the value restored by Unset; must match what getStr returns at the default
	// getStr reads the value persisted in the config file, bypassing any env var override.
	getStr func() string
	// envVarName is the environment variable that can override this key (e.g. "PINECONE_API_KEY").
	// Empty for keys with no env var binding. When non-empty, Set and Unset skip onChange because
	// writing the config file does not change the runtime value — the env var continues to win.
	envVarName string
	// validateStr normalises the incoming value and checks whether it differs from
	// the current stored value. It is pure (no I/O) and must be called before persistStr.
	// Returns ErrNoChange when the value is already current, or a validation error.
	// If nil, the value is passed through to persistStr unchanged.
	validateStr func(value string) (normalizedValue string, err error)
	// onChange is called with the prospective new value before it is persisted.
	// Returning an error aborts the operation; nothing is written to disk.
	onChange func(ctx context.Context, oldVal, newVal string) ([]string, error)
	// persistStr writes a normalised value returned by validateStr. It is only
	// called after validateStr and onChange have both succeeded.
	persistStr func(normalizedValue string)
}

// configKeys represent the set of valid configuration keys.
// This is used by lookupKey to validate keys on config commands,
// and the order of keys in the list command.
var configKeys = []string{
	"api-key",
	"color",
	"environment",
}

// configRegistry is a map of all config keys and their descriptors.
var configRegistry = map[string]keyDescriptor{
	"api-key": {
		Description: "Default API key for authenticating with Pinecone",
		LongDescription: help.Long(`
			Configure the CLI to authenticate with Pinecone using an API key.

			When set, the API key takes priority over any target context established by
			user login or service account credentials, and is used for all API calls.

			To clear the explicit API key, run 'pc config unset api-key'.

			The PINECONE_API_KEY environment variable takes precedence over any value
			stored here. When it is set, 'pc config set api-key' and 'pc config unset
			api-key' still update the stored preference, but authentication state is
			left unchanged because the environment variable controls the runtime value.
		`),
		Sensitive:  true,
		defaultVal: "",
		getStr: func() string {
			return secrets.DefaultAPIKey.GetStored()
		},
		envVarName: "PINECONE_API_KEY",
		persistStr: func(value string) {
			secrets.DefaultAPIKey.Set(value)
		},
		// Reconcile the auth context whenever the API key changes, matching the
		// behaviour of pc auth clear --api-key.
		onChange: func(_ context.Context, _, newVal string) ([]string, error) {
			if newVal != "" {
				state.AuthedUser.Update(func(u *state.TargetUser) {
					u.AuthContext = state.AuthDefaultAPIKey
				})
				return nil, nil
			}
			// Key was cleared — fall back to whichever credential remains.
			if secrets.ClientId.Get() != "" && secrets.ClientSecret.Get() != "" {
				state.AuthedUser.Update(func(u *state.TargetUser) {
					u.AuthContext = state.AuthServiceAccount
				})
			} else if secrets.GetOAuth2Token().AccessToken != "" {
				state.AuthedUser.Update(func(u *state.TargetUser) {
					u.AuthContext = state.AuthUserToken
				})
			} else {
				state.AuthedUser.Update(func(u *state.TargetUser) {
					u.AuthContext = state.AuthNone
				})
			}
			return nil, nil
		},
	},

	"color": {
		Description: "Enable or disable colored terminal output",
		ValidValues: []string{"true", "false", "on", "off", "1", "0"},
		defaultVal:  "true",
		getStr: func() string {
			return text.BoolToString(conf.Color.Get())
		},
		validateStr: func(value string) (string, error) {
			switch strings.ToLower(value) {
			case "true", "on", "1":
				return "true", nil
			case "false", "off", "0":
				return "false", nil
			default:
				return "", fmt.Errorf("invalid value %q for color; must be one of: true, false, on, off", value)
			}
		},
		persistStr: func(value string) {
			conf.Color.Set(value == "true")
		},
	},

	"environment": {
		Description: "Pinecone environment to target (production or staging)",
		LongDescription: help.Long(`
			Select which Pinecone environment the CLI talks to. Most users should
			leave this set to 'production'; 'staging' is intended for Pinecone
			internal development.

			This setting is hidden from 'pc config list' by default. Use
			'pc config list --all' to include it.

			Changing the environment clears your existing authentication state: any
			OAuth session is logged out, the default API key is cleared, and the
			target organization and project are reset. You will need to re-authenticate
			and re-target after switching.

			The PINECONE_ENVIRONMENT environment variable takes precedence over any
			value stored here. When it is set, 'pc config set environment' still
			updates the stored preference, but authentication state is left unchanged
			because the environment variable controls the runtime environment.
		`),
		Hidden:      true,
		ValidValues: []string{"production", "prod", "staging"},
		defaultVal:  "production",
		getStr: func() string {
			return conf.Environment.GetStored()
		},
		envVarName: "PINECONE_ENVIRONMENT",
		validateStr: func(value string) (string, error) {
			switch value {
			case "prod":
				value = "production"
			case "production", "staging":
				// canonical values
			default:
				return "", fmt.Errorf("invalid environment %q; must be one of: production, staging", value)
			}
			if conf.Environment.GetStored() == value {
				return "", ErrNoChange
			}
			return value, nil
		},
		persistStr: func(value string) {
			conf.Environment.Set(value)
		},
		onChange: func(ctx context.Context, _, _ string) ([]string, error) {
			var lines []string

			// Check for existing OAuth sessions and login credentials and clear them when the environment is changed.
			// Always clear any stored OAuth session when switching environments.
			// oauth.Token may fail (e.g. expired refresh, parse error) even when
			// tokens are still on disk, so we cannot rely on its success to decide
			// whether to call Logout — we call it unconditionally. The result is
			// only used to pick the right message.
			token, _ := oauth.Token(ctx)
			oauth.Logout()
			if token != nil && (token.AccessToken != "" || token.RefreshToken != "") {
				lines = append(lines, fmt.Sprintf("You have been logged out; to login again, run %s", style.Code("pc login")))
			} else {
				lines = append(lines, fmt.Sprintf("To login, run %s", style.Code("pc login")))
			}

			if secrets.DefaultAPIKey.Get() != "" {
				secrets.DefaultAPIKey.Clear()
				lines = append(lines, fmt.Sprintf("API key cleared; to set a new API key, run %s", style.Code("pc config set api-key <value>")))
			} else {
				lines = append(lines, fmt.Sprintf("To set a new API key, run %s", style.Code("pc config set api-key <value>")))
			}

			if state.TargetOrg.Get().Name != "" || state.TargetProj.Get().Name != "" {
				state.TargetOrg.Clear()
				state.TargetProj.Clear()
				lines = append(lines, fmt.Sprintf("Target organization and project cleared; to set a new target, run %s", style.Code("pc target -o myorg -p myproj")))
			}

			return lines, nil
		},
	},
}

// lookupKey returns the descriptor for name, or a descriptive error listing valid keys.
func lookupKey(name string) (keyDescriptor, error) {
	desc, ok := configRegistry[name]
	if !ok {
		return keyDescriptor{}, fmt.Errorf("unknown config key %q; valid keys are: %s", name, strings.Join(configKeys, ", "))
	}
	return desc, nil
}

// visibleKeys returns the set of config keys that are surfaced to the user.
func visibleKeys() []string {
	keys := make([]string, 0, len(configKeys))
	for _, key := range configKeys {
		if !configRegistry[key].Hidden {
			keys = append(keys, key)
		}
	}
	return keys
}

// displayValue formats a config value for human-readable output, substituting
// a placeholder when the value is empty. JSON output should use the raw value.
func displayValue(value string) string {
	if value == "" {
		return "<not set>"
	}
	return value
}

// ConfigEntry holds the key, value, and metadata for a single config setting, used by the list command.
type ConfigEntry struct {
	Key            string
	Value          string // effective value: the env var value when overriding, otherwise the stored file value
	EnvVarName     string // non-empty when an env var is available for this key through the config viper store
	EnvVarOverride bool   // true when an env var is currently overriding this key's stored file value
	Description    string
	Sensitive      bool
	Hidden         bool
}

// ConfigDescription holds full metadata for a config key, used by the describe command.
type ConfigDescription struct {
	Key             string
	Value           string // effective value: the env var value when overriding, otherwise the stored file value
	EnvVarName      string // non-empty when an env var is available for this key through the config viper store
	EnvVarOverride  bool   // true when an env var is currently overriding this key's stored file value
	Description     string
	LongDescription string
	Sensitive       bool
	ValidValues     []string
}

// ConfigService abstracts config registry operations for unit testing across
// the get, set, unset, list, and describe commands.
type ConfigService interface {
	// Get returns the effective value (env var if one is active, otherwise the stored file value)
	// along with the name of the overriding env var when one is active (empty string otherwise).
	Get(key string) (value string, sensitive bool, envVarName string, envVarOverride bool, err error)
	// GetStored returns the value persisted in the config file, bypassing env var overrides.
	GetStored(key string) (value string, sensitive bool, err error)
	Set(ctx context.Context, key, value string) (onChangeLines []string, err error)
	Unset(ctx context.Context, key string) (onChangeLines []string, err error)
	List(includeHidden bool) []ConfigEntry
	Describe(key string) (ConfigDescription, error)
}

type defaultConfigService struct{}

func newDefaultConfigService() ConfigService {
	return &defaultConfigService{}
}

func (s *defaultConfigService) Get(key string) (value string, sensitive bool, envVarName string, envVarOverride bool, err error) {
	desc, err := lookupKey(key)
	if err != nil {
		return "", false, "", false, err
	}
	if desc.envVarName != "" {
		if v := os.Getenv(desc.envVarName); v != "" {
			return v, desc.Sensitive, desc.envVarName, true, nil
		}
		return desc.getStr(), desc.Sensitive, desc.envVarName, false, nil
	}
	return desc.getStr(), desc.Sensitive, "", false, nil
}

func (s *defaultConfigService) GetStored(key string) (string, bool, error) {
	desc, err := lookupKey(key)
	if err != nil {
		return "", false, err
	}
	return desc.getStr(), desc.Sensitive, nil
}

func (s *defaultConfigService) Set(ctx context.Context, key, value string) ([]string, error) {
	desc, err := lookupKey(key)
	if err != nil {
		return nil, err
	}
	oldVal := desc.getStr()
	normalizedVal := value
	if desc.validateStr != nil {
		if normalizedVal, err = desc.validateStr(value); err != nil {
			return nil, err
		}
	}
	// Skip onChange when an env var is overriding the effective value: writing to
	// the config file won't change what the CLI uses at runtime, so side effects
	// like clearing auth state would be incorrect.
	var lines []string
	if desc.onChange != nil && (desc.envVarName == "" || os.Getenv(desc.envVarName) == "") {
		if lines, err = desc.onChange(ctx, oldVal, normalizedVal); err != nil {
			return nil, err
		}
	}
	desc.persistStr(normalizedVal)
	return lines, nil
}

func (s *defaultConfigService) Unset(ctx context.Context, key string) ([]string, error) {
	desc, err := lookupKey(key)
	if err != nil {
		return nil, err
	}
	oldVal := desc.getStr()
	newVal := desc.defaultVal
	if oldVal == newVal {
		return nil, ErrNoChange
	}
	// Skip onChange when an env var is overriding the effective value (same rationale as Set).
	var lines []string
	if desc.onChange != nil && (desc.envVarName == "" || os.Getenv(desc.envVarName) == "") {
		if lines, err = desc.onChange(ctx, oldVal, newVal); err != nil {
			return nil, err
		}
	}
	desc.persistStr(newVal)
	return lines, nil
}

func (s *defaultConfigService) List(includeHidden bool) []ConfigEntry {
	keys := visibleKeys()
	if includeHidden {
		keys = configKeys
	}
	entries := make([]ConfigEntry, 0, len(keys))
	for _, key := range keys {
		desc := configRegistry[key]
		entry := ConfigEntry{
			Key:         key,
			Value:       desc.getStr(),
			Description: desc.Description,
			Sensitive:   desc.Sensitive,
			Hidden:      desc.Hidden,
		}
		if desc.envVarName != "" {
			entry.EnvVarName = desc.envVarName
			if v := os.Getenv(desc.envVarName); v != "" {
				entry.Value = v
				entry.EnvVarOverride = true
			}
		}
		entries = append(entries, entry)
	}
	return entries
}

func (s *defaultConfigService) Describe(key string) (ConfigDescription, error) {
	desc, err := lookupKey(key)
	if err != nil {
		return ConfigDescription{}, err
	}
	d := ConfigDescription{
		Key:             key,
		Value:           desc.getStr(),
		Description:     desc.Description,
		LongDescription: desc.LongDescription,
		Sensitive:       desc.Sensitive,
		ValidValues:     desc.ValidValues,
	}
	if desc.envVarName != "" {
		d.EnvVarName = desc.envVarName
		if v := os.Getenv(desc.envVarName); v != "" {
			d.Value = v
			d.EnvVarOverride = true
		}
	}
	return d, nil
}
