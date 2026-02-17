package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func Load() Config {

	var k = koanf.New(".")

	// Load default values using the confmap provider.
	// We provide a flat map with the "." delimiter.
	// A nested map can be loaded by setting the delimiter to an empty string "".
	err := k.Load(confmap.Provider(defaultConfig, "."), nil)
	if err != nil {
		return Config{}
	}

	// Load YAML config and merge into the previously loaded config (because we can).
	yErr := k.Load(file.Provider("config.yml"), yaml.Parser())
	if yErr != nil {
		return Config{}
	}

	eErr := k.Load(env.Provider(".", env.Opt{
		Prefix: "GAMEAPP_",
		TransformFunc: func(k, v string) (string, any) {
			// Transform the key.
			k = strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(k, "GAMEAPP_")), "_", ".")
			// GAMEAPP_AUTH_SIGN__KEY
			// AUTH_SIGN_KEY
			// auth_sign__key
			// auth.sign..key
			str := strings.Replace(k, "..", "_", -1)
			// atun.sign_key
			if strings.Contains(v, " ") {
				return k, strings.Split(v, " ")
			}

			return str, v
		},
	}), nil)
	if eErr != nil {
		return Config{}
	}

	var cfg Config

	if uErr := k.Unmarshal("", &cfg); uErr != nil {
		panic(uErr)
	}

	return cfg

}
