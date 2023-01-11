package domain

import "errors"

const (
	OracleTablespaceType = iota
	OracleDirectoryType
	OracleDbLinkType
	OracleUserType
	OracleSynonymType
	OracleContextType
	OracleSequencesType
	OracleTypeType
	OracleTableType
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
	OracleTableFKType
	OracleScriptsBeforeType
	OracleScriptsAfterType
	OracleScriptsMigrationType
)

var UnknownObjectType = errors.New("can't extract object type from folder path")
var FileNotExists = errors.New("file not exists")
var SchemaNotFound = errors.New("can't parse schema from file")

type ServerSchema struct {
	Server string
	Schema string
}

type OracleObject struct {
	EpicModuleName   string
	ModuleName       string
	ObjectName       string
	ObjectType       int
	ServerSchemaList []ServerSchema
	File             File
	Errors           []string
	InstallOrder     int
}

type OracleObjectServerSchema struct {
	server string
	schema string
}
