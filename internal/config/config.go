package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/VadimGossip/gitPatchTool/internal/domain"
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
	var err error
	if err = viper.UnmarshalKey("work_options", &cfg); err != nil {
		return err
	}
	if err = viper.UnmarshalKey("work_options.path", &cfg.Path); err != nil {
		return err
	}
	if err = viper.UnmarshalKey("dictionaries.server_schema", &cfg.Dictionaries.ServerSchema); err != nil {
		return err
	}

	if err = viper.UnmarshalKey("dictionaries.server_schema", &cfg.Dictionaries.ServerSchema); err != nil {
		return err
	}

	cfg.Dictionaries.ServerSchemaAlias = make(map[domain.ServerSchema]string)
	for key, val := range cfg.Dictionaries.ServerSchema {
		cfg.Dictionaries.ServerSchemaAlias[domain.ServerSchema{Server: val.Server, Schema: val.Schema}] = key
	}

	var installFilename domain.ServerSchemaFilenameList
	if err = viper.UnmarshalKey("dictionaries.server_schema_filename_list.install", &installFilename); err != nil {
		return err
	}

	cfg.Dictionaries.InstallFilename, err = installFilename.BuildDictionary()
	if err != nil {
		return fmt.Errorf("install filename dictionary build error %s", err)
	}

	migrationFilename := make(domain.ServerSchemaFilenameList, 0)
	if err = viper.UnmarshalKey("dictionaries.server_schema_filename_list.migration", &migrationFilename); err != nil {
		return err
	}

	cfg.Dictionaries.MigrationFilename, err = migrationFilename.BuildDictionary()
	if err != nil {
		return fmt.Errorf("migration filename dictionary build  error %s", err)
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

func validate(cfg *domain.Config) error {
	if len(cfg.Dictionaries.ServerSchema) == 0 {
		return fmt.Errorf("empty server schema dictionary")
	}

	for key, val := range cfg.Dictionaries.ServerSchema {
		if key == "" {
			return fmt.Errorf("empty key of server schema dictionary")
		}

		if val.Server == "" {
			return fmt.Errorf("empty value server of schema dictionary")
		}

		if val.Schema == "" {
			return fmt.Errorf("empty value schema of schema dictionary")
		}
	}

	return nil
}

func Init(configDir string) (*domain.Config, error) {
	viper.SetConfigName("config")
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}

	cfg := &domain.Config{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	if err := checkAndFixPaths(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
