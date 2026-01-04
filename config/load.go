package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func Load() *Config {

	var k = koanf.New(".")

	// Load default values using the confmap provider.
	// We provide a flat map with the "." delimiter.
	// A nested map can be loaded by setting the delimiter to an empty string "".
	err := k.Load(confmap.Provider(defaultConfig, "."), nil)
	if err != nil {
		return nil
	}

	// Load YAML config and merge into the previously loaded config (because we can).
	yErr := k.Load(file.Provider("config.yml"), yaml.Parser())
	if yErr != nil {
		return nil
	}

	eErr := k.Load(env.Provider(".", env.Opt{
		Prefix: "GAMEAPP_",
		TransformFunc: func(k, v string) (string, any) {
			// Transform the key.
			k = strings.ReplaceAll(
				strings.ToLower(
					strings.TrimPrefix(k, "GAMEAPP_")), "_", ".")

			// Transform the value into slices, if they contain spaces.
			// Eg:_TAGS="foo bar baz" -> tags: ["foo", "bar", "baz"]
			// This is to demonstrate that string values can be transformed to any type
			// where necessary.
			if strings.Contains(v, " ") {
				return k, strings.Split(v, " ")
			}

			return k, v
		},
	}), nil)
	if eErr != nil {
		return nil
	}

	var cfg Config

	if uErr := k.Unmarshal("", &cfg); uErr != nil {
		panic(uErr)
	}

	return &cfg

}
