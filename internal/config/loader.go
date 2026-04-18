package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Load() (Config, error) {
	v := viper.New()
	defaults := Default()

	v.SetDefault("language", defaults.Language)
	v.SetDefault("format", defaults.Format)
	v.SetDefault("max_nodes", defaults.MaxNodes)

	v.SetConfigName("diagramr.config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix("DIAGRAMR")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
