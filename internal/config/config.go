package config

import (
	"errors"
	"fmt"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
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

func checkAndFixPaths(cfg *domain.Config) error {
	if !strings.HasSuffix(cfg.Path.RootDir, string(os.PathSeparator)) {
		cfg.Path.RootDir = cfg.Path.RootDir + string(os.PathSeparator)
	}
	if !strings.HasSuffix(cfg.Path.InstallDir, string(os.PathSeparator)) {
		cfg.Path.InstallDir = cfg.Path.InstallDir + string(os.PathSeparator)
	}

	if cfg.Mode == "patch" {
		if !strings.HasPrefix(cfg.Path.InstallDir, cfg.Path.RootDir) {
			return errors.New(fmt.Sprintf("install dir %s not in git dir %s", cfg.Path.InstallDir, cfg.Path.RootDir))
		}
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

	if err := checkAndFixPaths(&cfg); err != nil {
		return nil, err
	}
	logrus.Infof("Config %v", cfg)

	return &cfg, nil
}
