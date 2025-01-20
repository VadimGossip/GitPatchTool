package schema

import (
	"strings"

	"github.com/VadimGossip/gitPatchTool/internal/domain"
)

type Service interface {
	ParseSchemaAliasLine(schemaAliasLine string) map[domain.ServerSchema]struct{}
	GetInstallFilename(serverSchema domain.ServerSchema) string
	GetMigrationFilename(serverSchema domain.ServerSchema) string
	GetSchemaAlias(serverSchema domain.ServerSchema) string
}

type service struct {
	schemaDictCfg domain.DictionariesConfig
}

func NewService(schemaDictCfg domain.DictionariesConfig) *service {
	return &service{schemaDictCfg: schemaDictCfg}
}

func (s *service) getServerSchema(alias string) (domain.ServerSchema, error) {
	if val, ok := s.schemaDictCfg.ServerSchema[alias]; ok {
		return val, nil
	}
	return domain.ServerSchema{}, domain.UnknownSchemaStrItem
}

func (s *service) ParseSchemaAliasLine(schemaAliasLine string) map[domain.ServerSchema]struct{} {
	schemaAliasLine = strings.ToLower(schemaAliasLine)
	schemaAliasLine = strings.Replace(schemaAliasLine, "schema", "", -1)
	schemaAliasLine = strings.Replace(schemaAliasLine, ":", "", -1)
	schemaAliasLine = strings.Replace(schemaAliasLine, " ", "", -1)
	schemaAliasLine = strings.Replace(schemaAliasLine, "--", "", -1)
	schemaAliasLine = strings.Replace(schemaAliasLine, "/", ",", -1)
	schemaAliasLine = strings.Replace(schemaAliasLine, "\\", ",", -1)
	if len(schemaAliasLine) > 0 {
		parts := strings.Split(schemaAliasLine, ",")
		result := make(map[domain.ServerSchema]struct{})
		for _, alias := range parts {
			serverSchema, err := s.getServerSchema(alias)
			if err == nil {
				result[serverSchema] = struct{}{}
			}
		}
		return result
	}
	return nil
}

func (s *service) GetInstallFilename(serverSchema domain.ServerSchema) string {
	return s.schemaDictCfg.InstallFilename[serverSchema]
}

func (s *service) GetMigrationFilename(serverSchema domain.ServerSchema) string {
	return s.schemaDictCfg.MigrationFilename[serverSchema]
}

func (s *service) GetSchemaAlias(serverSchema domain.ServerSchema) string {
	return s.schemaDictCfg.ServerSchemaAlias[serverSchema]
}
