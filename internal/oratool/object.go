package oratool

import "github.com/VadimGossip/gitPatchTool/internal/domain"

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
	OracleScriptsBeforeType
	OracleScriptsAfterType
	OracleScriptsMigrationType
)

type Object struct {
	EpicModuleName string
	ModuleName     string
	ObjectType     int
	Schema         string
	Server         string
	File           domain.File
}
