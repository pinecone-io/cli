package config

import (
	"context"
	"errors"
	"fmt"
	"strings"

	conf "github.com/pinecone-io/cli/internal/pkg/utils/configuration/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

// ErrNoChange is returned by a keyDescriptor's setStr when the incoming value
// is equivalent to the current stored value and no write is needed.
var ErrNoChange = errors.New("no change")

// keyDescriptor describes a single user-configurable setting.
type keyDescriptor struct {
	Description     string
	LongDescription string // optional multi-paragraph detail shown by `pc config describe`
	Sensitive       bool
	Hidden          bool
	ValidValues     []string // non-nil: values shown in help; nil: any non-empty string accepted
	getStr          func() string
	setStr          func(value string) error // returns ErrNoChange or a validation error
	clearStr        func()
	// onChange is invoked after a successful setStr. It may call exit.Error for
	// fatal side-effect failures. Returns human-readable info lines for the user.
	onChange func(ctx context.Context, oldVal, newVal string) []string
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
		`),
		Sensitive: true,
		Hidden:    false,
		getStr: func() string {
			return secrets.DefaultAPIKey.Get()
		},
		setStr: func(value string) error {
			if value == "" {
				return fmt.Errorf("api-key value cannot be empty")
			}
			secrets.DefaultAPIKey.Set(value)
			return nil
		},
		clearStr: func() {
			secrets.DefaultAPIKey.Clear()
		},
	},
	"color": {
		Description: "Enable or disable colored terminal output",
		Sensitive:   false,
		Hidden:      false,
		ValidValues: []string{"true", "false"},
		getStr: func() string {
			return text.BoolToString(conf.Color.Get())
		},
		setStr: func(value string) error {
			var colorSetting bool
			switch strings.ToLower(value) {
			case "true", "on", "1":
				colorSetting = true
			case "false", "off", "0":
				colorSetting = false
			default:
				return fmt.Errorf("invalid value %q for color; must be one of: true, false", value)
			}
			conf.Color.Set(colorSetting)
			return nil
		},
		clearStr: func() {
			conf.Color.Clear()
		},
	},
	"environment": {
		Description: "Pinecone environment to target (production or staging)",
		LongDescription: help.Long(`
			Select which Pinecone environment the CLI talks to. Most users should
			leave this set to 'production'; 'staging' is intended for Pinecone
			internal development.

			Changing the environment clears your existing authentication state: any
			OAuth session is logged out, the default API key is cleared, and the
			target organization and project are reset. You will need to re-authenticate
			and re-target after switching.
		`),
		Sensitive:   false,
		Hidden:      true,
		ValidValues: []string{"production", "staging"},
		getStr: func() string {
			return conf.Environment.Get()
		},
		setStr: func(value string) error {
			switch value {
			case "production", "prod":
				value = "production"
			case "staging":
				// already the canonical value
			default:
				return fmt.Errorf("invalid environment %q; must be one of: production, staging", value)
			}
			if conf.Environment.Get() == value {
				return ErrNoChange
			}
			conf.Environment.Set(value)
			return nil
		},
		clearStr: func() {
			conf.Environment.Clear()
		},
		onChange: func(ctx context.Context, _, _ string) []string {
			var lines []string

			token, err := oauth.Token(ctx)
			if err != nil {
				msg.FailMsg("Error retrieving oauth token: %s", err)
				exit.Error(err, "error retrieving oauth token")
				return nil // unreachable
			}
			if token != nil && (token.AccessToken != "" || token.RefreshToken != "") {
				oauth.Logout()
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

			return lines
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
	keys := make([]string, 0, len(configRegistry))
	for key, desc := range configRegistry {
		if !desc.Hidden {
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

// ConfigEntry holds the key, value, and metadata for a single config setting,
// used by the list command.
type ConfigEntry struct {
	Key         string
	Value       string
	Description string
	Sensitive   bool
}

// ConfigDescription holds full metadata for a config key, used by the describe command.
type ConfigDescription struct {
	Key             string
	Value           string
	Description     string
	LongDescription string
	Sensitive       bool
	ValidValues     []string
}

// ConfigService abstracts config registry operations for unit testing across
// the get, set, unset, list, and describe commands.
type ConfigService interface {
	Get(key string) (value string, sensitive bool, err error)
	Set(ctx context.Context, key, value string) (onChangeLines []string, err error)
	Unset(ctx context.Context, key string) (onChangeLines []string, err error)
	List() []ConfigEntry
	Describe(key string) (ConfigDescription, error)
}

type defaultConfigService struct{}

func newDefaultConfigService() ConfigService {
	return &defaultConfigService{}
}

func (s *defaultConfigService) Get(key string) (string, bool, error) {
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
	if err := desc.setStr(value); err != nil {
		return nil, err
	}
	newVal := desc.getStr()
	if desc.onChange != nil {
		return desc.onChange(ctx, oldVal, newVal), nil
	}
	return nil, nil
}

func (s *defaultConfigService) Unset(ctx context.Context, key string) ([]string, error) {
	desc, err := lookupKey(key)
	if err != nil {
		return nil, err
	}
	oldVal := desc.getStr()
	desc.clearStr()
	newVal := desc.getStr()
	if desc.onChange != nil && oldVal != newVal {
		return desc.onChange(ctx, oldVal, newVal), nil
	}
	return nil, nil
}

func (s *defaultConfigService) List() []ConfigEntry {
	visibleKeys := visibleKeys()
	entries := make([]ConfigEntry, 0, len(visibleKeys))
	for _, key := range visibleKeys {
		desc := configRegistry[key]
		entries = append(entries, ConfigEntry{
			Key:         key,
			Value:       desc.getStr(),
			Description: desc.Description,
			Sensitive:   desc.Sensitive,
		})
	}
	return entries
}

func (s *defaultConfigService) Describe(key string) (ConfigDescription, error) {
	desc, err := lookupKey(key)
	if err != nil {
		return ConfigDescription{}, err
	}
	return ConfigDescription{
		Key:             key,
		Value:           desc.getStr(),
		Description:     desc.Description,
		LongDescription: desc.LongDescription,
		Sensitive:       desc.Sensitive,
		ValidValues:     desc.ValidValues,
	}, nil
}
