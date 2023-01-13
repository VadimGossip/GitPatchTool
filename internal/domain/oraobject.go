package domain

import (
	"errors"
)

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
var UnknownSchemaStrItem = errors.New("unknown schema str item")

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
	File             OracleFile
	Errors           []string
	InstallOrder     int
}

type OracleObjectServerSchema struct {
	server string
	schema string
}

var SchemaStrServerSchemaDict = map[string]ServerSchema{
	"core": {
		Server: "core",
		Schema: "vtbs",
	},
	"charger": {
		Server: "hpffm",
		Schema: "vtbs",
	},
	"hpffm": {
		Server: "hpffm",
		Schema: "vtbs",
	},
	"vtbs_bi": {
		Server: "hpffm",
		Schema: "vtbs_bi",
	},
	"vtbs_x_alaris": {
		Server: "hpffm",
		Schema: "vtbs_x_alaris",
	},
	"xalaris": {
		Server: "hpffm",
		Schema: "vtbs_x_alaris",
	},
	"adesk": {
		Server: "hpffm",
		Schema: "vtbs_adesk",
	},
	"vtbs_adesk": {
		Server: "hpffm",
		Schema: "vtbs_adesk",
	},
	"reporter": {
		Server: "hpffm",
		Schema: "vtbs_adesk",
	},
}

var ServerSchemaInstallFilenameDict = map[ServerSchema]string{
	ServerSchema{
		Server: "core",
		Schema: "vtbs",
	}: VtbsCoreInstallName,
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs",
	}: VtbsHpffmInstallName,
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs_adesk",
	}: VtbsAdeskHpffmInstallName,
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs_x_alaris",
	}: VtbsXAlarisHpffmInstallName,
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs_bi",
	}: VtbsBiHpffmInstallName,
}

var ServerSchemaMigrationFilenameDict = map[ServerSchema]string{
	ServerSchema{
		Server: "core",
		Schema: "vtbs",
	}: VtbsCoreMigrationName,
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs",
	}: VtbsHpffmMigrationName,
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs_adesk",
	}: VtbsAdeskHpffmMigrationName,
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs_x_alaris",
	}: VtbsXAlarisHpffmMigrationName,
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs_bi",
	}: VtbsBiHpffmMigrationName,
}

func GetServerSchemaBySchemaStrItem(schemaStrItem string) (ServerSchema, error) {
	if val, ok := SchemaStrServerSchemaDict[schemaStrItem]; ok {
		return val, nil
	}
	return ServerSchema{}, UnknownSchemaStrItem
}

func GetInstallFileNameByServerSchema(serverSchema ServerSchema) string {
	if val, ok := ServerSchemaInstallFilenameDict[serverSchema]; ok {
		return val
	}
	return ""
}

func GetMigrationFileNameByServerSchema(serverSchema ServerSchema) string {
	if val, ok := ServerSchemaMigrationFilenameDict[serverSchema]; ok {
		return val
	}
	return ""
}
