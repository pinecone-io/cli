package configuration

import (
	"encoding/json"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/spf13/viper"
)

type Property interface {
	Init()
	Clear()
}

type ConfigProperty[T any] struct {
	DefaultValue T
	KeyName      string
	ViperStore   *viper.Viper
}

func (c ConfigProperty[T]) Init() {
	log.Trace().Str("key", c.KeyName).Msg("Setting default value for property")
	c.ViperStore.SetDefault(c.KeyName, c.DefaultValue)
}

func (c ConfigProperty[T]) Set(value T) {
	log.Trace().Str("key", c.KeyName).Msg("Setting value for property")
	c.ViperStore.Set(c.KeyName, value)
	err := c.ViperStore.WriteConfig()
	if err != nil {
		exit.Error(err, "Error writing config file")
	}
}

func (c ConfigProperty[T]) Get() T {
	log.Trace().Str("key", c.KeyName).Msg("Reading value for property")
	return c.ViperStore.Get(c.KeyName).(T)
}

func (c ConfigProperty[T]) Clear() {
	c.Set(c.DefaultValue)
}

type MarshaledProperty[T any] struct {
	DefaultValue T
	KeyName      string
	ViperStore   *viper.Viper
}

func (c MarshaledProperty[T]) Init() {
	log.Trace().Str("key", c.KeyName).Msg("Setting default value for property")
	bytes, err := json.Marshal(c.DefaultValue)
	if err != nil {
		exit.Errorf(err, "Error marshalling value for property %s", c.KeyName)
	}
	c.ViperStore.SetDefault(c.KeyName, string(bytes))
}

func (c MarshaledProperty[T]) Set(value T) {
	log.Trace().Str("key", c.KeyName).Msg("Setting value for property")
	bytes, err := json.Marshal(value)
	if err != nil {
		exit.Errorf(err, "Error marshalling value for property %s", c.KeyName)
	}
	c.ViperStore.Set(c.KeyName, string(bytes))
	err = c.ViperStore.WriteConfig()
	if err != nil {
		exit.Errorf(err, "Error writing config file")
	}
}

// Update loads the current or default value, lets the caller mutate it, then writes it back
func (c MarshaledProperty[T]) Update(mut func(*T)) {
	curr := c.Get()
	mut(&curr)
	c.Set(curr)
}

func (c MarshaledProperty[T]) Get() T {
	log.Trace().Str("key", c.KeyName).Msg("Reading value for property")
	str := c.ViperStore.GetString(c.KeyName)
	if str == "" {
		return c.DefaultValue
	}

	var value T
	err := json.Unmarshal([]byte(str), &value)
	if err != nil {
		exit.Errorf(err, "Error unmarshalling value for property %s", c.KeyName)
	}
	return value
}

func (c MarshaledProperty[T]) Clear() {
	c.Set(c.DefaultValue)
}
