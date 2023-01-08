package config

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func parseConfigFile(configDir string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func unmarshal(cfg *domain.Config) error {
	if err := viper.UnmarshalKey("work_options", &cfg); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("work_options.path", &cfg.Path); err != nil {
		return err
	}

	return nil
}

func Init(configDir string) (*domain.Config, error) {
	viper.SetConfigName("config")
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}

	var cfg domain.Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	logrus.Infof("Config %v", cfg)

	return &cfg, nil
}
