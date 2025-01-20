package domain

import (
	"errors"
	"fmt"
)

const (
	OracleTablespaceType = iota + 1
	OracleDirectoryType
	OracleDbLinkType
	OracleUserType
	OracleContextType
	OracleSequencesType
	OracleTypeType
	OracleTableType
	OracleTableFKType
	OracleMLogType
	OracleMViewType
	OraclePackageType
	OracleViewType
	OracleTriggerType
	OracleVTaskType
	OracleRowType
	OracleRoleType
	OracleFunctionType
	OracleVClogType
	OracleSynonymType
	OracleScriptsBeforeType
	OracleScriptsAfterType
	OracleScriptsMigrationType
)

var UnknownObjectType = errors.New("can't extract object type from folder path")
var FileNotExists = errors.New("file not exists")
var SchemaNotFound = errors.New("can't parse schema from file")
var UnknownSchemaStrItem = errors.New("unknown schema str item")

type ServerSchema struct {
	Server string `mapstructure:"server"`
	Schema string `mapstructure:"schema"`
}

type ServerSchemaFilename struct {
	Server   string
	Schema   string
	Filename string
}

type ServerSchemaFilenameList []ServerSchemaFilename

func (s ServerSchemaFilenameList) BuildDictionary() (map[ServerSchema]string, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("empty server schema filename list")
	}

	result := make(map[ServerSchema]string)

	for idx := range s {
		if s[idx].Server == "" {
			return nil, fmt.Errorf("empty server")
		}

		if s[idx].Schema == "" {
			return nil, fmt.Errorf("empty schema")
		}

		if s[idx].Filename == "" {
			return nil, fmt.Errorf("empty filename")
		}

		result[ServerSchema{
			Server: s[idx].Server,
			Schema: s[idx].Schema,
		}] = s[idx].Filename
	}

	return result, nil
}

type OracleObject struct {
	EpicModuleName   string
	ModuleName       string
	ObjectName       string
	ObjectType       int
	ServerSchemaList []ServerSchema
	File             OracleFile
	Errors           []string
	InstallOrder     int
}

var DirOracleObjectTypeDict = map[string]int{
	"tablespaces":       OracleTablespaceType,
	"directories":       OracleDirectoryType,
	"dblinks":           OracleDbLinkType,
	"users":             OracleUserType,
	"contexts":          OracleContextType,
	"sequences":         OracleSequencesType,
	"types":             OracleTypeType,
	"tables":            OracleTableType,
	"mlogs":             OracleMLogType,
	"mviews":            OracleMViewType,
	"packages":          OraclePackageType,
	"views":             OracleViewType,
	"triggers":          OracleTriggerType,
	"vtbs_tasks":        OracleVTaskType,
	"rows":              OracleRowType,
	"roles":             OracleRoleType,
	"functions":         OracleFunctionType,
	"vtbs_clogs":        OracleVClogType,
	"tables.fk":         OracleTableFKType,
	"synonyms":          OracleSynonymType,
	"scripts_before":    OracleScriptsBeforeType,
	"scripts_after":     OracleScriptsAfterType,
	"scripts_migration": OracleScriptsMigrationType,
}

var OracleObjectTypeDirDict = map[int]string{
	OracleTablespaceType:       "tablespaces",
	OracleDirectoryType:        "directories",
	OracleDbLinkType:           "dblinks",
	OracleUserType:             "users",
	OracleSynonymType:          "synonyms",
	OracleContextType:          "contexts",
	OracleSequencesType:        "sequences",
	OracleTypeType:             "types",
	OracleTableType:            "tables",
	OracleMLogType:             "mlogs",
	OracleMViewType:            "mviews",
	OraclePackageType:          "packages",
	OracleViewType:             "views",
	OracleTriggerType:          "triggers",
	OracleVTaskType:            "vtbs_tasks",
	OracleRowType:              "rows",
	OracleRoleType:             "roles",
	OracleFunctionType:         "functions",
	OracleVClogType:            "vtbs_clogs",
	OracleTableFKType:          "tables.fk",
	OracleScriptsBeforeType:    "scripts_before",
	OracleScriptsAfterType:     "scripts_after",
	OracleScriptsMigrationType: "scripts_migration",
}
