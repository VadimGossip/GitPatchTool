package domain

import (
	"errors"
)

const (
	OracleTablespaceType = iota + 1
	OracleDirectoryType
	OracleDbLinkType
	OracleUserType
	OracleSynonymType
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

var SchemaStrItemServerSchemaDict = map[string]ServerSchema{
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

var ServerSchemaSchemaStrItemDict = map[ServerSchema]string{
	ServerSchema{
		Server: "core",
		Schema: "vtbs",
	}: "core",
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs",
	}: "charger",
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs_adesk",
	}: "vtbs_adesk",
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs_x_alaris",
	}: "vtbs_x_alaris",
	ServerSchema{
		Server: "hpffm",
		Schema: "vtbs_bi",
	}: "vtbs_bi",
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

var DirOracleObjectTypeDict = map[string]int{
	"tablespaces":       OracleTablespaceType,
	"directories":       OracleDirectoryType,
	"dblinks":           OracleDbLinkType,
	"users":             OracleUserType,
	"synonyms":          OracleSynonymType,
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
	OracleScriptsBeforeType:    "scripts_before:",
	OracleScriptsAfterType:     "scripts_after:",
	OracleScriptsMigrationType: "scripts_migration",
}

func GetServerSchemaBySchemaStrItem(schemaStrItem string) (ServerSchema, error) {
	if val, ok := SchemaStrItemServerSchemaDict[schemaStrItem]; ok {
		return val, nil
	}
	return ServerSchema{}, UnknownSchemaStrItem
}
