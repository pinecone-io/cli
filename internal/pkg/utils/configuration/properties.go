package configuration

import (
	"encoding/json"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
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
	c.ViperStore.SetDefault(c.KeyName, c.DefaultValue)
}

func (c ConfigProperty[T]) Set(value T) {
	c.ViperStore.Set(c.KeyName, value)
	err := c.ViperStore.WriteConfig()
	if err != nil {
		exit.Error(err)
	}
}

func (c ConfigProperty[T]) Get() T {
	return c.ViperStore.Get(c.KeyName).(T)
}

func (c ConfigProperty[T]) Clear() {
	c.Set(c.DefaultValue)
}

type MarshaledProperty[T any] struct {
	DefaultValue *T
	KeyName      string
	ViperStore   *viper.Viper
}

func (c MarshaledProperty[T]) Init() {
	c.ViperStore.SetDefault(c.KeyName, *c.DefaultValue)
}

func (c MarshaledProperty[T]) Set(value *T) {
	bytes, err := json.Marshal(value)
	if err != nil {
		exit.Error(err)
	}
	c.ViperStore.Set(c.KeyName, string(bytes))
	err = c.ViperStore.WriteConfig()
	if err != nil {
		exit.Error(err)
	}
}

func (c MarshaledProperty[T]) Get() T {
	bytes := []byte(c.ViperStore.GetString(c.KeyName))
	var value T
	err := json.Unmarshal(bytes, &value)
	if err != nil {
		exit.Error(err)
	}
	return value
}

func (c MarshaledProperty[T]) Clear() {
	c.Set(c.DefaultValue)
}
